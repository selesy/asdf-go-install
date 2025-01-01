package envtest_test

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/selesy/asdf-go-install/internal/env"
	"github.com/selesy/asdf-go-install/internal/env/envtest"
	"github.com/selesy/asdf-go-install/internal/logger/loggertest"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("with defaults", func(t *testing.T) {
		t.Parallel()

		log, _ := loggertest.New(t, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		e := envtest.New(t, log, []string{})

		expDir, expCfg := expected(t, "user")

		assert.Equal(t, expDir, e.Dir())
		assert.Equal(t, expDir, e.DataDir())
		assert.Equal(t, expCfg, e.ConfigFile())
		assert.Equal(t, ".tool-versions", e.DefaultToolVersionsFilename())
		assertZeroes(t, e)
	})

	t.Run("with overrides", func(t *testing.T) {
		t.Parallel()

		log, _ := loggertest.New(t, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		e := envtest.New(t, log, []string{
			"ASDF_DIR=/home/other/.asdf",
			"ASDF_DATA_DIR=/home/other/.asdf",
			"ASDF_CONFIG_FILE=/home/other/.asdfrc",
			"ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=" + ".other-versions",
		})

		expDir, expCfg := expected(t, "other")

		assert.Equal(t, expDir, e.Dir())
		assert.Equal(t, expDir, e.DataDir())
		assert.Equal(t, expCfg, e.ConfigFile())
		assert.Equal(t, ".other-versions", e.DefaultToolVersionsFilename())
		assertZeroes(t, e)
	})
}

func assertZeroes(t *testing.T, e *env.Env) {
	t.Helper()

	assert.Zero(t, e.InstallType())
	assert.Zero(t, e.InstallVersion())
	assert.Zero(t, e.InstallPath())
	assert.Zero(t, e.Concurrency())
	assert.Zero(t, e.DownloadPath())
	assert.Zero(t, e.PluginPath())
	assert.Zero(t, e.PluginSourceURL())
	assert.Zero(t, e.PluginPrevRef())
	assert.Zero(t, e.PluginPostRef())
	assert.Zero(t, e.CmdFile())
}

func expected(t *testing.T, username string) (dir, config string) {
	t.Helper()

	return filepath.Join("/", "home", username, ".asdf"),
		filepath.Join("/", "home", username, ".asdfrc")
}
