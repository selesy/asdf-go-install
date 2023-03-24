package agi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/selesy/asdf-go-install/internal/agi"
	"github.com/stretchr/testify/require"
)

func tempDir(t *testing.T, name string) string {
	t.Helper()

	const directoryPermissions = 0o775

	path := filepath.Join(os.TempDir(), name)
	require.NoError(t, os.Mkdir(path, directoryPermissions))

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(path))
	})

	return path
}

func unsetEnvVar(t *testing.T, key string) {
	t.Helper()

	val, ok := os.LookupEnv(key)
	if ok {
		t.Cleanup(func() {
			require.NoError(t, os.Setenv(key, val))
		})
	}

	require.NoError(t, os.Unsetenv(key))
}

func newTestPlugin(t *testing.T) (agi.Plugin, *testRecorder) {
	t.Helper()

	env, rec := newTestEnv(t)
	plugin := agi.NewPlugin(env)

	return plugin, rec
}

func TestPlugin_Add(t *testing.T) {
	t.Parallel()

	t.Run("errors if ASDF_DIR environment variable is not set", func(t *testing.T) {
		plugin, rec := newTestPlugin(t)

		unsetEnvVar(t, "ASDF_DIR")

		exp := "parsing the envvars for the called plugin function failed"
		require.EqualError(t, plugin.Add([]string{"1", "2", "3"}), exp)
		rec.show(t, "errors_if_asdfdir_environment_variable_is_missing")
	})

	t.Run("errors if there are not three arguments", func(t *testing.T) {
		t.Parallel()

		exp := "found the wrong number of arguments"

		t.Run("", func(t *testing.T) {
			t.Parallel()

			plugin, rec := newTestPlugin(t)
			require.EqualError(t, plugin.Add([]string{"one", "two"}), exp)
			rec.show(t, "errors_if_there_are_less_than_three_arguments")
		})

		t.Run("", func(t *testing.T) {
			t.Parallel()

			plugin, rec := newTestPlugin(t)
			require.EqualError(t, plugin.Add([]string{"one", "two", "three", "four"}), exp)
			rec.show(t, "errors_if_there_are_more_than_three_arguments")
		})
	})
}

func TestPlugin_MkDir(t *testing.T) {
	t.Parallel()

	t.Run("errors if path exists but is not a directory", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		exp := "failure while executing plugin command - exists but is not a directory - testdata/make_directory_if_not_exists_file"
		require.EqualError(t, plugin.MkDir("testdata/make_directory_if_not_exists_file"), exp)
		rec.show(t, "errors_if_path_exists_but_is_not_a_directory")
	})

	t.Run("errors if path exists but is not a directory", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		exp := "failure while executing plugin command\nmkdir testdata/does_not_exist/does_not_exist: no such file or directory"
		require.EqualError(t, plugin.MkDir("testdata/does_not_exist/does_not_exist"), exp)
		rec.show(t, "errors_if_directory_creation_fails")
	})

	t.Run("succeeds if path exists and is directory", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		require.NoError(t, plugin.MkDir("testdata/make_directory_if_not_exists_folder"))
		rec.show(t, "succeeds_if_path_exists_and_is_directory")
	})

	t.Run("succeeds if directory is created", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(tempDir(t, "succeeds_if_directory_is_created"), "new_directory")

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		require.NoError(t, plugin.MkDir(path))
		rec.show(t, "succeeds_if_directory_is_created")
	})
}

func TestPlugin_Symlinks(t *testing.T) {
	t.Parallel()

	t.Run("errors if path exists but is not target", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		exp := "failure while executing plugin command - symlink already exists but does not equal target - ./testdata/symlink_to_target_false != testdata/bin/asdf-go-install"
		require.EqualError(t, plugin.Symlinks("./testdata", "./testdata/symlink_to_target_false"), exp)
		rec.show(t, "errors_if_path_exists_but_is_not_target")
	})

	t.Run("succeeds if symlink already exists to correct target", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		require.NoError(t, plugin.Symlinks("./testdata", "./testdata/symlink_to_target_true"))
		rec.show(t, "succeeds_if_symlink_exists_to_correct_target")
	})

	t.Run("succeeds if symlink is created", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(tempDir(t, "succeeds_if_symlink_is_created"), "created_symlink")

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		require.NoError(t, plugin.Symlinks("./testdata", path))
		rec.show(t, "succeeds_if_symlnk_is_created")
	})
}
