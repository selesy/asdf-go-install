package agi

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	ConfigFilename        = ".config.json"
	ConfigFilePermissions = 0o775
)

type Config struct {
	Name    string `json:"name"`
	Package string `json:"package"`
}

var (
	ErrFailedToMarshalConfig   = errors.New("failed to marshall configuration file")
	ErrFailedToReadConfig      = errors.New("failed to read configuration file")
	ErrFailedToUnmarshalConfig = errors.New("failed to unmarshal configuration file")
	ErrFailedToWriteConfig     = errors.New("failed to write configuration file")
)

func (c *Config) Read(path string) error {
	configPath := filepath.Join(path, ConfigFilename)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return errors.Join(ErrFailedToReadConfig, err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return errors.Join(ErrFailedToUnmarshalConfig, err)
	}

	return nil
}

func (c *Config) Write(path string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return errors.Join(ErrFailedToMarshalConfig, err)
	}

	configPath := filepath.Join(path, ConfigFilename)
	if err := os.WriteFile(configPath, data, ConfigFilePermissions); err != nil {
		return errors.Join(ErrFailedToWriteConfig, err)
	}

	return nil
}
