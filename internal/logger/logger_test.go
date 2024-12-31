package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/selesy/asdf-go-install/internal/logger"
	"github.com/selesy/asdf-go-install/internal/logger/loggertest"
)

func TestCachingHandler(t *testing.T) {
	t.Parallel()

	h := logger.NewCachingHandler()

	log1 := slog.New(h)
	entries(t, log1)

	log2, buf := loggertest.New(t, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	require.NoError(t, h.Transfer(context.Background(), log2))

	golden.Assert(t, buf.String(), "test_logger.log")
}

func TestLogger(t *testing.T) {
	t.Parallel()

	log, buf := loggertest.New(t, &slog.HandlerOptions{
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
