# Host Action Mapping

## Purpose

This document defines how a host integration should map user conversation into AgentSkeleton actions.

It is written for `Codex` and `Claude Code` integrations that need to provide a zero-command user experience while keeping workflow semantics aligned with the AgentSkeleton core.

## Scope

This document defines:

- when a host should activate AgentSkeleton
- which CLI action the host should choose next
- when the host should ask the user for clarification
- when the host should revisit generated documents
- when the host should refresh a stale drafting package

It does not define:

- host-specific packaging details
- plugin installation details
- the internal implementation of host-side skills/plugins

## Core Rule

The host owns intent interpretation.

AgentSkeleton owns workflow state interpretation.

In practice:

- the host decides what the user is trying to do
- AgentSkeleton decides what the documentation workflow currently needs

## Repository Activation Rules

The host should consider AgentSkeleton active when at least one of the following is true:

- the repository already contains AgentSkeleton context under `.agentskeleton/` or `<output-dir>/.agentskeleton/`
- the repository contains explicit AgentSkeleton guidance in `AGENTS.md`, `CLAUDE.md`, or related docs
- the user explicitly asks to initialize, reshape, or organize repository documentation with AgentSkeleton

If none of the above is true, the host may still choose AgentSkeleton when:

- the user asks for AI-friendly repository documentation structure
- the user asks for repository documentation reshaping
- the user asks for agent collaboration files and repository guidance docs

## Initial Routing Rules

When AgentSkeleton is not yet initialized, the host should choose one of these entry actions:

- `init-docs` for a new project
- `reshape-docs` for an existing repository

The default decision rule is:

- use `init-docs` when the user clearly describes a new project, early repository setup, or greenfield initialization
- use `reshape-docs` when the user clearly describes an existing project, legacy repository, missing structure, or documentation cleanup

If the repository state and user message disagree, the host should ask one clarifying question before proceeding.

Examples:

- empty or near-empty repository plus "start this project" -> `init-docs`
- existing codebase plus "help me document and organize this repo" -> `reshape-docs`

## Main Flow Routing Rules

Once context exists, the host should route actions using the following order of preference:

1. resolve blocking missing context if drafting is not reliable
2. draft or update the current priority document
3. revisit review candidates triggered by the latest change batch
4. continue from `post_apply_plan`

The host should not keep asking open-ended questions if AgentSkeleton already indicates that drafting can proceed reliably.

## Action Selection Table

### Case 1: No Context Exists

Host action:

- choose `init-docs` or `reshape-docs`
- then run `plan`

Why:

- the host needs a first structured state snapshot before it can manage drafting

### Case 2: User Asks "What Should We Do Next?"

Host action:

- run `plan`

Why:

- `plan` is the workflow state snapshot and the correct entrypoint for prioritization

### Case 3: User Wants to Draft or Continue a Document

Host action:

- run `focus-doc`
- use the returned drafting package

Why:

- `focus-doc` is the authoritative drafting package for the current priority or requested document

### Case 4: Host Receives Structured Answers From Conversation

Host action:

- normalize answers into the response envelope
- run `response --apply`
- continue from `post_apply_plan`

Why:

- accepted answers must be written back into context before the workflow can advance correctly

### Case 5: Host Receives Invalid Structured Output

Host action:

- run the response validation path
- if decision is `retry`, use repair prompt flow
- if decision is `unresolved`, stop automatic progression and escalate to user clarification

Why:

- workflow state should not be corrupted by invalid structured writes

### Case 6: Existing Documents May Need Convergence

Host action:

- inspect `review_candidates` from `plan` or `post_apply_plan`
- choose review work only if the latest change batch materially affects generated docs

Why:

- review work is temporary and change-batch scoped, not a permanent backlog

## Clarification Rules

The host should ask the user a clarification question when one of the following is true:

