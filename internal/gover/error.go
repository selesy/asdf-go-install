package gover

import "errors"

// ErrContainsBuildMetadata is returned when a valid semantic version
// contains build metadata and therefore can't be parsed as a Go version.
var ErrContainsBuildMetadata = errors.New("version contains build metadata")

// ErrMissingLeadingV is returned when an otherwise valid semantic version
// is missing the leading "v" required by Go versions.
var ErrMissingLeadingV = errors.New("version is missing leading \"v\"")

// ErrNoStableVersion is returned when a Collection contains only
// pre-release versions (including pseudo-versions.)
var ErrNoStableVersion = errors.New("no stable Go versions were found in the collection")
