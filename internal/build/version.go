// Package build exposes version information injected at build time via ldflags.
package build

// These are set by GoReleaser (and make build) via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
