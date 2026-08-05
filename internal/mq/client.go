package mq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func Connect(ctx context.Context, url, queue string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &Client{conn: conn, channel: ch, queue: queue}, nil
}

func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		_ = c.conn.Close()
		return err
	}
	return c.conn.Close()
}

func (c *Client) Publish(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	return c.channel.PublishWithContext(ctx, "", c.queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (c *Client) Consume(ctx context.Context, handler func(ctx context.Context, msg Message) error) error {
	msgs, err := c.channel.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case delivery, ok := <-msgs:
			if !ok {
				return fmt.Errorf("consume channel closed")
			}

			var msg Message
			if err := json.Unmarshal(delivery.Body, &msg); err != nil {
				_ = delivery.Nack(false, false)
				continue
			}

			if err := handler(ctx, msg); err != nil {
				_ = delivery.Nack(false, true)
			} else {
				_ = delivery.Ack(false)
			}
		}
	}
}
