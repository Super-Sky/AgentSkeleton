# CLI Contract

## Purpose

This document defines the initial contract for the AgentSkeleton CLI.

The CLI is not responsible for writing complete repository documents by itself. Its job is to:

- organize project context
- decide what documentation should exist
- provide structured next-step guidance
- support Codex and Claude Code as host environments

## Core Design

- The CLI is the product core.
- Host models such as Codex and Claude Code are expected to consume the CLI output.
- The CLI should return structured, stable, machine-friendly output.
- Human-readable formatting may exist, but structured output is the baseline contract.

## Context File

The CLI should maintain a project context file.

Recommended location:

- `.agentskeleton/context.yaml`

The file stores the current state of the documentation guidance process.

## Context Schema

Initial fields:

```yaml
version: v0.0.0
project:
  name: ""
  summary: ""
  mode: new | legacy
  domain: ""
  primary_users: []
  host: codex | claude-code
documentation:
  phase: discovery | planning | drafting | refining
  generated_docs: []
  missing_docs: []
  release_version: v0.0.0
structure:
  strategy: recommended | existing
  recommended_layout: internal/app
  current_layout_summary: ""
conversation:
  answered_questions: []
  open_questions: []
```

## Command Set

Initial command areas:

- `plan`
- `next`
- `init-docs`
- `reshape-docs`

This contract defines `plan` and `next` first.

## `plan`

### Purpose

Summarize the current project state and produce a documentation plan.

### Input

- context file
- optional CLI flags for mode, host, or release version

### Output Requirements

`plan` should return:

- current project mode
- current documentation phase
- known project facts
- unresolved information gaps
- recommended document list
- why each document is needed
- recommended next actions

### Structured Output Shape

```yaml
command: plan
project_mode: new | legacy
documentation_phase: discovery | planning | drafting | refining
known_facts:
  - key: project_name
    value: MallHub
  - key: primary_users
    value:
      - mall operators
      - merchants
missing_information:
  - deployment_shape
  - ownership_model
recommended_documents:
  - path: README.md
    purpose: repository entrypoint and summary
    status: required
  - path: docs/domain-overview.md
    purpose: shared domain language for humans and models
    status: required
next_actions:
  - ask about deployment shape
  - ask about document ownership
  - draft README.md
```

## `next`

### Purpose

Provide the next structured questions that should be asked in conversation.

### Input

- context file
- optional stage override

### Output Requirements

`next` should return:

- current documentation phase
- current conversation goal
- ordered list of next questions
- why each question matters
- which documents depend on each answer

### Structured Output Shape

```yaml
command: next
documentation_phase: discovery
conversation_goal: clarify repository documentation scope
questions:
  - id: project_summary
    prompt: What is the one-sentence summary of this project?
    reason: The summary anchors README and domain overview drafts.
    affects:
      - README.md
      - docs/domain-overview.md
  - id: project_mode
    prompt: Is this a new project or an existing repository being reshaped?
    reason: This determines whether the CLI recommends structure or documents an existing one.
    affects:
      - docs/architecture.md
      - docs/legacy-reshape-guide.md
```

## `init-docs`

### Purpose

Initialize a documentation guidance session for a new project.

### Expected Behavior

- create `.agentskeleton/context.yaml` if missing
- set `project.mode` to `new`
- set `structure.strategy` to `recommended`
- seed an initial document plan

## `reshape-docs`

### Purpose

Initialize a documentation reshaping session for an existing repository.

### Expected Behavior

- create `.agentskeleton/context.yaml` if missing
- set `project.mode` to `legacy`
- set `structure.strategy` to `existing`
- seed a structure-inventory-first plan

## Output Format Policy

The CLI should support both:

- structured output for host-model consumption
- readable output for direct human inspection

But structured output is the stable contract.

Recommended formats:

- `yaml`
- `json`

## Non-Goals

- generating final business code
- forcing a single repository structure on legacy projects
- replacing host-model reasoning with hard-coded document text
