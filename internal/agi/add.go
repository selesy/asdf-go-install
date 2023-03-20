package agi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

	if len(args) != expectedArgsCount {
		p.env.log.Error("Expected 2 arguments")

		return ErrExitCodeBadArgumentCount
	}

	var (
		cmd  = args[0]
		name = args[1]
		pkg  = args[2]
	)

	// TODO: verify plugin name
	// TODO: verify plugin url is a Go package main

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
	trgtDir := strings.TrimSuffix(cmd, filepath.Join("lib", "commands", "command-add.bash"))
	if err := p.makeSymLinks(
		trgtDir,
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
		return fmt.Errorf("%w - %w - %s", ErrExitCodeCommandFailure, err, name)
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

func (p *plugin) makeSymLinks(trgtDir string, links ...string) error {
	trgt := filepath.Join(trgtDir, "bin", "asdf-go-install")

	for _, link := range links {
		err := os.Symlink(trgt, link)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w - %w - failed to create symlink %s", ErrExitCodeCommandFailure, err, trgt)
		}

		check, err := filepath.EvalSymlinks(link)
		if err != nil {
			return fmt.Errorf("%w - %w - failed to check existing symlink - %s", ErrExitCodeCommandFailure, err, link)
		}

		if check != trgt {
			return fmt.Errorf("%w - existing symlink already exists but does no equal target - %s != %s", ErrExitCodeCommandFailure, link, trgt)
		}

		p.env.log.Info("symlink already exists - skipping ", trgt)
	}

	return nil
}
