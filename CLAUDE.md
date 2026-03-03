# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`git-release` is a git plugin for opinionated semver release management with first-class release candidate support. Installed as `git-release`, it becomes available as `git release`.

The project is currently in the **design/spec phase** — the full specification lives in `README.md`. The implementation does not yet exist.

## Implementation Target

The tool must be a single self-contained shell script or compiled binary named `git-release`. Key constraints from the spec:

- **Tag state**: Always fetch from remote (`git ls-remote --tags <remote>`) before computing next version; never trust local tag state alone
- **Semver parsing**: Must handle full semver spec including pre-release comparison rules; prefer a dedicated library over hand-rolled parsing
- **Output**: All output to stdout; errors and warnings to stderr
- **Confirmation prompts**: Must read from `/dev/tty` directly (not stdin) so the command works safely in pipelines
- **No force-push**: If a tag already exists remotely, error with a clear message
- **Exit codes**: 0 success, 1 user abort, 2 error

## Tag Format

- Production: `v1.2.3` (no suffix)
- First RC: `v1.2.3-rc` (not `-rc.1`)
- Subsequent RCs: `v1.2.3-rc.2`, `v1.2.3-rc.3`, …
- Tag prefix is configurable via `git config release.tagPrefix`

## Commands to Implement

- `git release` — promote in-flight RC to prod, or bump patch directly
- `git release stage [minor|major]` — create/increment an RC
- `git release status` — show current state summary
- `git release ls` — list all release tags in descending semver order

All mutating commands support `--dry-run` / `-n` and `--yes` / `-y`.

## Configuration

Via `git config` (keys under `release.*`): `tagPrefix`, `remote`, `releaseBranch`, `requireClean`, `signTags`.

## Issue Tracking

This repo uses **beads** (`bd`) for issue tracking — not GitHub Issues, not markdown files.

```bash
bd ready                    # Find available work
bd show <id>                # View issue details
bd update <id> --status=in_progress  # Claim work
bd close <id>               # Complete work
bd sync --flush-only        # Export to JSONL before ending session
```

Do NOT use TodoWrite, TaskCreate, or markdown task lists. Always create a beads issue before writing code.

## Session Completion

Before ending a session, run:
```bash
bd sync --flush-only
```

Note: No git remote is configured — issues are saved locally only.

<!-- git-release prime -->
# git-release — Agent Quick Reference

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
