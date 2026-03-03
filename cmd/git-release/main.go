package main

import (
	"errors"
	"os"

	"github.com/skaramicke/git-release/internal/commands"
)

func main() {
	err := commands.Root().Execute()
	if err == nil {
		os.Exit(0)
	}
	if errors.Is(err, commands.ErrAbort) {
		os.Exit(1)
	}
	os.Exit(2)
}
