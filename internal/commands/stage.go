package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/skaramicke/git-release/internal/git"
	"github.com/skaramicke/git-release/internal/release"
)

func newStageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stage [scope]",
		Short: "Create or increment a release candidate",
		Long: `Creates or increments a release candidate tag.

Without a scope, continues whatever is already in flight (or starts a patch RC).
With a scope (minor, major), applies the requested bump to the latest prod tag.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runStage,
		Example: `  git release stage           # continue in-flight RC, or create patch-rc
  git release stage minor     # target next minor version
  git release stage major     # target next major version
  git release stage --dry-run # preview without pushing`,
	}
}

func runStage(cmd *cobra.Command, args []string) error {
	scope := release.ScopeNone
	if len(args) == 1 {
		var err error
		scope, err = parseScope(args[0])
		if err != nil {
			return err
		}
	}

	ctx, err := loadContext(cmd, true)
	if err != nil {
		cmd.PrintErrf("✗ %v\n", err)
		return err
	}

	next, err := release.NextStage(ctx.state, scope)
	if err != nil {
		ctx.print.Errorf("%v", err)
		return err
	}

	tagName := next.String(ctx.cfg.TagPrefix)

	if ctx.dryRun {
		ctx.print.DryRun(fmt.Sprintf("would create and push tag %s → HEAD", tagName))
		// Show stale note if applicable
		if ctx.state.InFlightRC != nil && scope != release.ScopeNone {
			old := ctx.state.InFlightRC.String(ctx.cfg.TagPrefix)
			if old != tagName {
				ctx.print.DryRun(fmt.Sprintf("note: %s would become stale (not deleted)", old))
			}
		}
		return nil
	}

	if err := git.CreateAndPushTag(ctx.dir, ctx.cfg.Remote, tagName, "HEAD", ctx.cfg.SignTags, false); err != nil {
		ctx.print.Errorf("failed to create/push tag: %v", err)
		return err
	}

	ctx.print.TagCreated(tagName, ctx.cfg.TagPrefix, false)

	// Note stale old RC
	if ctx.state.InFlightRC != nil {
		old := ctx.state.InFlightRC.String(ctx.cfg.TagPrefix)
		if old != tagName {
			ctx.print.Success(fmt.Sprintf("note: %s is now stale (not deleted)", old))
		}
	}

	return nil
}
