# AGENTS.md

Language:

- [English](AGENTS.md)
- [中文](AGENTS.zh-CN.md)

This repository is an AI-friendly scaffold project. Agent work in this repository should optimize for explicit structure, low ambiguity, and reusable conventions.

## Working Principles

- Prefer shared project conventions over tool-specific shortcuts.
- Treat the CLI as the primary product surface.
- Treat templates as core assets owned by the repository.
- Keep Codex and Claude Code support aligned unless divergence is necessary.
- Prefer readable repository structure over clever automation.
- Treat the repository as AI-authored by default, not manually coded.

## Repository Intent

This project exists to help users:

- start new projects with agent-friendly structure
- migrate existing projects into an agent-friendly structure
- generate agent collaboration files
- inspect repositories for missing structure and conventions

## Documentation Rules

- Put product intent in `README.md`.
- Put deeper product decisions in `docs/`.
- Keep agent-specific guidance minimal and explicit.
- Avoid hidden conventions that are not written in the repository.
- Host integration rules belong in `docs/host-integration.md`.

## Implementation Direction

- Main language: Go
- Primary interface: CLI
- Shared scaffold assets: `templates/`
- Implementation code: `cmd/` and `internal/`

## Stage-1 Expectation

The first stage should focus on repository definition, structure, and documentation before building out a full CLI implementation.

## Authoring Constraint

This project is expected to be implemented through AI agents.

- Humans define direction, standards, and acceptance criteria.
- Agents generate and modify repository contents.
- Manual source-code authoring is not the intended operating model.
