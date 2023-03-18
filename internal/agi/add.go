package agi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type Config struct {
	Name    string `json:"name"`
	Package string `json:"package"`
}

func (p *plugin) Add(args []string) ExitCode {
	const (
		configFilePermissions = 0o555
		expectedArgsCount     = 3
	)

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

	asdfDir, ok := os.LookupEnv("ASDF_DIR")
	if !ok {
		p.env.log.Error("Missing the ASDF_DIR environment variable")

		return ErrExitCodeEnvVarFailure
	}

	p.env.log.Debug("ASDF directory: ", asdfDir)

	pluginDir := filepath.Join(asdfDir, "plugins")

	installDir := filepath.Join(pluginDir, name)
	if err := p.makeDirectoryIfNotExists(installDir); err != nil {
		p.env.log.Error(err)

		return ErrExitCodeCommandFailure
	}

	binDir := filepath.Join(installDir, "bin")
	if err := p.makeDirectoryIfNotExists(binDir); err != nil {
		p.env.log.Error(err)

		return ErrExitCodeCommandFailure
	}

	// TODO: Write symlinks for list-all, download, install and help
	if err := p.makeSymLinks(
		pluginDir,
		filepath.Join(binDir, "download"),
		filepath.Join(binDir, "install"),
		filepath.Join(binDir, "list-all"),
	); err != nil {
		p.env.log.Error("failed to write symlink - %w", err)

		return ErrExitCodeCommandFailure
	}

	// TODO: Write README/help file

	// TODO: Write config file.
	cfg := &Config{
		Name:    name,
		Package: pkg,
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		p.env.log.Error("failed to marshal Config to JSON")

		return ErrExitCodeCommandFailure
	}

	configPath := filepath.Join(installDir, ".config")
	if err := os.WriteFile(configPath, data, configFilePermissions); err != nil {
		p.env.log.Error("failed to write .config file")

		return ErrExitCodeCommandFailure
	}

	p.env.log.Debug("Wrote config file to ", configPath)

	return ErrExitCodeNotImplemented
}

func (p *plugin) makeDirectoryIfNotExists(name string) error {
	const directoryPermissions = 0o775

	fi, err := os.Stat(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w - %w - %s", ErrExitCodeEnvVarFailure, err, name)
	}

	if err == nil && !fi.IsDir() {
		return fmt.Errorf("%w - exists but is not a directory - %s", ErrExitCodeCommandFailure, name)
	}

	if err == nil && fi.IsDir() {
		p.env.log.Debug("directory already exists - skipping ", name)

		return nil
	}

	if err := os.Mkdir(name, directoryPermissions); err != nil {
		return fmt.Errorf("%w - failed to create directory - %s", ErrExitCodeCommandFailure, name)
	}

	p.env.log.Debugf("created directory - %s", name)

	return nil
}

func (p plugin) makeSymLinks(asdfDir string, links ...string) error {
	target := filepath.Join(asdfDir, "go-install", "bin", "asdf-go-install")

	for _, link := range links {
		if err := os.Symlink(target, link); err != nil {
			return err
		}

		p.env.log.Debugf("created symlink - %s", link)
	}

	return nil
}