- the host cannot confidently choose between `init-docs` and `reshape-docs`
- the current priority is not ready because required context is still missing
- the user request conflicts with current workflow state
- the user requests a specific document but the required context is insufficient
- the response path ends in `unresolved`

The host should avoid asking for clarification when:

- the CLI already provides enough context to keep drafting
- the user request clearly matches the current priority
- the needed follow-up is purely workflow progression and not a business ambiguity

## Priority Advancement Rules

The host should usually advance using `current_priority`.

The host should not override `current_priority` unless:

- the user explicitly requests a different document
- the host is handling review convergence for the latest change batch
- the requested override is still safe given the available context

When the user requests a specific document, the host should:

- run `focus-doc --path <document>`
- check readiness and missing context
- only proceed if the package is draftable or the user has accepted placeholder-driven drafting

## Review Candidate Rules

`review_candidates` should be treated as temporary convergence work.

The host should consider review work when:

- the latest resolved context changes shared repository meaning
- a recently generated document likely depends on the newly resolved answers
- `review_after_draft` points to affected generated documents

The host should not treat `review_candidates` as:

- a permanent task list
- a repository-wide review backlog
- a replacement for current priority drafting

Default rule:

- finish blocking clarification first
- then complete the current draft action
- then revisit review candidates from the latest change batch

## Stale Draft Package Rules

The host must treat a drafting package as stale when:

- the current repository context moves beyond the package's `change_batch_id`
- new answers have been applied since the package was produced
- the host has reason to believe `required_context`, `missing_context`, or `review_after_draft` may have changed

When a package is stale, the host should:

- stop drafting from the old package
- run `focus-doc` again
- continue only from the refreshed package

The host should never silently continue using an outdated package after new context has been written back.

## Continuation Rules After `response --apply`

After `response --apply`, the host should prefer the returned `post_apply_plan` instead of reflexively re-running `plan`.

The host should then choose among:

- ask the next blocking question if missing context still blocks drafting
- draft the new `current_priority`
- revisit `review_candidates` if the latest change batch makes convergence necessary

This keeps the host aligned with the freshest workflow state.

## Recommended Default Flow

The default host behavior should be:

1. detect whether AgentSkeleton is active or should be activated
2. initialize with `init-docs` or `reshape-docs` if needed
3. run `plan`
4. if clarification is blocking, ask the next high-value question
5. otherwise run `focus-doc`
6. draft or update the target document
7. normalize structured answers and run `response --apply`
8. continue from `post_apply_plan`
9. revisit review candidates only when the latest change batch requires convergence

## Escalation Rules

The host should escalate to the user instead of continuing automatically when:

- project mode is ambiguous
- the structured response is unresolved after retry budget
- the requested draft would rely on invented facts
- the user intent conflicts with the active workflow and cannot be safely reconciled

Escalation should be short and explicit.

The host should explain:

- what decision is blocked
- why the workflow cannot safely continue
- what single clarification is needed next

## Anti-Patterns

Host integrations should avoid:

- exposing AgentSkeleton command names unless the user explicitly asks
- repeatedly re-running `plan` when `post_apply_plan` already exists
- treating `review_candidates` as a permanent backlog
- drafting from stale `focus-doc` packages
- asking broad open-ended questions when the missing context is narrow and explicit
- hard-coding host behavior that conflicts with the shared protocol

## Minimum Behavioral Contract

For `v0.1.0`, a compliant host integration should:

- activate AgentSkeleton in the right repository situations
- choose `init-docs` versus `reshape-docs` predictably
- use `plan` for prioritization
- use `focus-doc` for drafting input
- use `response --apply` for accepted structured results
- continue from `post_apply_plan`
- respect `review_candidates` as latest-batch convergence work
- refresh stale drafting packages
- escalate instead of inventing missing facts

## Conclusion

The zero-command experience depends on hosts making consistent routing decisions.

This action-mapping layer is therefore part of the product surface, even though it lives outside the CLI core.
