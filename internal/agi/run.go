package agi

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

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

	// ErrExitCodeInvalidArgument is an ExitCodeError that indicates one
	// of the command's arguments is invalid.
	ErrExitCodeInvalidArgument

	// ErrWrongNumberOfArguments is an ExitCodeError that indicates the
	// wrong number of arguments was found.
	ErrExitCodeBadArgumentCount

	// ErrExitCodeUnknownCommand is an ExitCodeError that indicates an
	// unknonwn command was parsed from the arguments.
	ErrExitCodeUnknownCommand

	// ErrExitCodeEnvVarFailure is an ExitCodeError that indicates there
	// was a problem parsing the environment variables needed for a
	// specific plugin command.
	//
	// See https://asdf-vm.com/plugins/create.html#environment-variables
	ErrExitCodeEnvVarFailure

	// ErrExitCodeCommandFailure is an ExitCodeError that indicates that
	// a problem occurred while running a plugin command.
	ErrExitCodeCommandFailure
)

const (
	// MsgErrBadArgumentCount is the message for ErrExitCodeBadArgumentCount.
	MsgErrBadArgumentCount = "found the wrong number of arguments"

	// MsgErrCommandFailure is the message for ErrExitCodeCommandFailure.
	MsgErrCommandFailure = "failure while executing plugin command"

	// MsgErrEnvVarParseFailed is the message for ErrExitCodeEnvVarFailure.
	MsgErrEnvVarParseFailed = "parsing the envvars for the called plugin function failed"

	// MsgErrFoundArgument is the message for ErrExitCodeFoundArguments.
	MsgErrFoundArgument = "the plugin executable should be called with no additional arguments"

	// MsgErrInvalidArgument is the message for ErrExitCodeInvalidArgument.
	MsgErrInvalidArgument = "an invalid argument was found"

	// MsgErrNoCommand is the message for ErrExitCodeNoCommand.
	MsgErrNoCommand = "the plugin executable should be called with the plugin function name"

	// MsgErrNotImplemented is the message for ErrExitCodeNotImplemented.
	MsgErrNotImplemented = "the called plugin function is not implemented"

	// MsgErrUnknownCommand is the message for ErrExitCodeUnknownCommand.
	MsgErrUnknownCommand = "the called plugin function is not known"
)

// Error returns the stringified message for the underlying error.
//
// Implements: error.
func (e ExitCodeError) Error() string {
	var msg string

	msg, ok := map[ExitCodeError]string{ //nolint:varnamelen
		ErrExitCodeBadArgumentCount: MsgErrBadArgumentCount,
		ErrExitCodeCommandFailure:   MsgErrCommandFailure,
		ErrExitCodeEnvVarFailure:    MsgErrEnvVarParseFailed,
		ErrExitCodeFoundArguments:   MsgErrFoundArgument,
		ErrExitCodeInvalidArgument:  MsgErrInvalidArgument,
		ErrExitCodeNoCommand:        MsgErrNoCommand,
		ErrExitCodeNotImplemented:   MsgErrNotImplemented,
		ErrExitCodeUnknownCommand:   MsgErrUnknownCommand,
	}[e]
	if !ok {
		panic(ErrPanicUnknownError)
	}

	return msg
}

type (
	ExtensionFunc func([]string) ExitCode
	// PluginFunc is the signature of all Plugin functions which takes no
	// arguments and returns an ExitCode.  All Plugin functions match this
	// signature - instead of arguments, each command retrieves its input
	// paraameters from environment variables.
	PluginFunc func() ExitCode
)

type (
	extensionFunc func(Plugin, []string) ExitCode
	pluginFunc    func(Plugin) ExitCode
)

var (
	_ pluginFunc    = Plugin.Download
	_ pluginFunc    = Plugin.HelpOverview
	_ pluginFunc    = Plugin.Install
	_ pluginFunc    = Plugin.ListAll
	_ extensionFunc = Plugin.Add
)

// Plugin implements the functions required to implement an asdf plugin.
type Plugin interface {
	Download() ExitCode
	HelpOverview() ExitCode
	Install() ExitCode
	ListAll() ExitCode
	Add([]string) ExitCode
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

	if e.log.Level == logrus.DebugLevel {
		cwd, err := os.Getwd()
		if err != nil {
			e.log.Error("failed to retrieve current working directory")

			return ErrExitCodeEnvVarFailure
		}

		e.log.Debug("Current working directory: ", cwd)

		for i, v := range args {
			e.log.Debugf("Argument (%d): %s", i, v)
		}

		for _, v := range os.Environ() {
			if strings.HasPrefix(v, "ASDF") {
				e.log.Debug("Environment variable: ", v)
			}
		}
	}

	if plugin == nil {
		panic(ErrPanicNilPlugin)
	}

	if len(args) < 1 {
		return ErrExitCodeNoCommand
	}

	_, fileName := filepath.Split(args[0])

	e.log.Debug("Executable name: ", fileName)

	fn, ok := map[string]PluginFunc{ //nolint:varnamelen
		"download":         verifyNoArguments(plugin.Download, args),
		"help":             verifyNoArguments(plugin.HelpOverview, args),
		"install":          verifyNoArguments(plugin.Install, args),
		"list-all":         verifyNoArguments(plugin.ListAll, args),
		"command-add.bash": func() ExitCode { return plugin.Add(args) },
	}[fileName]
	if !ok {
		e.log.Error(ErrExitCodeUnknownCommand, " - ", fileName)

		return ErrExitCodeUnknownCommand
	}

	return fn()
}

var _ Plugin = (*plugin)(nil)

type plugin struct {
	env *Env
}

func NewPlugin(env *Env) *plugin { //nolint:golint
	return &plugin{
		env: env,
	}
}

// Main executes the default plugin using the default execution context.
func Main() int {
	env := NewEnv()
	plugin := NewPlugin(env)

	return int(env.Execute(plugin, os.Args))
}

// HelpOverview provides a general description of the Plugin.
//
// See https://asdf-vm.com/plugins/create.html#bin-help-scripts
func (p *plugin) HelpOverview() ExitCode {
	return ErrExitCodeNotImplemented
}

func verifyNoArguments(fn PluginFunc, args []string) PluginFunc { //nolint:varnamelen
	return func() ExitCode {
		if len(args) > 1 {
			return ErrExitCodeFoundArguments
		}

		return fn()
	}
}
