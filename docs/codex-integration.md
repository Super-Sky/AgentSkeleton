# Codex Integration

## Purpose

This document defines the first official Codex integration artifact for AgentSkeleton.

For `v0.1.0`, the recommended Codex integration form is a repository-distributed skill that teaches Codex how to deliver the zero-command AgentSkeleton workflow.

## Integration Artifact

The current official Codex artifact lives at:

- `skills/agentskeleton-codex/SKILL.md`

Its purpose is to:

- activate AgentSkeleton at the right time
- route between `init-docs`, `reshape-docs`, `plan`, `focus-doc`, and `response --apply`
- keep the zero-command flow aligned with shared host specs
- avoid leaking CLI workflow management to the user

## Why A Skill First

For `v0.1.0`, a skill is the thinnest useful integration form because it:

- keeps the host layer lightweight
- preserves the CLI as the workflow engine
- avoids prematurely locking the product into a heavier plugin form
- lets Codex reuse the shared host adapter specs immediately

This does not rule out a future plugin.

It means the first official Codex integration should optimize for clarity, alignment, and speed to validation.

## Skill Responsibilities

The Codex skill should:

- detect whether AgentSkeleton should be activated
- choose `init-docs` versus `reshape-docs`
- choose `update` when existing documentation should be broadly refreshed from current repository truth
- use `plan` as the workflow snapshot
- use `focus-doc` as the drafting package source
- use `response --apply` for write-back
- continue from `post_apply_plan`
- refresh stale draft packages
- ask only targeted clarification questions

The skill should not:

- become a second workflow engine
- introduce host-only semantics that conflict with shared docs
- depend critically on experimental CLI surfaces

## Dependency Order

The skill should direct Codex to read these docs first:

1. `AGENTS.md`
2. `docs/host-action-mapping.md`
3. `docs/zero-command-flow.md`
4. `docs/protocol-stability.md`

These define the stable host-side behavior for `v0.1.0`.

## Enablement Model

The intended usage model is:

- the repository ships the Codex skill
- Codex loads the skill when the repository task fits AgentSkeleton
- the user remains in natural conversation
- AgentSkeleton commands stay behind the host layer

## v0.1.0 Success Criteria

The Codex integration is considered minimally successful for `v0.1.0` when:

- Codex can activate AgentSkeleton without user command selection
- Codex can initialize or resume the workflow state
- Codex can route reliably between clarification and drafting
- Codex can write accepted structured results back through `response --apply`
- Codex can continue from `post_apply_plan`
- Codex can refresh stale draft packages

## Next Steps

After this initial artifact exists, the next implementation steps should be:

1. validate the skill against at least one new-project flow
2. validate the skill against at least one legacy-reshape flow
3. decide whether Codex needs an additional plugin or wrapper form after the skill path is proven

## Conclusion

For `v0.1.0`, the official Codex integration should stay thin.

The skill is the host-facing behavior layer.

AgentSkeleton remains the workflow and protocol core.
