package agi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/selesy/asdf-go-install/internal/agi"
	"github.com/stretchr/testify/require"
)

func unsetEnvVar(t *testing.T, key string) {
	val, ok := os.LookupEnv(key)
	if ok {
		t.Cleanup(func() {
			require.NoError(t, os.Setenv(key, val))
		})
	}

	require.NoError(t, os.Unsetenv(key))
}

func TestPlugin_Add(t *testing.T) {
	t.Parallel()

	t.Run("errors if ASDF_DIR environment variable is not set", func(t *testing.T) {
		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)

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

			env, rec := newTestEnv(t)
			plugin := agi.NewPlugin(env)
			require.EqualError(t, plugin.Add([]string{"one", "two"}), exp)
			rec.show(t, "errors_if_there_are_less_than_three_arguments")
		})

		t.Run("", func(t *testing.T) {
			t.Parallel()

			env, rec := newTestEnv(t)
			plugin := agi.NewPlugin(env)
			require.EqualError(t, plugin.Add([]string{"one", "two", "three", "four"}), exp)
			rec.show(t, "errors_if_there_are_more_than_three_arguments")
		})
	})
}

func TestPlugin_MakeDirectoryIfNotExist(t *testing.T) {
	t.Parallel()

	t.Run("errors if path exists but is not a directory", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		exp := "failure while executing plugin command - exists but is not a directory - testdata/make_directory_if_not_exists_file"
		require.EqualError(t, plugin.MakeDirectoryIfNotExists("testdata/make_directory_if_not_exists_file"), exp)
		rec.show(t, "errors_if_path_exists_but_is_not_a_directory")
	})

	t.Run("errors if path exists but is not a directory", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		exp := "failure while executing plugin command - failed to create directory - testdata/does_not_exist/does_not_exist"
		require.EqualError(t, plugin.MakeDirectoryIfNotExists("testdata/does_not_exist/does_not_exist"), exp)
		rec.show(t, "errors_if_directory_creation_fails")
	})

	t.Run("succeeds if path exists and is directory", func(t *testing.T) {
		t.Parallel()

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		require.NoError(t, plugin.MakeDirectoryIfNotExists("testdata/make_directory_if_not_exists_folder"))
		rec.show(t, "succeeds_if_path_exists_and_is_directory")
	})
	t.Run("succeeds if directory is created", func(t *testing.T) {
		t.Parallel()

		tmp := t.TempDir()
		path := filepath.Join(tmp, "new_directory")

		env, rec := newTestEnv(t)
		plugin := agi.NewPlugin(env)
		require.NoError(t, plugin.MakeDirectoryIfNotExists(path))
		rec.show(t, "succeeds_if_path_exists_and_is_directorydirectory_is_created")
	})
}
