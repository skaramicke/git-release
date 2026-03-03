package commands

import (
	"github.com/spf13/cobra"
)

const primeOutput = `# git-release — Agent Quick Reference

## What it does
git-release manages semver tags and release candidates on the current repo.
Tags follow the pattern: v1.2.3 (prod) and v1.2.3-rc / v1.2.3-rc.2 (RC).

## Read-only commands (safe to run any time)

  git release status        Show current state: latest prod, in-flight RC,
                            and what the next stage/release would create.

  git release ls            List all release tags descending, annotated with
                            [latest prod], [in-flight RC], [superseded RC].

## Mutating commands (fetch remote tag state first, push on success)

  git release stage         Continue whatever is in-flight (or start a patch RC).
  git release stage minor   Target the next minor version.
  git release stage major   Target the next major version.

  git release               Promote in-flight RC → prod, or create a direct patch
                            release if no RC is in flight.
  git release minor         Direct minor release (prompts if RC is in flight).
  git release major         Direct major release (prompts if RC is in flight).

## Flags (all mutating commands)

  -n, --dry-run   Print what would happen without pushing any tags.
  -y, --yes       Skip confirmation prompts (use in CI pipelines).

## Typical workflow

  git release status              # where are we?
  git release stage               # create v1.2.3-rc
  git release stage               # create v1.2.3-rc.2
  git release                     # promote → v1.2.3

## Rules to know

  - "stage" without a scope always continues whatever is in flight.
  - Staging a lower scope than the in-flight RC is an error.
  - "release" without a scope promotes the RC; never bypasses it silently.
  - Stale RCs are left in place (not deleted) when a higher-scope RC supersedes them.
  - All operations are idempotent on read; never force-push.
`

func newPrimeCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "prime",
		Short:  "Print a quick-reference guide for AI agents",
		Long:   "Prints a concise reference of all git-release commands and rules, intended to be injected into an AI agent's context at the start of a session.",
		Args:   cobra.NoArgs,
		Hidden: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Print(primeOutput)
			return nil
		},
	}
}
