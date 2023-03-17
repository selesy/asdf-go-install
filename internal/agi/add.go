package agi

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func (p *plugin) Add(name, url string) ExitCode {
	baseDir, ok := os.LookupEnv("ASDF_DIR")
	if !ok {
		return ErrExitCodeEnvVarFailure
	}

	pluginDir := filepath.Join(baseDir, "name")

	if err := makePluginDir(pluginDir); err != nil {
		return ErrExitCodeEnvVarFailure
	}

	return ErrExitCodeNotImplemented
}

func makePluginDir(name string) error {
	fi, err := os.Stat(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w - %w - %s", ErrExitCodeEnvVarFailure, err, name)
	}

	if fi.IsDir() {
		return fmt.Errorf("%w - exists but is not a directory - %s", ErrExitCodeEnvVarFailure, name)
	}

	if err := os.Mkdir(name, fs.FileMode(os.O_RDWR)); err != nil {
		return fmt.Errorf("%w - failed to create directory - %s", ErrExitCodeEnvVarFailure, name)
	}

	return ErrExitCodeNotImplemented
}
