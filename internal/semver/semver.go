// Package semver provides version parsing, comparison, and manipulation for
// git-release tags. Supports production versions (v1.2.3) and release
// candidates (v1.2.3-rc, v1.2.3-rc.2, …).
package semver

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Version represents a parsed release tag. RC=1 means the first candidate
// (-rc with no number); RC=2 means -rc.2, etc.
type Version struct {
	Major, Minor, Patch int
	IsRC                bool
	RC                  int // 1 = first RC (-rc), 2+ = numbered (-rc.N)
}

// ParseTag parses a tag like "v1.2.3" or "v1.2.3-rc.2" using the given prefix.
// Returns an error if the tag doesn't match the expected format.
func ParseTag(tag, prefix string) (Version, error) {
	if !strings.HasPrefix(tag, prefix) {
		return Version{}, fmt.Errorf("tag %q does not start with prefix %q", tag, prefix)
	}
	rest := tag[len(prefix):]

	var v Version

	// Split off RC suffix if present
	base := rest
	if idx := strings.Index(rest, "-rc"); idx != -1 {
		v.IsRC = true
		base = rest[:idx]
		suffix := rest[idx+3:] // everything after "-rc"
		switch {
		case suffix == "":
			v.RC = 1
		case strings.HasPrefix(suffix, "."):
			n, err := strconv.Atoi(suffix[1:])
			if err != nil || n < 2 {
				return Version{}, fmt.Errorf("invalid RC suffix in tag %q", tag)
			}
			v.RC = n
		default:
			return Version{}, fmt.Errorf("invalid RC suffix in tag %q", tag)
		}
	}

	parts := strings.Split(base, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("tag %q: expected major.minor.patch, got %q", tag, base)
	}

	nums := []*int{&v.Major, &v.Minor, &v.Patch}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return Version{}, fmt.Errorf("tag %q: non-numeric component %q", tag, p)
		}
		*nums[i] = n
	}

	return v, nil
}

// String formats the version as a tag string using the given prefix.
func (v Version) String(prefix string) string {
	base := fmt.Sprintf("%s%d.%d.%d", prefix, v.Major, v.Minor, v.Patch)
	if !v.IsRC {
		return base
	}
	if v.RC == 1 {
		return base + "-rc"
	}
	return fmt.Sprintf("%s-rc.%d", base, v.RC)
}

// Base returns the production version with the same major.minor.patch,
// stripping any RC designation.
func (v Version) Base() Version {
	return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch}
}

// GreaterThan returns true if v is strictly greater than other using semver
// ordering. For the same major.minor.patch, a production release is greater
// than any RC; higher RC numbers are greater than lower ones.
func (v Version) GreaterThan(other Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	if v.Patch != other.Patch {
		return v.Patch > other.Patch
	}
	// Same base version — prod > RC; higher RC > lower RC
	switch {
	case !v.IsRC && !other.IsRC:
		return false
	case !v.IsRC && other.IsRC:
		return true
	case v.IsRC && !other.IsRC:
		return false
	default:
		return v.RC > other.RC
	}
}

// Equal returns true if v and other represent the same version.
func (v Version) Equal(other Version) bool {
	return v.Major == other.Major &&
		v.Minor == other.Minor &&
		v.Patch == other.Patch &&
		v.IsRC == other.IsRC &&
		v.RC == other.RC
}

// BumpPatch returns a new version with the patch component incremented.
func (v Version) BumpPatch() Version {
	return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
}

// BumpMinor returns a new version with the minor component incremented and
// patch reset to 0.
func (v Version) BumpMinor() Version {
	return Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
}

// BumpMajor returns a new version with the major component incremented and
// minor + patch reset to 0.
func (v Version) BumpMajor() Version {
	return Version{Major: v.Major + 1, Minor: 0, Patch: 0}
}

// FirstRC returns the first release candidate for this production version.
func (v Version) FirstRC() Version {
	return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch, IsRC: true, RC: 1}
}

// NextRC returns the next release candidate in sequence.
// Calling NextRC on a non-RC version returns its first RC.
func (v Version) NextRC() Version {
	if !v.IsRC {
		return v.FirstRC()
	}
	return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch, IsRC: true, RC: v.RC + 1}
}

// SortDesc sorts a slice of versions in descending order (highest first).
func SortDesc(versions []Version) {
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].GreaterThan(versions[j])
	})
}
