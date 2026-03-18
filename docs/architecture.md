# Architecture Notes

## Core Model

AgentSkeleton has two layers:

- documentation blueprints stored in the repository
- a CLI that guides, plans, and validates documentation work

The repository is expected to evolve through AI-assisted documentation work, with humans acting as product owners, reviewers, and constraint setters.

## Support Model

The project should support both Codex and Claude Code through:

- one shared structure
- one shared documentation model
- small adapter-specific files only where needed

## Delivery Model

- Users work inside Codex or Claude Code conversations.
- AgentSkeleton runs as a CLI inside those environments.
- The CLI outputs structured prompts, plans, and document expectations.
- The host model turns those outputs into actual repository documents.

## Repository Structure Handling

- For new projects, AgentSkeleton may recommend a default structure to make documentation and collaboration easier.
- The current preferred default for application-style new projects is an `internal/app`-oriented layout.
- For existing projects, AgentSkeleton should document the structure that already exists instead of imposing a replacement architecture.

## Initial Directory Intent

- `cmd/`: CLI entrypoints
- `internal/`: internal implementation packages
- `templates/`: documentation blueprints and guidance assets
- `docs/`: product and design documentation
