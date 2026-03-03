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
