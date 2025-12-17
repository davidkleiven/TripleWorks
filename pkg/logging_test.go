package pkg

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtxHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), MethodKey, "PUT")
	ctx = context.WithValue(ctx, HostKey, "127.0.0.1")

	var buf bytes.Buffer
	handler := CtxHandler{slog.NewJSONHandler(&buf, nil)}
	handler.Handle(ctx, slog.Record{Message: "this is a message"})
	assert.Contains(t, buf.String(), "PUT")
	assert.Contains(t, buf.String(), "127.0.0.1")
}
