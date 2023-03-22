package agi_test

import (
	"testing"

	"github.com/selesy/asdf-go-install/internal/agi"
	"github.com/stretchr/testify/require"
)

func getConfig() *agi.Config {
	return &agi.Config{
		Name:    "example",
		Package: "https://example.com/org/repo.git",
	}
}

//nolint:paralleltest
func TestConfig_Write_Read(t *testing.T) {
	path := t.TempDir()
	exp := getConfig()

	t.Run("Write", func(t *testing.T) {
		require.NoError(t, exp.Write(path))
	})

	t.Run("Read", func(t *testing.T) {
		act := new(agi.Config)

		require.NoError(t, act.Read(path))
		require.Equal(t, exp, act)
	})
}

func TestConfig_Read(t *testing.T) {
	t.Parallel()

	t.Run("succeeds with valid file", func(t *testing.T) {
		t.Parallel()

		act, exp := new(agi.Config), getConfig()
		require.NoError(t, act.Read("./testdata/succeeds_with_valid_file"))
		require.Equal(t, exp, act)
	})

	t.Run("errors on unmarshal with invalid JSON", func(t *testing.T) {
		t.Parallel()

		const exp = "failed to unmarshal configuration file\nunexpected end of JSON input"

		act := new(agi.Config)
		err := act.Read("./testdata/errors_on_unmarshal_with_invalid_json")
		require.EqualError(t, err, exp)
	})

	t.Run("errors from non-existent file", func(t *testing.T) {
		t.Parallel()

		const exp = "failed to read configuration file\n" +
			"open testdata/errors_on_non_existent_file/.config.json: no such file or directory"

		act := new(agi.Config)
		err := act.Read("./testdata/errors_on_non_existent_file")
		require.EqualError(t, err, exp)
	})
}

func TestConfig_Write(t *testing.T) {
	t.Parallel()

	// t.Run("errors with bad struct", func(t *testing.T) {
	// 	t.Parallel()

	// 	act := &agi.Config{
	// 		Name:    "\b",
	// 		Package: "\b",
	// 	}
	// 	path := t.TempDir()
	// 	err := act.Write(path)
	// 	require.EqualError(t, err, "")
	// })

	t.Run("errors with bad path", func(t *testing.T) {
		t.Parallel()

		const exp = "failed to write configuration file\n" +
			"open testdata/dir-does-not-exist/.config.json: no such file or directory"

		act, path := new(agi.Config), "./testdata/dir-does-not-exist"
		err := act.Write(path)
		require.EqualError(t, err, exp)
	})
}
