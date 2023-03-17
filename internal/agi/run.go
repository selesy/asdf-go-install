package agi

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	// ErrPanicNilPlugin provides an error message when the developer
	// accidentally calls Env.Execute without providing a Plugin.
	ErrPanicNilPlugin = errors.New("nil plugin - programming error calling Env.Execute()")

	// ErrPanicUnknownError happens if an exit code (with error) does not
	// provide an error message.
	ErrPanicUnknownError = errors.New("exit code without message - programming error when calling ErrExitCode.Error()")
)

type (
	// ExitCode is a POSIX return code emitted when the plugin exits without
	// an error.
	ExitCode int

	// ExitCodeError is an error that indicates the plugin exited with an
	// error.
	ExitCodeError = ExitCode
)

const (
	// ExitCodeOK returns an ExitCode that indicates the plugin finished
	// its execution successfully.
	ExitCodeOK ExitCode = 0

	// ErrExitCoddeNotImplemented is an ErrExitCode that indicates part of
	// the called code has not yet been implemented.
	ErrExitCodeNotImplemented ExitCodeError = iota + 1

	// ErrExitCodeNoCommand is an ExitCodeError that indicates no argument
	// was passed to the command parser.
	ErrExitCodeNoCommand

	// ErrExitCodeFoundArguments is an ExitCodeError that indicates extra
	// arguments were included on the command line.
	ErrExitCodeFoundArguments

	// ErrExitCodeUnknownCommand is an ExitCodeError that indicates an
	// unknonwn command was parsed from the arguments.
	ErrExitCodeUnknownCommand

	// ErrExitCodeEnvVarFailure is an ExitCodeError that indicates there
	// was a problem parsing the environment variables needed for a
	// specific plugin command.
	//
	// See https://asdf-vm.com/plugins/create.html#environment-variables
	ErrExitCodeEnvVarFailure
)

const (
	// MsgErrEnvVarParseFailed is the message for ErrExitCodeEnvVarFailure.
	MsgErrEnvVarParseFailed = "parsing the envvars for the called plugin function failed"

	// MsgErrFoundArgument is the message for ErrExitCodeFoundArguments.
	MsgErrFoundArgument = "the plugin executable should be called with no additional arguments"

	// MsgErrNoCommand is the message for ErrExitCodeNoCommand.
	MsgErrNoCommand = "the plugin executable should be called with the plugin function name"

	// MsgErrNotImplemented is the message for ErrExitCodeNotImplemented.
	MsgErrNotImplemented = "the called plugin function  is not implemented"

	// MsgErrUnknownCommand is the message for ErrExitCodeUnknownCommand.
	MsgErrUnknownCommand = "the called plugin function is not known"
)

// Error returns the stringified message for the underlying error.
//
// Implements: error.
func (e ExitCodeError) Error() string {
	var msg string

	msg, ok := map[ExitCode]string{ //nolint:varnamelen
		ErrExitCodeEnvVarFailure:  MsgErrEnvVarParseFailed,
		ErrExitCodeFoundArguments: MsgErrFoundArgument,
		ErrExitCodeNoCommand:      MsgErrNoCommand,
		ErrExitCodeNotImplemented: MsgErrNotImplemented,
		ErrExitCodeUnknownCommand: MsgErrUnknownCommand,
	}[e]
	if !ok {
		panic(ErrPanicUnknownError)
	}

	return msg
}

// PluginFunc is the signature of all Plugin functions which takes no
// arguments and returns an ExitCode.  All Plugin functions match this
// signature - instead of arguments, each command retrieves its input
// paraameters from environment variables.
type PluginFunc func() ExitCode

type pluginFunc func(Plugin) ExitCode

var (
	_ pluginFunc = Plugin.Download
	_ pluginFunc = Plugin.HelpOverview
	_ pluginFunc = Plugin.Install
	_ pluginFunc = Plugin.ListAll
)

// Plugin implements the functions required to implement an asdf plugin.
type Plugin interface {
	Download() ExitCode
	HelpOverview() ExitCode
	Install() ExitCode
	ListAll() ExitCode
}

// EnvOption is a function that alters the Env passed as an argument.
type EnvOption func(*Env)

// WithStdout provides an EnvOption that sets the output writer within
// the Env.
func WithStdout(stdout io.Writer) EnvOption {
	return func(e *Env) {
		e.out = stdout
	}
}

// WithStderr provides an EnvOption that sets the error writer within
// the Env.  All logging is passed to the error writer regardless of
// the logging level.
func WithStderr(stderr io.Writer) EnvOption {
	return func(e *Env) {
		e.err = stderr
		e.log.SetOutput(stderr)
	}
}

// WithLogFormatter is an EnvOption that sets the logging system's
// formatter.
func WithLogFormatter(formatter logrus.Formatter) EnvOption {
	return func(e *Env) {
		e.log.Formatter = formatter
	}
}

// Env is the execution context for an asdf Plugin.
type Env struct {
	out io.Writer
	err io.Writer
	log *logrus.Logger
}

// NewEnv returns a default execution context if no EnvOptions are
// provided.  Otherwise, the default Env is altered by the provided
// options.
func NewEnv(opts ...EnvOption) *Env {
	fmtr := &logrus.JSONFormatter{} //nolint:exhaustivestruct,exhaustruct
	fmtr.TimestampFormat = "2006-01-02T15:04:05Z"

	env := &Env{
		out: os.Stdout,
		err: os.Stderr,
		log: logrus.New(),
	}

	env.log.SetOutput(env.err)
	env.log.SetFormatter(fmtr)

	_, ok := os.LookupEnv("ASDF_DEBUG")
	if ok {
		env.log.SetLevel(logrus.DebugLevel)
	}

	for _, opt := range opts {
		opt(env)
	}

	return env
}

// Execute parses the asdf script name from the first argument and runs
// the associated plugin function with in the provided environment.
func (e *Env) Execute(plugin Plugin, args []string) ExitCode {
	e.log.Trace("Execute()")

	if plugin == nil {
		panic(ErrPanicNilPlugin)
	}

	e.log.Debug("Arguments: ", args)

	if len(args) < 1 {
		return ErrExitCodeNoCommand
	}

	if len(args) > 1 {
		return ErrExitCodeFoundArguments
	}

	_, fileName := filepath.Split(args[0])

	e.log.Debug("Executable name: ", fileName)

	fn, ok := map[string]PluginFunc{ //nolint:varnamelen
		"download": plugin.Download,
		"help":     plugin.HelpOverview,
		"install":  plugin.Install,
		"list-all": plugin.ListAll,
	}[fileName]
	if !ok {
		return ErrExitCodeUnknownCommand
	}

	return fn()
}

var _ Plugin = (*plugin)(nil)

type plugin struct {
	env *Env
}

// Main executes the default plugin using the default execution context.
func Main() int {
	env := NewEnv()

	plugin := &plugin{
		env: env,
	}

	return int(env.Execute(plugin, os.Args))
}

// HelpOverview provides a general description of the Plugin.
//
// See https://asdf-vm.com/plugins/create.html#bin-help-scripts
func (p *plugin) HelpOverview() ExitCode {
	return ErrExitCodeNotImplemented
}
