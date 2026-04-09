# Zero-Command Flow

## Purpose

This document defines the target zero-command workflow for AgentSkeleton inside `Codex` and `Claude Code`.

It describes the user-facing experience that `v0.1.0` should deliver and the minimum host behavior required to make that experience consistent.

## Experience Goal

The user should experience AgentSkeleton as a conversation-native workflow.

The user should not need to:

- memorize AgentSkeleton command names
- decide which CLI step to run next
- manually manage workflow state
- manually detect stale drafting packages
- manually maintain convergence between generated documents

The host integration should absorb those responsibilities.

## Core User Story

The default zero-command story is:

1. the user explains a documentation goal in conversation
2. the host detects that AgentSkeleton should be used
3. the host initializes or resumes AgentSkeleton automatically
4. the host asks only the highest-value clarification questions
5. the host drafts or updates the next document using AgentSkeleton context
6. the host writes accepted structured progress back into workflow state
7. the host continues the loop until the current documentation objective is satisfied

## Scope

This document covers:

- the main zero-command loop
- initialization behavior
- drafting behavior
- clarification behavior
- convergence behavior
- continuation after context updates

It does not cover:

- plugin packaging
- installation instructions
- internal skill implementation details

## Flow States

The host should reason about the zero-command flow using these states:

- activation
- initialization
- planning
- clarification
- drafting
- applying
- convergence
- escalation

These are host behavior states, not necessarily one-to-one CLI commands.

## State 1: Activation

### Trigger

The user enters a repository conversation and expresses a documentation or repository-structure goal.

### Host Responsibilities

- detect whether AgentSkeleton is already active in the repository
- determine whether AgentSkeleton should be activated for this task
- decide whether the request fits AgentSkeleton or should stay outside its scope

### Expected User Experience

The user simply describes the goal.

The host does not ask the user to choose a command.

## State 2: Initialization

### Trigger

AgentSkeleton is not yet initialized for the repository.

### Host Responsibilities

- determine whether the repository should use `init-docs` or `reshape-docs`
- ask one clarification question if the mode is ambiguous
- run the chosen initialization action
- establish the first workflow context

### Expected User Experience

The user experiences this as a setup step handled by the host.

The host may ask a single high-value question such as whether the repository is new or existing, but should not turn initialization into a command tutorial.

## State 3: Planning

### Trigger

Context exists and the host needs the current workflow snapshot.

### Host Responsibilities

- run `plan`
- inspect `current_priority`
- inspect `missing_information`
- inspect `review_candidates`
- decide whether the next best action is clarification, drafting, or convergence

### Expected User Experience

The host presents a clear next step rather than exposing raw planning mechanics.

## State 4: Clarification

### Trigger

Drafting is not yet reliable because required context is missing or routing is ambiguous.

### Host Responsibilities

- ask the next highest-value question only
- keep the question narrow and tied to workflow needs
- avoid broad exploratory questioning when the missing context is specific
- normalize accepted answers into the response envelope

### Expected User Experience

The host asks focused follow-up questions that feel necessary and contextual.

The user should not feel like they are filling out a form mechanically.

## State 5: Drafting

### Trigger

The current priority or requested document is ready to draft, or ready enough with explicit placeholders.

### Host Responsibilities

- run `focus-doc`
- consume the drafting package
- use `required_context`, `missing_context`, `available_context`, and `suggested_sections`
- draft or update the target document without inventing missing facts

### Expected User Experience

The host should naturally continue the work rather than asking the user which command to run or what data structure to inspect.

## State 6: Applying

### Trigger

The host has accepted structured answers or completed a meaningful workflow update that should be written back into context.

### Host Responsibilities

- normalize structured answers into the response envelope
- run `response --apply`
- consume `post_apply_plan`
- continue from the returned workflow state

### Expected User Experience

The user experiences forward progress without needing to manually synchronize internal state.

## State 7: Convergence

### Trigger

The latest change batch affects already-generated documents.

### Host Responsibilities

- inspect `review_candidates` or `review_after_draft`
- decide whether convergence work should happen now or immediately after the current draft
- treat convergence as temporary latest-batch work, not as a permanent backlog

### Expected User Experience

The host may say that a previous document should be refreshed because new information changed shared meaning.

The user should not need to manually track affected docs.

## State 8: Escalation

### Trigger

The workflow cannot safely continue automatically.

### Host Responsibilities

- stop autonomous progression
- explain the single blocking issue
- ask for one targeted clarification

### Escalation Examples

- project mode is ambiguous
- response validation ends in `unresolved`
- the requested draft would require invented facts
- the requested task conflicts with the active workflow state

### Expected User Experience

Escalation should feel like a careful pause, not like a workflow failure dump.

## Canonical Main Loop

The default zero-command loop should be:

1. activate AgentSkeleton when the repository and task fit
2. initialize with `init-docs` or `reshape-docs` if needed
3. run `plan`
4. if drafting is blocked, ask the highest-value clarification question
5. apply accepted structured answers through `response --apply`
6. continue from `post_apply_plan`
7. when drafting is ready, run `focus-doc`
8. draft the target document
9. apply new structured progress
10. revisit latest-batch convergence only when required

This loop should continue until the current documentation goal is satisfied.

## Drafting Package Freshness

The host must treat a `focus-doc` package as temporary.

The package should be refreshed when:

- `change_batch_id` is no longer current
- new answers have been applied since the package was generated
- the host switched tasks and later returned to the same draft

The host should not continue drafting from a stale package.

## Conversation Style Expectations

The host should:

- speak in task-level language
- ask the minimum necessary number of questions
- explain next steps in user terms
- avoid referencing internal protocol details unless the user explicitly asks

The host should not:

- expose raw CLI sequences by default
- require the user to inspect workflow fields manually
- ask the user to manage change batches or drafting package freshness

## Success Criteria For v0.1.0

The zero-command flow is considered real when both `Codex` and `Claude Code` can do all of the following:

- activate AgentSkeleton without user command selection
- initialize or resume context automatically
- route between clarification and drafting correctly
- use `focus-doc` as the drafting source of truth
- apply structured progress with `response --apply`
- continue from `post_apply_plan`
- revisit convergence work only when the latest change batch requires it
- pause and escalate instead of inventing missing facts

## Non-Goals

For `v0.1.0`, the zero-command flow does not need to include:

- full autonomous repository completion without user feedback
- advanced multi-document parallel planning
- broad customization of conversation policy
- deep marketplace/plugin ecosystem support

## Conclusion

The zero-command flow is the user-facing product promise for `v0.1.0`.

The CLI core enables that promise, but the host integration must make it feel natural, quiet, and reliable.
