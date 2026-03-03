package main

import (
	"errors"
	"os"

	"github.com/skaramicke/git-release/internal/build"
	"github.com/skaramicke/git-release/internal/commands"
)

// These vars are populated by GoReleaser via -ldflags.
// They live here so the linker path matches: main.version etc.
var (
	version = build.Version
	commit  = build.Commit
	date    = build.Date
)

func main() {
	build.Version = version
	build.Commit = commit
	build.Date = date

	err := commands.Root().Execute()
	if err == nil {
		os.Exit(0)
	}
	if errors.Is(err, commands.ErrAbort) {
		os.Exit(1)
	}
	os.Exit(2)
}
