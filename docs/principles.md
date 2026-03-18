# Principles

## Baseline

AgentSkeleton is defined by a small set of non-negotiable product principles.

## AI-First Collaboration

- Large models are expected to draft repository documents.
- Humans provide direction, constraints, review, and acceptance.
- The product guides documentation work rather than generating business code.

## Dual-Agent Support

- Codex and Claude Code are both first-class workflows.
- The repository should preserve one shared structure wherever possible.
- Agent-specific files should exist only when a real difference must be expressed.

## Documentation-Blueprint Product Model

- Documentation blueprints are not secondary assets. They are the core of the product.
- The CLI exists to guide, inspect, and maintain those assets.
- Repository structure should reflect this product model clearly.

## Structure Strategy

- New projects may use recommended repository patterns.
- The preferred default for new application-style projects is an `internal/app`-oriented structure.
- Existing projects should not be forced into a new code layout by default.
- For legacy projects, the product should help explain and document the current structure before suggesting structural change.

## CLI-Core, Host-Assisted Delivery

- The product core is a CLI.
- Codex and Claude Code are primary host environments, not the product itself.
- The CLI should produce structured guidance that host models can immediately use.

## Explicit Repository Rules

- Important context must live in the repository, not in private memory.
- Naming, layout, and workflow constraints should be documented.
- Hidden conventions are a defect.

## External-Ready Quality Bar

- The project should be usable by external users, not only by the internal team.
- Internal use should benefit from the same clarity required by external use.
- Early versions should still keep scope narrow and defaults strong.
