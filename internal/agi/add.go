package agi

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

func (p *plugin) Add(args []string) ExitCode {
	const expectedArgsCount = 3

	var (
		name = args[1]
		pkg  = args[2]
	)

	if len(args) != expectedArgsCount {
		p.env.log.Error("Expected 2 arguments")

		return ErrExitCodeBadArgumentCount
	}

	// TODO: verify plugin name

	if _, err := url.Parse(pkg); err != nil {
		p.env.log.Error("url must be valid Go package")

		return ErrExitCodeInvalidArgument
	}

	baseDir, ok := os.LookupEnv("ASDF_DIR")
	if !ok {
		p.env.log.Error("Missing the ASDF_DIR environment variable")

		return ErrExitCodeEnvVarFailure
	}

	p.env.log.Debug("ASDF directory: ", baseDir)

	pluginDir := filepath.Join(baseDir, "plugins", name)
	if err := p.mkDirIfNotExist(pluginDir); err != nil {
		p.env.log.Error(err)

		return ErrExitCodeCommandFailure
	}

	binDir := filepath.Join(pluginDir, "bin")
	if err := p.mkDirIfNotExist(binDir); err != nil {
		p.env.log.Error(err)

		return ErrExitCodeCommandFailure
	}

	// TODO: Write symlinks for list-all, download, install and help
	if err := p.makeSymLinks(
		filepath.Join("..", "..", name, "bin", "download"),
		filepath.Join("..", "..", name, "bin", "install"),
		filepath.Join("..", "..", name, "bin", "list-all"),
	); err != nil {
		p.env.log.Error("failed to write symlink - %w", err)

		return ErrExitCodeCommandFailure
	}

	// TODO: Write README/help file
	// TODO: Write config file.

	return ErrExitCodeNotImplemented
}

func (p *plugin) mkDirIfNotExist(name string) error {
	const directoryPermissions = 0o775

	fi, err := os.Stat(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w - %w - %s", ErrExitCodeEnvVarFailure, err, name)
	}

	if err == nil && fi.IsDir() {
		return fmt.Errorf("%w - exists but is not a directory - %s", ErrExitCodeEnvVarFailure, name)
	}

	if err := os.Mkdir(name, directoryPermissions); err != nil {
		return fmt.Errorf("%w - failed to create directory - %s", ErrExitCodeEnvVarFailure, name)
	}

	p.env.log.Debugf("created directory - %s", name)

	return nil
}

func (p plugin) makeSymLinks(links ...string) error {
	for _, link := range links {
		if err := os.Symlink("../../go-install/bin/asdf-go-install", link); err != nil {
			return err
		}
	}

	return nil
}
