package pkg

import (
	"context"
	"log/slog"
)

type CtxKey string

const (
	MethodKey = "method"
	HostKey   = "host"
)

type CtxHandler struct{ slog.Handler }

func (h CtxHandler) Handle(ctx context.Context, r slog.Record) error {
	if v := ctx.Value(MethodKey); v != nil {
		r.AddAttrs(slog.Any("method", v))
	}
	if v := ctx.Value(HostKey); v != nil {
		r.AddAttrs(slog.Any("host", v))
	}
	return h.Handler.Handle(ctx, r)
}
