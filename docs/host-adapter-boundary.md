# Host Adapter Boundary

## Purpose

This document defines the product boundary between `AgentSkeleton` and host-side integrations such as `Codex` skills/plugins and `Claude Code` skills/plugins.

The goal is to keep the product focused on documentation orchestration while still enabling a zero-command, conversation-driven user experience.

## Product Direction

The intended end state is:

- users work through conversation only
- hosts such as Codex and Claude Code decide when to invoke AgentSkeleton
- AgentSkeleton remains a structured engine instead of becoming a general-purpose conversational agent

In other words:

- the user experience should feel agent-native
- the product core should remain protocol-oriented

## Core Decision

`AgentSkeleton` should not expand into a full natural-language agent runtime.

Instead:

- host skills/plugins should interpret user intent
- host models should ask clarifying questions when needed
- AgentSkeleton should receive and return structured workflow state

This boundary prevents the project from becoming an overlapping second agent system.

## Why This Boundary Exists

### 1. Product Focus

If AgentSkeleton tries to own natural-language understanding, it stops being a documentation workflow engine and starts becoming an agent platform.

That would increase complexity across:

- prompting
- dialogue management
- ambiguity resolution
- long-running conversational state

### 2. Host Strength Reuse

Codex and Claude Code already provide strong capabilities for:

- natural-language understanding
- repository-aware reasoning
- clarification turns
- conversation memory

AgentSkeleton should reuse those strengths instead of rebuilding them.

### 3. Shared Protocol Stability

A structured core is easier to:

- test
- version
- validate
- keep aligned across multiple hosts

This matters more than embedding host-specific interaction logic into the core.

## Responsibility Split

### Host Skill / Plugin Responsibilities

The host layer should own:

- interpreting user requests from conversation
- deciding whether the repository should use AgentSkeleton
- choosing when to call `init-docs` or `reshape-docs`
- deciding when to call `plan`, `focus-doc`, and `response --apply`
- asking the user for missing business context
- deciding whether to advance current priority or revisit review targets
- presenting progress in a natural conversational form

### AgentSkeleton Responsibilities

AgentSkeleton should own:

- repository documentation state
- structured workflow progression
- document prioritization
- drafting package generation
- review target calculation
- response validation and apply behavior
- post-apply continuation state

### Shared Responsibility Boundary

The host may decide what the user is trying to do.

AgentSkeleton should decide what the current documentation workflow state means.

## What AgentSkeleton Should Avoid

To preserve the boundary, AgentSkeleton should avoid becoming responsible for:

- open-ended chat interaction
- generic natural-language command parsing
- free-form intent classification across unrelated tasks
- host-specific UX policy
- conversational memory outside repository state

If those concerns move into the core, the product boundary becomes unclear and host integrations become harder to maintain.

## Expected Host Integration Model

The preferred integration model is:

1. the user expresses a goal in conversation
2. the host skill/plugin maps that goal to AgentSkeleton actions
3. AgentSkeleton returns structured workflow guidance
4. the host model drafts or updates repository documents
5. the host writes the result back through `response --apply`
6. the host continues from the returned `post_apply_plan`

The user should not need to know the command names behind this loop.

## Minimum Protocol Expectation

For this boundary to work well, AgentSkeleton should continue optimizing for machine-consumable outputs.

Important outputs should stay explicit and stable, especially:

- current workflow state
- current priority
- review candidates
- drafting package context
- missing context
- change batch identity
- post-apply continuation state

Critical workflow meaning should not be hidden in long-form prose when a structured field can express it.

## Plugin and Skill Positioning

Codex and Claude Code integrations may take the form of:

- native plugins
- reusable skills
- host-side wrappers
- integrated command policies

The exact packaging may differ by host, but the boundary should remain the same:

- the host integration owns intent handling
- AgentSkeleton owns workflow orchestration

## Strategic Implication

This means the product should be described as:

- a documentation orchestration engine
- a CLI/protocol core
- a host-integrated workflow system

It should not be described as:

- a standalone conversational agent for arbitrary natural-language tasks

## Near-Term Guidance

In the current stage, the project should prioritize:

- clearer host adapter rules
- stronger structured output contracts
- stable workflow progression semantics
- thin Codex and Claude Code integrations

Natural-language UX quality should primarily improve through better host-side skills/plugins, not by turning the core into an independent agent runtime.

## Conclusion

The zero-command experience should happen at the host layer.

The durable product logic should remain in AgentSkeleton.

That split gives the project the best chance to be:

- easier to use
- easier to evolve
- easier to test
- easier to keep aligned across Codex and Claude Code
