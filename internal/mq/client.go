package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultRetryDelay = 2 * time.Second
	maxRetryDelay     = 30 * time.Second
	confirmTimeout    = 10 * time.Second
)

type Options struct {
	Prefetch   int
	Consumers  int
	MaxRetries int
}

type Client struct {
	url        string
	queue      string
	dlq        string
	prefetch   int
	consumers  int
	maxRetries int

	publishMu sync.Mutex

	mu       sync.RWMutex
	conn     *amqp.Connection
	ch       *amqp.Channel
	confirms <-chan amqp.Confirmation
	done     chan struct{}
}

func Connect(ctx context.Context, url, queue string, opts Options) (*Client, error) {
	c := &Client{
		url:        url,
		queue:      queue,
		dlq:        queue + ".dlq",
		prefetch:   opts.Prefetch,
		consumers:  opts.Consumers,
		maxRetries: opts.MaxRetries,
		done:       make(chan struct{}),
	}
	if c.prefetch <= 0 {
		c.prefetch = 10
	}
	if c.consumers <= 0 {
		c.consumers = 1
	}
	if c.maxRetries <= 0 {
		c.maxRetries = 5
	}

	if err := c.connect(ctx); err != nil {
		return nil, err
	}

	go c.reconnectLoop(ctx)
	return c, nil
}

func (c *Client) connect(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open channel: %w", err)
	}

	for _, q := range []string{c.queue, c.dlq} {
		if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return fmt.Errorf("declare queue %s: %w", q, err)
		}
	}

	if err := ch.Qos(c.prefetch, 0, false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("set qos: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("enable publisher confirms: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.ch = ch
	c.confirms = ch.NotifyPublish(make(chan amqp.Confirmation, c.prefetch))
	c.mu.Unlock()

	return nil
}

func (c *Client) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	c.mu.Lock()
	ch, conn := c.ch, c.conn
	c.ch, c.conn = nil, nil
	c.mu.Unlock()

	var err error
	if ch != nil {
		if cerr := ch.Close(); cerr != nil {
			err = cerr
		}
	}
	if conn != nil {
		if cerr := conn.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	return err
}

func (c *Client) reconnectLoop(ctx context.Context) {
	delay := defaultRetryDelay
	for {
		c.mu.RLock()
		conn, ch := c.conn, c.ch
		c.mu.RUnlock()
		if conn == nil || ch == nil {
			return
		}

		connClose := conn.NotifyClose(make(chan *amqp.Error, 1))
		chClose := ch.NotifyClose(make(chan *amqp.Error, 1))

		select {
		case <-c.done:
			return
		case <-ctx.Done():
			return
		case <-connClose:
		case <-chClose:
		}

		c.mu.Lock()
		c.conn, c.ch, c.confirms = nil, nil, nil
		c.mu.Unlock()

		for {
			select {
			case <-c.done:
				return
			case <-ctx.Done():
				return
			default:
			}

			if err := c.connect(ctx); err != nil {
				slog.Warn("mq reconnect failed", "error", err, "next_retry", delay.String())
				select {
				case <-time.After(delay):
				case <-c.done:
					return
				case <-ctx.Done():
					return
				}
				delay *= 2
				if delay > maxRetryDelay {
					delay = maxRetryDelay
				}
				continue
			}

			slog.Info("mq reconnected")
			delay = defaultRetryDelay
			break
		}
	}
}

func (c *Client) Publish(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	return c.publishRaw(ctx, c.queue, body)
}

func (c *Client) publishRaw(ctx context.Context, queue string, body []byte) error {
	c.publishMu.Lock()
	defer c.publishMu.Unlock()

	c.mu.RLock()
	ch, confirms := c.ch, c.confirms
	c.mu.RUnlock()
	if ch == nil {
		return fmt.Errorf("mq not connected")
	}

	if err := ch.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("publish to %s: %w", queue, err)
	}

	select {
	case conf := <-confirms:
		if !conf.Ack {
			return fmt.Errorf("publish to %s nacked by broker", queue)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(confirmTimeout):
		return fmt.Errorf("publish to %s: confirm timeout", queue)
	}
}

func (c *Client) Consume(ctx context.Context, handler func(ctx context.Context, msg Message) error) error {
	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		c.mu.RLock()
		ch := c.ch
		c.mu.RUnlock()
		if ch == nil {
			select {
			case <-time.After(defaultRetryDelay):
			case <-ctx.Done():
				return nil
			}
			continue
		}

		msgs, err := ch.Consume(c.queue, "", false, false, false, false, nil)
		if err != nil {
			slog.Warn("mq consume failed", "error", err)
			select {
			case <-time.After(defaultRetryDelay):
			case <-ctx.Done():
				return nil
			}
			continue
		}

		if err := c.runWorkers(ctx, msgs, handler); err != nil {
			return err
		}
	}
}

func (c *Client) runWorkers(ctx context.Context, msgs <-chan amqp.Delivery, handler func(ctx context.Context, msg Message) error) error {
	var wg sync.WaitGroup
	done := make(chan struct{})
	for i := 0; i < c.consumers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-done:
					return
				case delivery, ok := <-msgs:
					if !ok {
						return
					}
					c.handle(ctx, delivery, handler)
				}
			}
		}()
	}

	wg.Wait()
	close(done)
	return ctx.Err()
}

func (c *Client) handle(ctx context.Context, delivery amqp.Delivery, handler func(ctx context.Context, msg Message) error) {
	var msg Message
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		_ = delivery.Ack(false)
		return
	}

	if err := handler(ctx, msg); err != nil {
		slog.Warn("mq job failed", "type", msg.Type, "error", err, "retry", msg.RetryCount)

		if msg.RetryCount < c.maxRetries {
			msg.RetryCount++
			body, err := json.Marshal(msg)
			if err != nil {
				_ = delivery.Nack(false, true)
				return
			}
			if err := c.publishRaw(ctx, c.queue, body); err != nil {
				_ = delivery.Nack(false, true)
				return
			}
			_ = delivery.Ack(false)
			return
		}

		body, err := json.Marshal(msg)
		if err != nil {
			_ = delivery.Ack(false)
			return
		}
		if err := c.publishRaw(ctx, c.dlq, body); err != nil {
			_ = delivery.Nack(false, true)
			return
		}
		_ = delivery.Ack(false)
		return
	}

	_ = delivery.Ack(false)
}
