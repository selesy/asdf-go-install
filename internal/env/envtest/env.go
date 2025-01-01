// Package envtest prepares an env.Env that can be used during testing.
package envtest

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/selesy/asdf-go-install/internal/env"
)

// New creates an Env for use during testing.
//
// Four environment variables are required - these are assigned as
// follows:
//
//   - ASDF_DIR defaults to /home/user/.asdf
//   - ASDF_DATA_DIR defaults to /home/user/.asdf
//   - ASDF_CONFIG_FILE defaults to /home/user/.asdfrc
//   - ASDF_DEFAULT_TOOL_VERSIONS_FILE defaults to .tool-versions
//
// Any or all of these variables can be overridden by including
// new values for them in the provided environment variable list.
func New(t *testing.T, log *slog.Logger, environ []string) *env.Env {
	t.Helper()

	var vars = make(map[string]struct{}, len(environ))

	for _, v := range environ {
		left, _, _ := strings.Cut(v, "=")
		vars[left] = struct{}{}
	}

	for _, k := range keys() {
		if _, ok := vars[k.str()]; !ok {
			environ = append(environ, k.def())
		}
	}

	e, err := env.New(log, environ)
	require.NoError(t, err)

	return e
}

type key string

const (
	dirKey                         key = "ASDF_DIR"
	dataDirKey                     key = "ASDF_DATA_DIR"
	configFileKey                  key = "ASDF_CONFIG_FILE"
	defaultToolVersionsFilenameKey key = "ASDF_DEFAULT_TOOL_VERSIONS_FILENAME"
)

func (k key) def() string {
	return string(k) + "=" + defs()[k]
}

func (k key) str() string {
	return string(k)
}

func keys() []key {
	vals := defs()

	var keys = make([]key, 0, len(vals))
	for k := range vals {
		keys = append(keys, k)
	}

	return keys
}

func defs() map[key]string {
	return map[key]string{
		dirKey:                         "/home/user/.asdf",
		dataDirKey:                     "/home/user/.asdf",
		configFileKey:                  "/home/user/.asdfrc",
		defaultToolVersionsFilenameKey: ".tool-versions",
	}
}
