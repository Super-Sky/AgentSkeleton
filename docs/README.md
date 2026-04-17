# Document Index

This directory stores AgentSkeleton product documents, workflow contracts, host integration rules, and release-track artifacts.

The repository now uses a four-layer document model:

- current truth
- current feature docs
- long-lived project docs
- version snapshots

AI and humans should not read these layers in random order.

## Default Reading Order

Read in this order unless a task explicitly needs historical release context:

1. `current-capabilities.md`
2. `features/README.md`
3. long-lived project docs
4. version snapshot docs only when needed

## Layer Definitions

### 1. Current Truth

Use:

- `current-capabilities.md`

Purpose:

- describe what the repository currently and concretely supports
- summarize the current release-track state
- avoid forcing readers to reconstruct the present from old version snapshots

### 2. Current Feature Docs

Use:

- `features/README.md`

Purpose:

- describe the current major capability areas
- act as the main feature-level entrypoint for the current repository truth
- avoid scattering current behavior across multiple release directories

### 3. Long-Lived Project Docs

Use:

- `principles.md`
- `architecture.md`
- `cli-contract.md`
- `host-adapter-boundary.md`
- `host-action-mapping.md`
- `zero-command-flow.md`
- `protocol-stability.md`

Purpose:

- define durable product positioning, workflow semantics, and host contracts
- outlive a single release

### 4. Version Snapshots

Use:

- `v0.1.0-gap-analysis.md`
- `v0.1.0-implementation-plan.md`

Purpose:

- capture release-specific planning and scope
- preserve snapshot context for one release track
- not replace the current truth docs

## Current Main Entry Files

- `current-capabilities.md`
- `features/README.md`
- `principles.md`
- `architecture.md`
- `cli-contract.md`
- `host-integration.md`
- `host-action-mapping.md`
- `zero-command-flow.md`
- `protocol-stability.md`

## Validation And Release Assets

- `host-validation-scenarios.md`
- `host-validation-report-template.md`
- `validation-reports/README.md`
- `known-limitations.md`

## Skill And Host Integration Docs

- `codex-integration.md`
- `claude-integration.md`
- `skill-sync.md`

## Rule

If a document describes the current product truth, it should live in:

- `current-capabilities.md`
or
- `features/*.md`

If a document only explains one release planning moment, it should not become the main current-truth entrypoint.
