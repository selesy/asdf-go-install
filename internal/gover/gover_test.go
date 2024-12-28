package gover_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selesy/asdf-go-install/internal/gover"
)

func TestNewVersion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ver       string
		expOrig   string
		expStr    string
		expErr    error
		expPre    bool
		expPseudo bool
		expRel    bool
	}{
		"pass with leading v": {
			ver:       "v1.0.0",
			expOrig:   "v1.0.0",
			expStr:    "1.0.0",
			expErr:    nil,
			expPre:    false,
			expPseudo: false,
			expRel:    true,
		},
		"pass with prerelease": {
			ver:       "v1.0.0-pre.1",
			expOrig:   "v1.0.0-pre.1",
			expStr:    "1.0.0-pre.1",
			expErr:    nil,
			expPre:    true,
			expPseudo: false,
			expRel:    false,
		},
		"pass with pseudoversion": {
			ver:       "v0.0.0-20170915032832-14c0d48ead0c",
			expOrig:   "v0.0.0-20170915032832-14c0d48ead0c",
			expStr:    "0.0.0-20170915032832-14c0d48ead0c",
			expErr:    nil,
			expPre:    true,
			expPseudo: true,
			expRel:    false,
		},
		"fail without leading v": {
			ver:    "1.0.0",
			expErr: gover.ErrMissingLeadingV,
		},
		"fail with metadata": {
			ver:    "v1.0.0+41f3b2d",
			expErr: gover.ErrContainsBuildMetadata,
		},
		"fail if not semantic": {
			ver:    "vA.B.C",
			expErr: semver.ErrInvalidCharacters,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v, err := gover.NewVersion(test.ver)
			require.ErrorIs(t, err, test.expErr)

			if err != nil {
				return
			}

			assert.Equal(t, test.expOrig, v.Original())
			assert.Equal(t, test.expStr, v.String())
			assert.Equal(t, test.expPre, gover.IsPrerelease(v))
			assert.Equal(t, test.expPseudo, gover.IsPseudoVersion(v))
			assert.Equal(t, test.expRel, gover.IsRelease(v))
		})
	}
}

func TestSortCollection(t *testing.T) {
	t.Parallel()

	mustNew := func(s string) *semver.Version {
		ver, err := gover.NewVersion(s)
		if err != nil {
			t.Fatal(err)
		}

		return ver
	}

	v4_5_6 := mustNew("v4.5.6")
	v2_2_4 := mustNew("v2.2.4")
	v0_0_0_20170915032832_14c0d48ead0c := mustNew("v0.0.0-20170915032832-14c0d48ead0c")
	v1_2_3 := mustNew("v1.2.3")
	v1_2_3_20170915032832_14c0d48ead0c := mustNew("v1.2.3-20170915032832-14c0d48ead0c")
	v1_2_1pre3 := mustNew("v1.2.1-pre.3")
	v1_2_1pre1 := mustNew("v1.2.1-pre.1")
	v1_2_1 := mustNew("v1.2.1")
	v1_2_1pre2 := mustNew("v1.2.1-pre.2")
	v5_6_7_20170915032832_14c0d48ead0c := mustNew("v5.6.7-20170915032832-14c0d48ead0c")
	v3_4_5 := mustNew("v3.4.5")

	col := gover.NewCollection(
		v4_5_6,
		v2_2_4,
		v0_0_0_20170915032832_14c0d48ead0c,
		v1_2_3,
		v1_2_3_20170915032832_14c0d48ead0c,
		v1_2_1pre3,
		v1_2_1pre1,
		v1_2_1,
		v1_2_1pre2,
		v5_6_7_20170915032832_14c0d48ead0c,
		v3_4_5,
	)

	assert.Equal(t, 11, col.Len())

	t.Run("All", func(t *testing.T) {
		t.Parallel()

		exp := []*semver.Version{
			v0_0_0_20170915032832_14c0d48ead0c,
			v1_2_1pre1,
			v1_2_1pre2,
			v1_2_1pre3,
			v1_2_1,
			v1_2_3_20170915032832_14c0d48ead0c,
			v1_2_3,
			v2_2_4,
			v3_4_5,
			v4_5_6,
			v5_6_7_20170915032832_14c0d48ead0c,
		}

		assert.Equal(t, exp, col.All())
	})

	t.Run("LatestStable", func(t *testing.T) {
		t.Parallel()

		act, err := col.LatestStable()
		require.NoError(t, err)
		assert.Equal(t, v4_5_6, act)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		exp := "v0.0.0-20170915032832-14c0d48ead0c v1.2.1-pre.1 v1.2.1-pre.2 v1.2.1-pre.3 v1.2.1 v1.2.3-20170915032832-14c0d48ead0c v1.2.3 v2.2.4 v3.4.5 v4.5.6 v5.6.7-20170915032832-14c0d48ead0c"

		assert.Equal(t, exp, col.String())
	})

	t.Run("No LatestStable", func(t *testing.T) {
		t.Parallel()

		col := gover.NewCollection(
			v0_0_0_20170915032832_14c0d48ead0c,
			v1_2_3_20170915032832_14c0d48ead0c,
			v1_2_1pre3,
			v1_2_1pre1,
			v1_2_1pre2,
			v5_6_7_20170915032832_14c0d48ead0c,
		)

		ver, err := col.LatestStable()
		require.ErrorIs(t, err, gover.ErrNoStableVersion)
		assert.Nil(t, ver)
	})
}
