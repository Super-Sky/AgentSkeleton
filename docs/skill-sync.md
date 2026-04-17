# Skill Sync

## Purpose

This document defines the source-of-truth relationship between shared docs and host skills.

## Current Model

The repository currently uses:

- shared docs in `docs/*`
- host-specific entry skills in `skills/*`

The skills should remain thin.

## Source Of Truth Rules

- shared workflow rules live in `docs/*`
- skills should point to shared docs rather than restating long bodies of rules
- current host behavior truth should be updated in docs first
- skills should then be updated to reflect those doc changes

## Current Skills

- `skills/agentskeleton-codex/SKILL.md`
- `skills/agentskeleton-claude/SKILL.md`

## Sync Expectations

When any of the following change:

- host routing rules
- zero-command flow
- stable protocol boundary
- escalation policy

the maintainer should review both host skills for sync.

## Anti-Pattern

Do not let Codex and Claude skills evolve separate workflow semantics unless the divergence is explicit and justified.
