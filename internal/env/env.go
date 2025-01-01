// Package env parses the environment variables that control the behavior
// of the asdf-go-install plugin.
package env

import (
	"encoding"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"reflect"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/caarlos0/env/v10"
	"github.com/go-git/go-git/v5/plumbing/hash"
	"github.com/go-playground/validator/v10"
)

// Env contains the values used by the ASDF "scripts" and "extension
// commands" after they're parsed from passed environment variables.
type Env struct {
	agiVar  agiVar
	asdfVar asdfVar
}

// New parses the (relevant) environment variables available from the OS
// and creates an immutable instance of values that are used during
// plugin execution.
func New(log *slog.Logger, environ []string) (*Env, error) {
	log = log.WithGroup("env")

	for _, v := range environ {
		if !strings.HasPrefix(v, "ASDF") && !strings.HasPrefix(v, "AGI") {
			continue
		}

		key, val, _ := strings.Cut(v, "=")
		log.Debug(
			"candidate environment variable",
			slog.String("key", key),
			slog.String("value", val),
		)
	}

	var asdfVar asdfVar

	if err := env.ParseWithOptions(&asdfVar, env.Options{
		Prefix:                "ASDF_",
		Environment:           env.ToMap(environ),
		UseFieldNameByDefault: true,
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf((*semver.Version)(nil)): parseSemVer,
			reflect.TypeOf((*url.URL)(nil)):        parseURL,
			reflect.TypeOf((*hash.Hash)(nil)):      parseGitHash, // TODO: not current used
		},
	}); err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(asdfVar); err != nil {
		return nil, err
	}

	var agiVar agiVar

	if err := env.ParseWithOptions(&agiVar, env.Options{
		Prefix:                "AGI_",
		Environment:           env.ToMap(environ),
		UseFieldNameByDefault: true,
		FuncMap:               map[reflect.Type]env.ParserFunc{},
	}); err != nil {
		return nil, err
	}

	e := &Env{
		agiVar:  agiVar,
		asdfVar: asdfVar,
	}

	resolvedEnvironment := func(attr slog.Attr) {
		log.Debug("Resolved environment variable", attr)
	}

	resolvedEnvironment(slog.String("Dir", e.Dir()))
	resolvedEnvironment(slog.String("ConfigFile", e.ConfigFile()))
	resolvedEnvironment(slog.String("DataDir", e.DataDir()))
	resolvedEnvironment(slog.String("DefaultToolVersionsFilename", e.DefaultToolVersionsFilename()))
	resolvedEnvironment(slog.Any("InstallType", e.InstallType()))
	// *semver.MarshalText doesn't allow nil as a receiver
	resolvedEnvironment(slog.Any("InstallVersion", fmt.Sprintf("%v", e.InstallVersion())))
	resolvedEnvironment(slog.String("InstallPath", e.InstallPath()))
	resolvedEnvironment(slog.Int("Concurrency", e.Concurrency()))
	resolvedEnvironment(slog.String("DownloadPath", e.DownloadPath()))
	resolvedEnvironment(slog.String("PluginPath", e.PluginPath()))
	resolvedEnvironment(slog.Any("PluginSourceURL", e.PluginSourceURL()))
	resolvedEnvironment(slog.String("PluginPrevRef", e.PluginPrevRef()))
	resolvedEnvironment(slog.String("PluginPostRef", e.PluginPostRef()))
	resolvedEnvironment(slog.String("CmdFile", e.CmdFile()))

	return e, nil
}

// CmdFile resolves to the full path of the file being executed.
func (e *Env) CmdFile() string {
	return e.asdfVar.CmdFile
}

// Concurrency returns the number of cores to use when compiling the
// source code.
//
// Useful for setting make -j
func (e *Env) Concurrency() int {
	return e.asdfVar.Concurrency
}

// ConfigFile returns the full path of the asdf configuration file.
func (e *Env) ConfigFile() string {
	return e.asdfVar.ConfigFile
}

// DataDir returns the full path of the asdf data directory.
func (e *Env) DataDir() string {
	return e.asdfVar.DataDir
}

// DefaultToolVersionsFilename returns the expected asdf configuration
// file name.
func (e *Env) DefaultToolVersionsFilename() string {
	return e.asdfVar.DefaultToolVersionsFilename
}

// Dir returns the full root path of the asdf installation.
func (e *Env) Dir() string {
	return e.asdfVar.Dir
}

// DownloadPath returns the path to the source code or binary that was
// downloaded by bin/download.
func (e *Env) DownloadPath() string {
	return e.asdfVar.DownloadPath
}

// InstallType returns either InstallTypeVersion or InstallTypeRef.
func (e *Env) InstallType() InstallType {
	return e.asdfVar.InstallType
}

// InstallPath returns the path to where the tool should, or has been,
// installed.
func (e *Env) InstallPath() string {
	return e.asdfVar.InstallPath
}

// InstallVersion returns the full version number or Git reference
// depending on the value returned by InstallType().
func (e *Env) InstallVersion() *semver.Version {
	return e.asdfVar.InstallVersion
}

