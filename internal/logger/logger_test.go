package logger_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/selesy/asdf-go-install/internal/logger"
)

func TestCachingHandler(t *testing.T) {
	t.Parallel()

	h := logger.NewCachingHandler()

	log1 := slog.New(h)
	entries(t, log1)

	log2, buf := NewLogger(t, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	require.NoError(t, h.Transfer(context.Background(), log2))

	golden.Assert(t, buf.String(), "test_logger.log")
}

func TestLogger(t *testing.T) {
	t.Parallel()

	log, buf := NewLogger(t, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	entries(t, log)

	golden.Assert(t, buf.String(), "test_logger.log")
}

func entries(t *testing.T, log *slog.Logger) {
	t.Helper()

	log.Debug("Debug should print")
	log.Info("Info should print")
	log.Warn("Warn should print")
	log.Error("Error should print")

	log = log.With(slog.String("k1", "v1"), slog.String("k2", "v2"))
	log.Info("WithAttrs attributes should work")

	log = log.WithGroup("grp")
	log.Info("WithGroup should add a prefix to attributes", slog.String("k3", "v3"))

	log = log.With(slog.String("k4", "v4"))
	log.Info("WithGroup and WithAttrs", slog.String("k5", "v5"))
}

func NewLogger(t *testing.T, opts *slog.HandlerOptions) (*slog.Logger, *bytes.Buffer) {
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

func (h *handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

func (h *handler) Handle(ctx context.Context, rec slog.Record) error {
	rec.Time = *h.ts
	err := h.h.Handle(ctx, rec)
	*h.ts = h.ts.Add(time.Second)

	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		h:  h.h.WithAttrs(attrs),
		ts: h.ts,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		h:  h.h.WithGroup(name),
		ts: h.ts,
	}
}
