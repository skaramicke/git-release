# Agent Instructions

## Maintaining This File

**This file is your primary context.** Sessions are recycled (killed and restarted) rather than compacted, so every new session starts fresh with only this file, beads, and Claude's memories. If this file is wrong or stale, you will do the wrong thing.

**Rules:**
- When you change architecture, add commands, rename files, or shift patterns — update this file in the same commit
- Keep it concise. Remove outdated sections rather than accumulating. Target <100 lines of prose
- Focus on: what the project IS, how to build/test, key patterns, and gotchas
- Do NOT add session-specific notes, TODOs, or in-progress work here (use beads for that)

This project uses **grit tasks** for issue tracking.

## Quick Reference

```bash
grit tasks ready              # Find available work
grit tasks show <id>          # View issue details
grit tasks update <id> --status in_progress  # Claim work
grit tasks close <id>         # Complete work
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

