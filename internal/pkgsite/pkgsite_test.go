package pkgsite_test

import (
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selesy/asdf-go-install/internal/config/configtest"
	"github.com/selesy/asdf-go-install/internal/gover"
	"github.com/selesy/asdf-go-install/internal/pkgsite"
)

func TestRepostory(t *testing.T) {
	t.Parallel()

	const expRepo = "https://go.googlesource.com/vuln"

	cfg, _, _ := configtest.NewConfig(t, []string{}, []string{})

	repo, err := pkgsite.Repository(cfg, "golang.org/x/vuln/cmd/govulncheck")
	require.NoError(t, err)
	assert.Equal(t, expRepo, repo.String())
}

func TestVersions(t *testing.T) {
	t.Parallel()

	const (
		expVersStr = "v0.1.0 v0.2.0 v1.0.0 v1.0.1 v1.0.2 v1.0.3 v1.0.4 v1.1.0 v1.1.1 v1.1.2 v1.1.3"
	)

	var expVers []*semver.Version

	_ = expVers

	for _, expVerStr := range strings.Split(expVersStr, " ") {
		expVer, err := gover.NewVersion(expVerStr)
		require.NoError(t, err)

		expVers = append(expVers, expVer)
	}

	cfg, _, _ := configtest.NewConfig(t, []string{}, []string{})

	vers, err := pkgsite.Versions(cfg, "golang.org/x/vuln/cmd/govulncheck")
	require.NoError(t, err)

	all := vers.All()
	assert.GreaterOrEqual(t, len(all), len(expVers))

	for i, expVer := range expVers {
		assert.Equal(t, expVer, all[i])
	}
}
