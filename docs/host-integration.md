# Host Integration Guide

## Purpose

This document explains how `AgentSkeleton` should be used inside `Codex` and `Claude Code` once it is installed in a project.

It is written for host agents, not for end users who should memorize CLI commands.

The core rule is simple:

- users express goals through conversation
- the host agent decides when to call `AgentSkeleton`
- the CLI emits structured guidance
- the model writes the actual documentation

## Host Positioning

The relationship should be understood as:

- the CLI is the engine
- Codex / Claude Code are the real user-facing entrypoints
- the user experiences conversation, not command syntax

Therefore, the host should not expect the user to manually remember:

- `init-docs`
- `plan`
- `focus-doc`
- `response --apply`

Those commands should be selected automatically by the host.

## Repository Discovery

When a host enters a project that uses AgentSkeleton, it should inspect:

1. `AGENTS.md`
2. `CLAUDE.md`
3. `<output-dir>/.agentskeleton/context.yaml`

Their roles are:

- `AGENTS.md`: shared working rules
- `CLAUDE.md`: Claude Code-specific additions
- `.agentskeleton/context.yaml`: current documentation guidance state

If `.agentskeleton/context.yaml` does not exist, the host should determine whether the user is:

- starting a new project
or
- reshaping documentation for an existing repository

Then it should automatically call:

- `init-docs`
or
- `reshape-docs`

## Recommended Main Flow

The current stable main flow is:

1. `init-docs` / `reshape-docs`
2. `plan`
3. `focus-doc`
4. host-model drafting
5. `response --apply`

This is the core protocol surface for the current MVP.

## Host Responsibilities Per Command

### 1. `init-docs` / `reshape-docs`

The host should automatically call these when:

- the project has no `.agentskeleton/context.yaml`
- the user clearly says this is a new project
- the user clearly says this is a legacy repository being reshaped

This should not be pushed back onto the user as a memorized CLI step.

### 2. `plan`

The host should treat `plan` as the current state snapshot.

The main fields to consume are:

- `recommended_documents`
- `current_priority`
- `review_candidates`

Where:

- `current_priority` means the next document to move forward
- `review_candidates` means temporary backtracking targets triggered by the latest change batch

### 3. `focus-doc`

The host should call `focus-doc` before drafting any document.

The main fields to consume are:

- `change_batch_id`
- `change_batch_inputs`
- `required_context`
- `missing_context`
- `available_context`
- `suggested_sections`
- `review_after_draft`

The host should treat this as the document drafting package.

### 4. Drafting

The host model is responsible for the actual draft.

It must follow these rules:

- do not invent unresolved facts
- use explicit placeholders when context is missing
- prefer reusable and publishable language
- do not turn the guidance system into a business code generator

### 5. `response --apply`

When the host receives a structured answer, it should call `response --apply` to write it back into context.

The main field to consume afterward is:

- `post_apply_plan`

The host should continue from this refreshed plan instead of always running `plan` again.

## Stale Draft Package Handling

The host must treat draft packages as expirable.

Rule:

- if the latest context batch moves beyond `focus-doc.change_batch_id`
- then the draft package is stale
- the host must call `focus-doc` again

The host should not continue drafting or converging documents from an outdated package.

## Review Model

`review_candidates` and `review_after_draft` are not permanent backlog markers.

They mean:

- which generated documents should be revisited for the latest change batch only
- which generated documents may need convergence after the current draft

The host should not cache them as long-lived reviewed/unreviewed flags.

## Codex Guidance

For Codex:

- read `AGENTS.md` first
- if the repository enables AgentSkeleton, follow the main flow
- treat `focus-doc` as the preferred drafting input rather than inferring structure directly from `plan`
- after `response --apply`, continue from `post_apply_plan`

## Claude Code Guidance

For Claude Code:

- read `CLAUDE.md` and `AGENTS.md` first
- use Claude-specific behavior only when necessary
- default to the shared main flow
- use `focus-doc` and `post_apply_plan` as primary progression signals

## Current Recommendation

At this stage, the host should primarily use:

- `init-docs` / `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

These are currently better treated as auxiliary or experimental:

- `workflow`
- `--write-plan-files`
- `--persist-trace`
- `--auto-repair`

## Conclusion

The goal of host integration is not to teach users a CLI.

It is to make the CLI almost invisible.

If the integration is successful, the user should mainly feel that:

- the agent asks better questions
- the agent organizes documentation more clearly
- the agent can keep moving forward while also converging affected documents

Instead of feeling that a new tool must be operated manually.
