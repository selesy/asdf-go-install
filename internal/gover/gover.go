// Package gover provides a parser and utility functions for Go module
// version numbers.
package gover

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/selesy/asdf-go-install/internal/config"
)

const (
	// GoVersionPrefix is the prefix required for Go version strings.
	GoVersionPrefix = "v"

	// PseudoVersionRegexp is a pattern that matches the suffix present
	// on Go pseudo-versions.
	PseudoVersionRegexp = "^[0-9]{14}-[0-9a-f]{12}$"
)

var pseudoVersionRegexp = regexp.MustCompile(PseudoVersionRegexp)

// NewVersion creates a semver.Version using the Go [module version numbering]
// syntax.
//
// A Go module version number differs from a true [semantic version] number
// in three minor ways:
//
//  1. The leading "v", which is optional in a lax semantic version and
//     not allowed in a strict semantic version, is required by a Go
//     module version number.
//
//  2. Build metadata, which is appended after a "+" in a semantic version
//     number, is not allowed in a Go version number.
//
//  3. A Go pseudo-version number is a specialization of a semantic
//     version with a carefully formatted pre-release suffix.
//
// [module version numbering]: https://go.dev/doc/modules/version-numbers
// [semantic version]: https://semver.org/
func NewVersion(v string) (*semver.Version, error) {
	// A Go versin requires the leading
	if !strings.HasPrefix(v, GoVersionPrefix) {
		return nil, fmt.Errorf("%w: parsed %s", ErrMissingLeadingV, v)
	}

	// A Go version should parse as a strict semver if the leading v is
	// removed.
	ver, err := semver.StrictNewVersion(strings.TrimPrefix(v, "v"))
	if err != nil {
		return nil, err
	}

	// A Go version must not contain build metadata.
	if ver.Metadata() != "" {
		return nil, fmt.Errorf("%w: parsed %s", ErrContainsBuildMetadata, v)
	}

	// Assuming the above checks succeed, return a lax semver so that
	// Original() will still return a proper Go version.
	return semver.NewVersion(v)
}

// IsPrerelease returns a boolean value indicating whether the Go
// version has a pre-release suffix
func IsPrerelease(v *semver.Version) bool {
	return v.Prerelease() != ""
}

// IsPseudoVersion returns a boolean value indicating whether the
// Go version has a pre-release suffix that's formatted as
// a pseudo-version.
func IsPseudoVersion(v *semver.Version) bool {
	return pseudoVersionRegexp.Match([]byte(v.Prerelease()))
}

// IsRelease returns a boolean value indicating whether the Go version
// references a release (is missing a pre-release suffix.)
func IsRelease(v *semver.Version) bool {
	return !IsPrerelease(v)
}

// Stores a sorted collection of Go module version numbers.
type Collection struct {
	col semver.Collection
}

// NewCollection sorts and stores the provided Go module version numbers.
func NewCollection(vers ...*semver.Version) *Collection {
	col := Collection{
		col: semver.Collection(vers),
	}

	sort.Sort(col.col)

	return &col
}

// All returns the underlying slice of Go module version numbers.
func (c *Collection) All() []*semver.Version {
	return c.col
}

// LatestStable returns the Go module version number for the most recent
// released version of the module.
//
// If there is no release version in the collection, an ErrNoStableVersion
// error is returned.
func (c *Collection) LatestStable() (*semver.Version, error) {
	var col = semver.Collection(make([]*semver.Version, len(c.col)))
	copy(col, c.col)

	slices.Reverse(col)

	for _, ver := range col {
		if IsRelease(ver) {
			return ver, nil
		}
	}

	return nil, ErrNoStableVersion
}

// Len returns the number of Go module version numbers stored in the
// Collection.
func (c *Collection) Len() int {
	return len(c.col)
}

// String returns a space-delimited string representation of the Go
// module version Collection starting with the lowest verion ending with
// the most recent version.
func (c *Collection) String() string {
	var vers = make([]string, len(c.col))

	for i, ver := range c.col {
		vers[i] = ver.Original()
	}

	return strings.Join(vers, " ")
}

type Collector func(cfg *config.Config, pkg string) (*Collection, error)
