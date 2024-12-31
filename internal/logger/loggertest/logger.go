// Package loggertest provides an idempotent slog.Logger for use in
// testing.
package loggertest

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// New creates a slog.Logger that can be used for test output.
//
// To make tests idempotent, a custom handler is used that overrides the
// way each logs timestamp is created.  Instead of using the system's
// clock, each log message has the time incremented by one second,
// starting at the Unix epoch.
func New(t *testing.T, opts *slog.HandlerOptions) (*slog.Logger, *bytes.Buffer) {
	t.Helper()

	buf := &bytes.Buffer{}

	ts, err := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
	require.NoError(t, err)

	return slog.New(&handler{
		h:  slog.NewTextHandler(buf, opts),
		ts: &ts,
	}), buf
}

var _ slog.Handler = (*handler)(nil)

type handler struct {
	h  slog.Handler
	ts *time.Time
}

// Enabled implements slog.Handler.
func (h *handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

// Handle implements slog.Handler.
func (h *handler) Handle(ctx context.Context, rec slog.Record) error {
	rec.Time = *h.ts
	err := h.h.Handle(ctx, rec)
	*h.ts = h.ts.Add(time.Second)

	return err
}

// WithAttrs implements slog.Handler.
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		h:  h.h.WithAttrs(attrs),
		ts: h.ts,
	}
}

// WithGroup implements slog.Handler.
func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		h:  h.h.WithGroup(name),
		ts: h.ts,
	}
}
