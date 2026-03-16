# Architecture Notes

## Core Model

AgentSkeleton has two layers:

- template assets stored in the repository
- a CLI that applies, generates, and validates those assets

The repository is expected to evolve through AI-authored changes, with humans acting as product owners, reviewers, and constraint setters rather than primary code authors.

## Support Model

The project should support both Codex and Claude Code through:

- one shared structure
- one shared documentation model
- small adapter-specific files only where needed

## Initial Directory Intent

- `cmd/`: CLI entrypoints
- `internal/`: internal implementation packages
- `templates/`: scaffold assets and file templates
- `examples/`: example generated structures or reference projects
- `docs/`: product and design documentation
