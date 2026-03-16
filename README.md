# AgentSkeleton

Language:

- English: `README.md`
- 中文: `README.zh-CN.md`

AgentSkeleton is an AI-first project scaffold for building and maintaining software with agent collaboration in mind.

It is built around a simple model:

- templates are the foundation
- a CLI is the main user interface
- both Codex and Claude Code are supported by design

It is designed for two primary use cases:

- initialize new projects with a clear, agent-readable structure
- migrate existing projects into a structure that is easier for AI tools and humans to maintain together

The repository template and bundled assets are the foundation that the CLI will generate, inspect, and maintain.

## What It Is

AgentSkeleton is not a business application. It is a scaffold for project construction and maintenance.

Its job is to make projects easier to understand, safer to evolve, and more predictable for both humans and AI agents.

## Goals

- Provide a reusable project template as the baseline asset.
- Deliver a CLI as the primary user-facing interface.
- Support both Codex-style agent workflows and Claude Code workflows.
- Keep the core structure shared across agent modes, with minimal adapter-specific differences.
- Help both greenfield projects and legacy project migration.
- Keep the project fully AI-authored, without manual source-code implementation.

## Core Principles

- AI-first authoring
- Shared structure across agent modes
- Explicit repository rules
- Stable defaults before heavy customization
- Migration support, not only greenfield setup

## Non-Goals For The First Stage

- Full automation for every project type.
- Deep language-specific generators for many stacks at once.
- A heavy plugin system before the core workflow is proven.

## MVP Scope

The first milestone focuses on definition and structure.

- Project initialization for new repositories.
- AI migration guidance for existing repositories.
- Generation of collaboration files such as `AGENTS.md` and `CLAUDE.md`.
- Repository structure scanning and basic diagnostics.

## Supported Agent Modes

AgentSkeleton will support:

- Codex / agent mode
- Claude / Claude Code mode

The support strategy is:

- One shared core structure
- One shared template foundation
- Small adapter-specific instruction files where required

This keeps maintenance cost under control and avoids maintaining two separate project systems.

## Authoring Principle

AgentSkeleton is intended to be built through AI collaboration rather than manual coding.

The default expectation is:

- no hand-written production code
- repository changes are proposed and implemented by AI agents
- humans define goals, constraints, and acceptance standards

This rule should shape both the product design and the repository workflow. See `docs/principles.md` for the full baseline.

## Repository Layout

```text
.
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── cmd/
├── docs/
├── examples/
├── internal/
└── templates/
```

## Planned CLI Direction

The CLI is expected to become the main entry point for users. Initial command areas are likely to include:

- `init`: create a new project from the scaffold
- `migrate`: adapt an existing project to the scaffold
- `generate`: create agent-facing support files
- `scan`: inspect repository structure and report gaps

## Current Status

This repository is in the initial definition stage. The first push should establish:

- core documentation
- baseline structure
- naming conventions
- agent support strategy
- CLI direction

## Contributing Direction

In the early phase, decisions should favor clarity, explicit conventions, and stable defaults over premature flexibility.
