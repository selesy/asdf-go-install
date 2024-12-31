package manifest_test

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/selesy/asdf-go-install/internal/config/configtest"
	"github.com/selesy/asdf-go-install/internal/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

const (
	name = "go-enum"
	pkg  = "github.com/abice/" + name
)

func TestNew(t *testing.T) {
	t.Parallel()

	man := manifest.New(name, pkg, packageURL(t))
	assert.Equal(t, expectedManifestVersion(t), man.ManifestVersion())
	assert.Equal(t, name, man.PluginName())
	assert.Equal(t, pkg, man.PluginPackage())
	assert.Equal(t, packageURL(t), man.GitRepository())
	assert.Nil(t, man.GitReference())
}

func TestRead(t *testing.T) {
	t.Parallel()

	cfg, _, _ := configtest.NewConfig(t, []string{"ASDF_DATA_DIR=testdata"}, []string{})

	man, err := manifest.Read(cfg, name)
	require.NoError(t, err)
	assert.Equal(t, expectedManifestVersion(t), man.ManifestVersion())
	assert.Equal(t, name, man.PluginName())
	assert.Equal(t, pkg, man.PluginPackage())
	assert.Equal(t, packageURL(t), man.GitRepository())
	assert.Equal(t, tagReference(t), man.GitReference())
}

func TestManifest_WithGitReference(t *testing.T) {
	t.Parallel()

	man1 := manifest.New(name, pkg, packageURL(t))
	man2 := man1.WithGitReference(tagReference(t))

	assert.Nil(t, man1.GitReference())
	assert.Equal(t, tagReference(t), man2.GitReference())
}

func TestManifest_Write(t *testing.T) {
	t.Parallel()

	dataDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dataDir, "plugins", name), 0o755))

	cfg, _, _ := configtest.NewConfig(t, []string{"ASDF_DATA_DIR=" + dataDir}, []string{})

	man := manifest.New(name, pkg, packageURL(t))
	man = man.WithGitReference(tagReference(t))

	require.NoError(t, man.Write(cfg, name))

	exp := golden.Get(t, filepath.Join("plugins", name, manifest.ManifestFilename))

	act, err := os.ReadFile(filepath.Join(dataDir, "plugins", name, manifest.ManifestFilename))
	require.NoError(t, err)
	assert.JSONEq(t, string(exp), string(act))
}

func expectedManifestVersion(t *testing.T) *semver.Version {
	t.Helper()

	expVers, err := semver.NewVersion("v1")
	require.NoError(t, err)

	return expVers
}

func packageURL(t *testing.T) *url.URL {
	t.Helper()

	pkgURL, err := url.Parse("https://" + pkg + ".git")
	require.NoError(t, err)

	return pkgURL
}

func tagReference(t *testing.T) *plumbing.Reference {
	t.Helper()

	return plumbing.NewReferenceFromStrings("v0.6.0", "919e61c0174b91303753ee3898569a01abb32c97")
}