// LogFormat returns the format of the logger's output.
func (e *Env) LogFormat() LogFormat {
	return e.agiVar.LogFormat
}

// LogLevel returns the logger's Level.
func (e *Env) LogLevel() slog.Level {
	return e.agiVar.LogLevel
}

// LogOutput returns the destination for the logger's records.
func (e *Env) LogOutput() string {
	return e.agiVar.LogOutput
}

// LogSource indicates whether the log records should include
// file and line number information.
func (e *Env) LogSource() bool {
	return e.agiVar.LogSource
}

// PluginPath returns the path where the plugin was installed.
func (e *Env) PluginPath() string {
	return e.asdfVar.PluginPath
}

// PluginPostRef returns the updated git-ref of the plugin's Git
// repository.
func (e *Env) PluginPostRef() string {
	return e.asdfVar.PluginPostRef
}

// PluginPrevRef returns the previous git-ref of the plugin's Git
// repository.
func (e *Env) PluginPrevRef() string {
	return e.asdfVar.PluginPrevRef
}

// PluginSourceURL returns the URL of the plugin's git repository.
func (e *Env) PluginSourceURL() *url.URL {
	return e.asdfVar.PluginSourceURL
}

func parseGitHash(s string) (any, error) {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	if len(decoded) != hash.Size {
		return nil, errors.New("blah")
	}

	return nil, nil
}

func parseSemVer(s string) (any, error) {
	return semver.StrictNewVersion(s)
}

func parseURL(s string) (any, error) {
	return url.Parse(s)
}

type agiVar struct {
	LogFormat LogFormat
	LogLevel  slog.Level
	LogOutput string
	LogSource bool
}

var _ encoding.TextUnmarshaler = (*LogFormat)(nil)

// LogFormat represents the desired output formatting of each log's
// record.
type LogFormat int

const (
	// LogFormatColorized indicates that the log records should be text
	// formatted and should be made more visually appealing using ANSI
	// colors.
	LogFormatColorized LogFormat = iota + 1
	// LogFormatJSON indicates that the log records should be formatted
	// as JSON.
	LogFormatJSON
	// LogFormatText indicates that the log records should be formatted
	// as plain (ASCII) text.
	LogFormatText
)

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *LogFormat) UnmarshalText(p []byte) error {
	switch strings.ToLower(string(p)) {
	case "colorized":
		*f = LogFormatColorized

		return nil
	case "json":
		*f = LogFormatJSON

		return nil
	case "text":
		*f = LogFormatText

		return nil
	default:
		return fmt.Errorf("%w: from \"%s\"", ErrInvalidLogFormat, string(p))
	}
}

type asdfVar struct {
	// Set by asdf.sh sourced into ~/.bashrc (e.g)
	Dir string `validate:"required"`

	// Set during any asdf execution
	ConfigFile                  string `validate:"required"`
	DataDir                     string `validate:"required"`
	DefaultToolVersionsFilename string `validate:"required"`

	// Set (or not) as described at https://asdf-vm.com/plugins/create.html#environment-variables-overview
	// Note that not all the environment variables below are provided
	// for every script.  See the individual script documentation for
	// more details.
	InstallType     InstallType
	InstallVersion  *semver.Version
	InstallPath     string
	Concurrency     int
	DownloadPath    string
	PluginPath      string
	PluginSourceURL *url.URL
	PluginPrevRef   string // TODO: type this more strongly?
	PluginPostRef   string // TODO: type this more strongly?
	CmdFile         string

	// Proposed logging flags: https://github.com/asdf-vm/asdf/issues/702#issuecomment-814234517
	Verbose       bool
	VerboseOutput string
}

const (
	installTypeVersion = "version"
	installTypeRef     = "ref"
)

var (
	_ encoding.TextMarshaler   = (*InstallType)(nil)
	_ encoding.TextUnmarshaler = (*InstallType)(nil)
)

// InstallType represents the type of Git reference that the plugin will
// be installed from.
type InstallType int

const (
	// InstallTypeVersion indicates that the plugin will be installed from
	// a Git tag reference structured as a Go version number.
	InstallTypeVersion InstallType = iota + 1
	// InstallTypeRef indicates that the plugin will be installed from a
	// Git reference which might include commits or branches.
	InstallTypeRef
)

// MarshalText implements encoding.TextMarshaler.
func (it *InstallType) MarshalText() ([]byte, error) {
	switch *it {
	case InstallTypeVersion:
		return []byte(installTypeVersion), nil
	case InstallTypeRef:
		return []byte(installTypeRef), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrMarshalFailed, it)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (it *InstallType) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case installTypeVersion:
		*it = InstallTypeVersion
	case installTypeRef:
		*it = InstallTypeRef
	default:
		return fmt.Errorf("%w: %s", ErrUnmarshalFailed, text)
	}

	return nil
}

// String returns the text representation of valid InstallType values
// or an empty string for invalid values.
func (it *InstallType) String() string {
	text, err := it.MarshalText()
	if err != nil {
		return ""
	}

	return string(text)
}
