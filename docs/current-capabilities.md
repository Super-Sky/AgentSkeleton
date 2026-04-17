# Current Capabilities

This document describes the current repository truth for AgentSkeleton.

It is not a version snapshot.

It should answer one question first:

- what does the repository currently support right now

## Current Product State

AgentSkeleton currently provides:

- a Go CLI core for documentation workflow orchestration
- a structured context model under `.agentskeleton`
- planning, drafting-package, and response-apply workflow commands
- an `update` command for safe context refresh from repository structure
- host integration specifications for zero-command usage
- initial Codex and Claude Code skill artifacts
- validation scenarios, report templates, and seeded validation reports
- CI, release-build scaffolding, and local smoke-test tooling

## Current Stable Workflow Core

The narrow stable workflow core is:

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

These commands are the intended public host contract for the `v0.1.0` release track.

The repository also currently includes:

- `update`

`update` is a refresh helper for already-initialized context. It is useful, but should currently be treated as a helper command rather than the narrowest host contract.

## Current Host Integration State

The repository currently includes:

- `skills/agentskeleton-codex/SKILL.md`
- `skills/agentskeleton-claude/SKILL.md`
- shared host rules in:
  - `host-action-mapping.md`
  - `zero-command-flow.md`
  - `protocol-stability.md`

This means the host integration model is defined and scaffolded.

It does not yet mean that real host validation is complete.

## Current Validation State

The repository currently has:

- validation scenarios
- validation report templates
- seeded report files
- local CLI smoke validation evidence
- real local CLI validation evidence on:
  - `sast_server`
  - `sast_task_assistant`

What is still missing:

- real Codex host-run validation evidence
- real Claude Code host-run validation evidence

## Current Release State

The repository currently has:

- CI workflow
- tag-triggered release build workflow
- CLI version command
- local build, test, smoke, and release-build targets
- known limitations and changelog scaffolding

This means release hardening is underway.

It does not yet mean `v0.1.0` is ready to cut.

## Current Blocking Gaps

The main remaining release blockers are:

- capture real Codex validation evidence
- capture real Claude Code validation evidence
- adjust host skills based on observed host behavior
- verify tagged release behavior in GitHub Actions

## Reading Guidance

Read next:

1. `features/README.md`
2. `host-action-mapping.md`
3. `zero-command-flow.md`
4. `protocol-stability.md`

Use release-track docs only when you need planning detail or scope history.
