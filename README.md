# AgentSkeleton

Language:

- [English](README.md)
- [中文](README.zh-CN.md)

AgentSkeleton is an AI-first documentation guidance tool for building AI-friendly project structures through conversation.

It is built around a simple model:

- documentation blueprints are the foundation
- a CLI is the product core
- Codex and Claude Code are the primary collaboration hosts

It is designed for two primary use cases:

- guide new projects into a clear, AI-friendly documentation structure
- reshape existing projects into a documented structure that is easier for AI tools and humans to maintain together

The CLI does not replace large models. It organizes the guidance flow, captures structured context, and tells the model what documentation should exist next.

## What It Is

AgentSkeleton is not a business application. It is a documentation architecture guidance system.

Its job is to help users build AI-friendly repository documentation through guided conversation, without touching business code.

## Goals

- Provide reusable documentation blueprints as baseline assets.
- Deliver a CLI as the primary product interface.
- Support both Codex-style agent workflows and Claude Code workflows.
- Keep the core structure shared across agent modes, with minimal adapter-specific differences.
- Help both greenfield projects and legacy project documentation reshaping.
- Keep the product focused on guidance and documentation, not business code generation.

## Core Principles

- AI-first collaboration
- Shared structure across agent modes
- Explicit repository rules
- Stable defaults before heavy customization
- Documentation reshaping, not only greenfield setup
- New projects may adopt recommended structures; legacy projects should document and respect the structure they already have

## Non-Goals For The First Stage

- Direct business code generation.
- Deep language-specific project scaffolding.
- A heavy plugin system before the core workflow is proven.

## MVP Scope

The first milestone focuses on definition and structure.

- Guided documentation setup for new repositories.
- Documentation reshaping guidance for existing repositories.
- Generation of collaboration files such as `README.md`, `AGENTS.md`, and `CLAUDE.md`.
- Structured question flow that tells the model what to document next.

## Supported Agent Modes

AgentSkeleton will support:

- Codex / agent mode
- Claude / Claude Code mode

The support strategy is:

- One shared core structure
- One shared documentation blueprint foundation
- Small adapter-specific instruction files where required

This keeps maintenance cost under control and avoids maintaining two separate project systems.

## Product Model

AgentSkeleton is intended to work alongside large models rather than replacing them.

The default expectation is:

- the CLI guides the conversation and organizes context
- Codex or Claude Code writes the actual document drafts
- humans define goals, constraints, and acceptance standards

This rule should shape both the product design and the repository workflow. See `docs/principles.md` for the full baseline.

## Installation

For the current release track, install AgentSkeleton by building the CLI locally:

```bash
go build -o agentskeleton ./cmd/agentskeleton
```

Or use the included local targets:

```bash
make build
make test
make smoke
```

To inspect the built CLI version:

```bash
./agentskeleton version --format json
```

To refresh existing context from repository facts without asking the same structural questions again:

```bash
./agentskeleton update --project /path/to/project --output-dir /path/to/output --format json
```

The repository also includes GitHub Actions workflows for CI and tagged release builds under `.github/workflows/`.

## v0.1.0 Direction

The target for `v0.1.0` is the first public zero-command release for Codex and Claude Code.

That means:

- the CLI remains the protocol core
- host integrations carry the zero-command experience
- users should stay in conversation instead of manually operating AgentSkeleton commands

See these docs for the current release-track definition:

- `docs/v0.1.0-gap-analysis.md`
- `docs/v0.1.0-implementation-plan.md`
- `docs/known-limitations.md`

## Reading Order

Default reading order for current repository truth:

1. `docs/current-capabilities.md`
2. `docs/features/README.md`
3. `docs/README.md`

Use release-track docs only when you need a version snapshot rather than the current repository state.

## Reading Order

Default reading order for current repository truth:

1. `docs/current-capabilities.md`
2. `docs/features/README.md`
3. `docs/README.md`

Use release-track docs only when you need a version snapshot rather than the current repository state.

## Repository Layout

```text
.
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── cmd/
├── docs/
├── internal/
└── templates/
```

## Planned CLI Direction

The CLI is expected to become the main entry point for users. Initial command areas are likely to include:

- `init-docs`: guide a new project into an AI-friendly documentation structure
- `reshape-docs`: guide an existing project through documentation reshaping
- `plan`: summarize what documents should exist next
- `next`: provide the next structured questions for the conversation
- `response`: validate/evaluate model output and optionally apply accepted answers
- `prompt`: generate initial or repair prompts from context
- `workflow`: run one bundled step (`plan + prompt + next`) with optional response apply, planned file materialization via `--write-plan-files`, retry repair packaging via `--auto-repair`, and process snapshots via `--persist-trace`
- `update`: infer safe context updates from the existing repository and refresh the plan without another confirmation round
- `plan` and `workflow` now expose `current_priority` so the host model knows which document should be drafted next
- `focus-doc` turns that priority into a drafting context package, while `review_candidates` expose backtracking work for already-generated documents
- `focus-doc` also exposes `review_after_draft`, so forward drafting and backward convergence are planned together
- `response --apply` now returns `post_apply_plan`, so hosts can continue immediately after a successful write-back
- review work is scoped to the latest change batch instead of being cached as a permanent state
- `focus-doc` carries the current `change_batch_id` so hosts can detect stale drafting packages

See `docs/agent-prompts.md` for the host-model prompt expectations so structured retries work.
See `docs/cli-runbook.md` for the end-to-end command sequence.
See `docs/host-integration.md` for how Codex and Claude Code should integrate with AgentSkeleton in practice.
See `docs/codex-integration.md` and `docs/claude-integration.md` for the current official host integration artifacts.
Path defaults: use `--project` to identify the target project; AgentSkeleton process artifacts live under `<output-dir>/.agentskeleton`, while final docs and skills belong under `<output-dir>/...`.

The first minimal CLI skeleton now exists in `cmd/agentskeleton`, with the output contract defined in `docs/cli-contract.md`.

## Current Status

This repository has moved beyond initial definition and now includes:

- the runnable CLI core
- host integration specs for zero-command behavior
- initial Codex and Claude Code skill artifacts
- shared host validation scenarios and report templates
- initial CI and tagged release workflow scaffolding

## Contributing Direction

In the early phase, decisions should favor clarity, explicit conventions, and stable defaults over premature flexibility.

## Commit Convention

Use a Jira-style identifier in commit messages.

Recommended format:

- `docs [AG-001]: align product positioning and blueprint strategy`
- `feat [AG-001]: add initial documentation guidance flow`
- `fix [AG-001]: correct document planning output`
