package agi_test

import (
	"path/filepath"
	"testing"

	"github.com/selesy/asdf-go-install/internal/agi"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	path := filepath.Join(t.TempDir())
	exp := &agi.Config{
		Name:    "example",
		Package: "https://example.com/path",
	}

	t.Run("Write", func(t *testing.T) {
		require.NoError(t, exp.Write(path))
	})

	t.Run("Read", func(t *testing.T) {
		act := &agi.Config{}

		require.NoError(t, act.Read(path))
	})
}
