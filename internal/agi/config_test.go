package agi_test

import (
	"testing"

	"github.com/selesy/asdf-go-install/internal/agi"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestConfig(t *testing.T) {
	path := t.TempDir()
	exp := &agi.Config{
		Name:    "example",
		Package: "https://example.com/path",
	}

	t.Run("Write", func(t *testing.T) {
		require.NoError(t, exp.Write(path))
	})

	t.Run("Read", func(t *testing.T) {
		act := new(agi.Config)

		require.NoError(t, act.Read(path))
	})
}
