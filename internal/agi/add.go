package agi

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (p *plugin) Add(args []string) ExitCode {
	const (
		expectedArgsCount = 3
	)

	if len(args) != expectedArgsCount {
		p.env.log.Error("the Add() command expects 3 arguments")

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
	if err := p.MkDir(installDir); err != nil {
		p.env.log.Error(err)

		return ErrExitCodeCommandFailure
	}

	binDir := filepath.Join(installDir, "bin")
	if err := p.MkDir(binDir); err != nil {
		p.env.log.Error(err)

		return ErrExitCodeCommandFailure
	}

	// Write symlinks for list-all, download, install and help
	trgtDir := strings.TrimSuffix(cmd, filepath.Join("lib", "commands", "command-add.bash"))
	if err := p.Symlinks(
		trgtDir,
		filepath.Join(binDir, "download"),
		filepath.Join(binDir, "install"),
		filepath.Join(binDir, "list-all"),
	); err != nil {
		p.env.log.Error("failed to write symlink - %w", err)

		return ErrExitCodeCommandFailure
	}

	// Write config file.
	cfg := &Config{
		Name:    name,
		Package: pkg,
	}

	if err := cfg.Write(installDir); err != nil {
		return ErrExitCodeCommandFailure
	}

	p.env.log.Debug("Wrote config file to ", filepath.Join(installDir, ConfigFilename))

	// TODO: Write README/help file

	return ExitCodeOK
}

// MkDir gently creates the directory specified by path.  No error is
// returned if the path already exists and the path leads to a directory.
func (p *plugin) MkDir(path string) error {
	const directoryPermissions = 0o775

	err := os.Mkdir(path, directoryPermissions)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return errors.Join(ErrExitCodeCommandFailure, err)
	}

	if err == nil {
		p.env.log.Info("Created directory: ", path)

		return nil
	}

	info, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w - %w - %s", ErrExitCodeCommandFailure, err, path)
	}

	if err == nil && !info.IsDir() {
		return fmt.Errorf("%w - exists but is not a directory - %s", ErrExitCodeCommandFailure, path)
	}

	p.env.log.Debug("directory already exists - skipping ", path)

	return nil
}

// Symlinks gently creates one or more symbolic links all pointing to
// the same target.  If the path already exists and is a symlink to
// the specified target, no error is returned.
func (p *plugin) Symlinks(trgtDir string, links ...string) error {
	trgt := filepath.Join(trgtDir, "bin", "asdf-go-install")

	for _, link := range links {
		err := os.Symlink(trgt, link)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w - %w - failed to create symlink %s", ErrExitCodeCommandFailure, err, trgt)
		}

		if err == nil {
			p.env.log.Info("created symlink ", link)

			return nil
		}

		check, err := filepath.EvalSymlinks(link)
		if err != nil {
			return fmt.Errorf("%w - %w - failed to check existing symlink - %s", ErrExitCodeCommandFailure, err, link)
		}

		if check != trgt {
			return fmt.Errorf("%w - symlink already exists but does not equal target - %s != %s", ErrExitCodeCommandFailure, link, trgt)
		}

		p.env.log.Info("symlink already exists - skipping ", trgt)
	}

	return nil
}
