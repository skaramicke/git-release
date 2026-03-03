package commands

import (
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current release state",
		Long:  `Prints a summary of the current release state: latest prod tag, in-flight RC, and what the next stage and release would be.`,
		Args:  cobra.NoArgs,
		RunE:  runStatus,
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx, err := loadContext(cmd, false)
	if err != nil {
		cmd.PrintErrf("✗ %v\n", err)
		return err
	}

	ctx.print.Status(ctx.state, ctx.cfg.TagPrefix)
	return nil
}
