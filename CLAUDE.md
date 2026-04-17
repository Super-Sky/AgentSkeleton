# CLAUDE.md

Language:

- [English](CLAUDE.md)
- [中文](CLAUDE.zh-CN.md)

This repository supports Claude Code as one of the expected agent workflows.

Claude-specific guidance should remain thin. The default rule is to follow the shared repository conventions first and only introduce Claude-specific instructions where they are truly necessary.

## Priority Order

1. Follow the repository structure and documented project intent.
2. Reuse shared templates and conventions.
3. Keep Claude-specific deviations explicit and minimal.
4. Preserve the AI-authored workflow and avoid assuming manual coding as the default path.

## Repository Role

AgentSkeleton is a scaffold project for:

- new project initialization
- legacy project AI migration
- agent collaboration file generation
- repository structure inspection

## Expected Outcome

Claude Code should be able to work effectively in this repository without requiring a second, separate project system.

See `docs/host-integration.md` for the shared host-side integration flow.
Read current repository truth through:

- `docs/current-capabilities.md`
- `docs/features/README.md`
- `docs/README.md`
