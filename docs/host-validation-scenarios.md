# Host Validation Scenarios

## Purpose

This document defines the shared validation scenarios for AgentSkeleton host integrations in `Codex` and `Claude Code`.

Its job is to verify that the official host artifacts actually deliver the intended zero-command workflow for `v0.1.0`.

## Validation Goal

A host integration passes validation when it can deliver the expected AgentSkeleton workflow behavior without requiring the user to manually manage CLI commands.

Validation should focus on:

- user-facing zero-command experience
- correct host action routing
- correct use of stable protocol fields
- correct continuation after context changes
- correct escalation behavior when automation should pause

## Scope

These scenarios should be run for:

- the official Codex integration
- the official Claude Code integration

The same scenario intent should apply to both hosts, even if the host-specific execution environment differs.

## Validation Rules

During validation:

- do not help the host by naming AgentSkeleton commands unless the scenario explicitly requires it
- evaluate whether the host chooses the correct action on its own
- treat command exposure to the user as a product regression unless explicitly requested
- verify behavior against the shared host specs, not against host-specific improvisation

## Scenario 1: New Project Initialization

### Goal

Verify that the host can activate AgentSkeleton and choose `init-docs` for a clearly new project.

### Setup

- repository has no AgentSkeleton context yet
- repository is empty or near-empty
- user asks for help starting a new project with AI-friendly documentation structure

### Expected Host Behavior

- detects that AgentSkeleton should be activated
- chooses `init-docs` without asking the user to select a command
- establishes context
- runs `plan`
- asks only the next necessary clarification question if drafting is blocked

### Failure Signals

- asks the user which CLI command to run
- chooses `reshape-docs`
- remains in broad exploratory questioning instead of moving into workflow

## Scenario 2: Legacy Repository Reshape

### Goal

Verify that the host can activate AgentSkeleton and choose `reshape-docs` for an existing repository.

### Setup

- repository has existing code and docs but no AgentSkeleton context
- user asks to reorganize or document the repo for agent-friendly collaboration

### Expected Host Behavior

- detects that AgentSkeleton should be activated
- chooses `reshape-docs`
- establishes context
- runs `plan`
- continues toward clarification or drafting based on returned workflow state

### Failure Signals

- chooses `init-docs` without justification
- asks the user to manually inspect or create context files
- exposes command-level workflow to the user

## Scenario 3: Continue Current Priority Draft

### Goal

Verify that the host uses `focus-doc` for drafting instead of improvising from `plan` alone.

### Setup

- repository already has AgentSkeleton context
- `plan` indicates a clear `current_priority`
- required context is sufficient for drafting
- user asks to continue the documentation work

### Expected Host Behavior

- runs `plan` or uses already-fresh workflow state
- chooses the current priority document
- runs `focus-doc`
- drafts using the returned package

### Failure Signals

- drafts directly from `plan` without `focus-doc`
- asks the user which document should be drafted next when `current_priority` is already clear
- exposes drafting package internals instead of translating them into user-facing progress

## Scenario 4: Blocking Clarification Before Drafting

### Goal

Verify that the host asks the minimum necessary clarification question when drafting is not yet reliable.

### Setup

- repository already has AgentSkeleton context
- current priority is not ready because required context is missing
- user asks to continue or draft the next doc

### Expected Host Behavior

- identifies the blocking missing context
- asks one high-value clarification question
- avoids broad exploratory questioning
- does not try to draft through missing facts unless placeholders are explicitly acceptable

### Failure Signals

- asks multiple unfocused questions at once
- drafts anyway while inventing facts
- asks the user to figure out the next command

## Scenario 5: Apply Structured Answers And Continue

### Goal

Verify that the host writes accepted structured answers back into context and continues from `post_apply_plan`.

### Setup

- repository already has AgentSkeleton context
- the user provides answers that can be normalized into the response envelope

### Expected Host Behavior

- normalizes the answers
- runs `response --apply`
- uses `post_apply_plan` as the default continuation source
- proceeds to the next blocking question, current priority, or latest-batch convergence work

### Failure Signals

- does not write accepted answers back into context
- re-runs `plan` reflexively without using `post_apply_plan`
- loses the workflow thread after apply

## Scenario 6: Review Candidate Convergence

### Goal

Verify that the host treats `review_candidates` as latest-batch convergence work instead of a permanent backlog.

### Setup

- repository already has generated docs
- newly applied answers affect shared repository meaning
- `review_candidates` appear in `plan` or `post_apply_plan`

### Expected Host Behavior

- recognizes that previously generated docs may need refresh
- treats the review set as temporary and change-batch scoped
- does not present it as a permanent review queue

### Failure Signals

- ignores clearly relevant review candidates
- treats review candidates as a long-lived backlog
- replaces current drafting work with unrelated review churn

## Scenario 7: Stale Draft Package Refresh

### Goal

Verify that the host refreshes stale `focus-doc` packages after context changes.

### Setup

- the host already obtained a `focus-doc` package
- new answers are applied, changing the active `change_batch_id`
- the host returns to the same draft task

### Expected Host Behavior

- recognizes the old package is stale
- re-runs `focus-doc`
- continues only from the refreshed package

### Failure Signals

- continues drafting from the stale package
- does not compare or react to `change_batch_id`
- silently produces work based on outdated context

## Scenario 8: Escalation On Unresolved State

### Goal

Verify that the host stops and escalates when automation should no longer continue.

### Setup

- response validation reaches `unresolved`
or
- project mode remains ambiguous after the available evidence

### Expected Host Behavior

- pauses automatic workflow progression
- explains the single blocking issue
- asks one targeted clarification question

### Failure Signals

- keeps guessing
- invents missing business facts
- dumps internal workflow details instead of a focused escalation

## Scenario 9: User Requests A Specific Document

### Goal

Verify that the host can safely honor a user override without abandoning workflow safety.

### Setup

- repository already has AgentSkeleton context
- user asks for a specific document that is not necessarily the current priority

### Expected Host Behavior

- runs `focus-doc --path <document>`
- checks readiness and missing context
- drafts only if it is safe or placeholder-driven drafting is acceptable

### Failure Signals

- ignores the user override without explanation
- drafts the requested doc without checking readiness
- treats any user override as permission to bypass workflow safety

## Scenario 10: Quiet Zero-Command Experience

### Goal

Verify the overall product feel, not just routing correctness.

### Setup

- run any of the previous scenarios in a realistic host conversation

### Expected Host Behavior

- keeps command names mostly hidden
- explains progress in task language
- makes the workflow feel guided but not mechanical

### Failure Signals

- repeatedly exposes raw commands
- asks the user to manage internal state
- makes the experience feel like CLI coaching instead of hosted workflow

## Pass Criteria For M5

`M5` should only be considered validated when:

- Codex passes the core scenarios
- Claude Code passes the core scenarios
- both hosts behave consistently with shared specs
- zero-command behavior is preserved across initialization, drafting, apply, convergence, and escalation

## Recommended Evidence To Capture

For each validation run, capture:

- the user prompt that triggered the scenario
- the host response sequence
- which AgentSkeleton actions were selected
- whether the host exposed commands to the user
- whether the host used the correct stable protocol fields
- final pass/fail notes

## Conclusion

Host integration validation is not just about whether commands can run.

It is about whether the host can turn the AgentSkeleton protocol core into a reliable zero-command product experience.
