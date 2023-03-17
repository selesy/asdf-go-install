package agi_test

import (
	"bytes"
	"testing"

	"github.com/selesy/asdf-go-install/internal/agi"
	"github.com/selesy/asdf-go-install/internal/agi/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

type testRecorder struct {
	out *bytes.Buffer
	err *bytes.Buffer
}

func (r *testRecorder) show(t *testing.T, path string) {
	t.Helper()

	golden.Assert(t, r.out.String(), path+"_out.txt")
	golden.Assert(t, r.err.String(), path+"_err.log")
}

func newTestEnv(t *testing.T) (*agi.Env, *testRecorder) {
	t.Helper()

	recorder := &testRecorder{
		out: &bytes.Buffer{},
		err: &bytes.Buffer{},
	}

	fmtr := &logrus.JSONFormatter{} //nolint:exhaustruct,exhaustivestruct
	fmtr.TimestampFormat = "Timestamp"

	env := agi.NewEnv(
		agi.WithStdout(recorder.out),
		agi.WithStderr(recorder.err),
		agi.WithLogFormatter(fmtr),
	)

	return env, recorder
}

func TestExitCodeError(t *testing.T) {
	t.Parallel()

	t.Run("err is equal to underlying message", func(t *testing.T) {
		t.Parallel()

		require.EqualError(t, agi.ErrExitCodeFoundArguments, agi.MsgErrFoundArgument)
	})

	t.Run("panis with unknown exit code", func(t *testing.T) {
		t.Parallel()

		require.Panics(t, func() {
			_ = agi.ExitCodeError(-1).Error()
		})
		require.PanicsWithError(t, agi.ErrPanicUnknownError.Error(), func() {
			_ = agi.ExitCodeError(-1).Error()
		})
	})
}

func TestEnv_Run(t *testing.T) { //nolint:tparallel
	t.Setenv("ASDF_DEBUG", "")

	t.Run("panics with nil plugin", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		require.Panics(t, func() {
			_ = env.Execute(nil, []string{"install"})
		})

		rec.show(t, "panics_with_nil_plugin")
	})

	t.Run("errors without command argument", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		mock := mocks.NewPlugin(t)
		require.ErrorIs(t, env.Execute(mock, []string{}), agi.ErrExitCodeNoCommand)

		rec.show(t, "errors_without_command_argument")
	})

	t.Run("errors with extra argument(s)", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		mock := mocks.NewPlugin(t)
		require.ErrorIs(t, env.Execute(mock, []string{"install", "argument"}), agi.ErrExitCodeFoundArguments)

		rec.show(t, "errors_with_extra_arguments")
	})

	t.Run("errors on unknown command", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		mock := mocks.NewPlugin(t)
		require.ErrorIs(t, env.Execute(mock, []string{"unknown"}), agi.ErrExitCodeUnknownCommand)

		rec.show(t, "errors_on_unknown_command")
	})

	t.Run("executes a known command", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		mock := mocks.NewPlugin(t)
		mock.EXPECT().Install().Return(0)
		_ = env.Execute(mock, []string{"install"})

		rec.show(t, "executes_a_known_command")
	})
}
