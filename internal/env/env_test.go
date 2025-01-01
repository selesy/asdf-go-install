package env_test

import (
	"log/slog"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/selesy/asdf-go-install/internal/env"
	"github.com/selesy/asdf-go-install/internal/logger/loggertest"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("passes with only required environment variables", func(t *testing.T) {
		t.Parallel()

		var (
			userDir                        = filepath.Join("/", "home", "user")
			dirVal                         = filepath.Join(userDir, ".asdf")
			dataDirVal                     = dirVal
			configFileVal                  = filepath.Join(userDir, ".asdfrc")
			defaultToolVersionsFilenameVal = ".tool-versions"
		)

		log, buf := loggertest.New(t, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		envVars := []string{
			"ASDF_DIR=" + dirVal,
			"ASDF_DATA_DIR=" + dataDirVal,
			"ASDF_CONFIG_FILE=" + configFileVal,
			"ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=" + defaultToolVersionsFilenameVal,
		}

		e, err := env.New(log, envVars)
		require.NoError(t, err)
		assert.NotNil(t, e)

		assert.Equal(t, dirVal, e.Dir())
		assert.Equal(t, dataDirVal, e.DataDir())
		assert.Equal(t, configFileVal, e.ConfigFile())
		assert.Equal(t, defaultToolVersionsFilenameVal, e.DefaultToolVersionsFilename())

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

		golden.Assert(t, buf.String(), "default-env-vars.log")
	})

	t.Run("fails without required environment variables", func(t *testing.T) {
		t.Parallel()

		log, _ := loggertest.New(t, &slog.HandlerOptions{})

		envVars := []string{
			"NOT_STUFF=this envvar is not a candidate",
			"AGI_STUFF=this envvar is a candidate",
			"ASDF_STUFF=this envvar is a candidate",
		}

		env, err := env.New(log, envVars)
		t.Log(reflect.TypeOf(err).Elem().Name())
		require.IsType(t, validator.ValidationErrors{}, err)
		assert.Len(t, err, 4)
		assert.Nil(t, env)
	})
}

func TestLogFormat_UnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inp    string
		expErr error
		expVal env.LogFormat
	}{
		{name: "valid LogFormat", inp: "json", expErr: nil, expVal: env.LogFormatJSON},
		{name: "any case", inp: "cOlOrIzEd", expErr: nil, expVal: env.LogFormatColorized},
		{name: "fails with unknown value", inp: "unknown", expErr: env.ErrInvalidLogFormat, expVal: env.LogFormat(0)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var f env.LogFormat

			err := f.UnmarshalText([]byte(test.inp))
			require.ErrorIs(t, err, test.expErr)
			assert.Equal(t, test.expVal, f)
		})
	}
}

func TestInstallType_UnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inp    string
		expErr error
		expVal env.InstallType
	}{
		{name: "valid InstallType", inp: "version", expErr: nil, expVal: env.InstallTypeVersion},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var i env.InstallType

			err := i.UnmarshalText([]byte(test.inp))
			require.ErrorIs(t, err, test.expErr)
			assert.Equal(t, test.expVal, i)
		})
	}
}
