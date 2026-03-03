// Package commands contains all cobra CLI commands for git-release.
package commands

import (
	"github.com/spf13/cobra"
)

var (
	flagDryRun bool
	flagYes    bool
)

// Root returns the root cobra command.
func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "git-release [scope]",
		Short: "Opinionated semver release management",
		Long: `git-release manages semver tags and release candidates.

Run without a subcommand to promote an in-flight RC to production,
or create a patch release directly if no RC is in flight.

Scopes: patch (default), minor, major`,
		Args:              cobra.MaximumNArgs(1),
		RunE:              runRelease,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	root.PersistentFlags().BoolVarP(&flagDryRun, "dry-run", "n", false, "print what would happen without creating or pushing any tags")
	root.PersistentFlags().BoolVarP(&flagYes, "yes", "y", false, "answer all confirmation prompts affirmatively (for CI)")

	root.AddCommand(
		newStageCmd(),
		newStatusCmd(),
		newLsCmd(),
		newPrimeCmd(),
		newUpdateCmd(),
		newVersionCmd(),
	)

	return root
}
