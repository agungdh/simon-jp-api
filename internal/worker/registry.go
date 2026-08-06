package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"simon-jp-api/internal/mq"
)

type Handler interface {
	Type() string
	Handle(ctx context.Context, data json.RawMessage) error
}

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

func (r *Registry) Register(handler Handler) {
	r.handlers[handler.Type()] = handler
}

func (r *Registry) Dispatch(ctx context.Context, msg mq.Message) error {
	handler, ok := r.handlers[msg.Type]
	if !ok {
		slog.Warn("no handler for message type, dropping", "type", msg.Type)
		return nil
	}
	if err := handler.Handle(ctx, msg.Data); err != nil {
		return fmt.Errorf("handler %s: %w", msg.Type, err)
	}
	return nil
}
