package worker

import (
	"context"
	"encoding/json"
	"fmt"

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
		return fmt.Errorf("no handler registered for type %q", msg.Type)
	}
	return handler.Handle(ctx, msg.Data)
}
