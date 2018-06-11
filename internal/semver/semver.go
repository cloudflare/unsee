package semver

import (
	"strings"

	"github.com/blang/semver"
)

// MustParse is like semver.MustParse, but first attempts to trim
// distribution metadata from the version string.
func MustParse(s string) semver.Version {
	return semver.MustParse(trimDistribution(s))
}

// MustParseRange is like semver.MustParseRange
func MustParseRange(s string) semver.Range {
	return semver.MustParseRange(s)
}

// trimDistribution attempts to trim the passed version string of
// and distribution-specific version numbers.
func trimDistribution(version string) string {
	return strings.SplitN(version, "~", 2)[0]
}
