---
name: agentskeleton-claude
description: Use this skill when working in a repository that uses AgentSkeleton and the goal is to provide a zero-command Claude Code experience. It tells Claude Code when to activate AgentSkeleton, when to use init-docs or reshape-docs, how to route between plan, focus-doc, and response --apply, and how to continue safely without exposing CLI workflow to the user.
---

# AgentSkeleton for Claude Code

Use this skill when:

- the repository already uses AgentSkeleton
- the user wants AI-friendly repository documentation structure
- the user wants to initialize or reshape documentation through conversation
- the task is to continue an AgentSkeleton-guided documentation workflow in Claude Code

Do not use this skill when:

- the task is unrelated to repository documentation workflow
- the user is asking for general coding help outside AgentSkeleton scope

## Core Role

Your job is to make AgentSkeleton feel invisible to the user.

That means:

- interpret user intent in conversation
- choose the right AgentSkeleton command path
- keep workflow state aligned with the CLI
- avoid asking the user to remember or choose commands

Do not turn AgentSkeleton into a general-purpose chat agent.

## Read First

Before acting, read these repository docs in this order:

1. `CLAUDE.md`
2. `AGENTS.md`
3. `docs/host-action-mapping.md`
4. `docs/zero-command-flow.md`
5. `docs/protocol-stability.md`

Read these when needed:

- `docs/host-integration.md` for the shared host integration model
- `docs/host-adapter-boundary.md` for product boundary decisions
- `docs/cli-runbook.md` for command examples

## Default Workflow

Use this default loop:

1. detect whether AgentSkeleton is already active
2. if no context exists, choose `init-docs` or `reshape-docs`
3. run `plan`
4. if missing context blocks drafting, ask the highest-value clarification question
5. otherwise run `focus-doc`
6. draft or update the target document
7. normalize accepted structured answers
8. run `response --apply`
9. continue from `post_apply_plan`
10. revisit `review_candidates` only when the latest change batch requires convergence

## Activation Rules

Treat AgentSkeleton as active when:

- `.agentskeleton/` or `<output-dir>/.agentskeleton/` exists
- the repository docs explicitly describe AgentSkeleton usage
- the user explicitly asks to use AgentSkeleton

You may also activate it when the user asks to:

- structure repository documentation
- reshape a legacy repository into an AI-friendly documented state
- create agent collaboration files and supporting repository docs

## Routing Rules

- Use `init-docs` for clearly new or greenfield projects.
- Use `reshape-docs` for clearly existing or legacy repositories.
- Use `plan` for workflow state and prioritization.
- Use `focus-doc` before drafting any document.
- Use `response --apply` to write accepted structured answers back into context.
- Use `post_apply_plan` as the default continuation source after apply.

Do not default to `workflow` for the narrow stable host path in `v0.1.0`.

## Claude-Specific Guidance

Keep Claude-specific behavior thin.

- Prefer shared repository rules first.
- Only use Claude-specific deviations when they are clearly necessary.
- Do not fork the workflow semantics away from the shared AgentSkeleton host specs.

## Clarification Rules

Ask one narrow clarification question when:

- project mode is ambiguous
- the current priority is blocked by missing required context
- the user requests a specific document that is not yet safe to draft
- response validation ends in `unresolved`

Avoid broad exploratory questioning when the missing context is explicit.

## Freshness Rules

Treat `focus-doc` packages as temporary.

Refresh the package when:

- `change_batch_id` is stale
- new answers have been applied
- you return to the same draft after context changed

Never continue drafting from a stale package.

## Convergence Rules

Treat `review_candidates` and `review_after_draft` as temporary latest-batch convergence work.

Do not treat them as a permanent backlog.

Default order:

1. unblock clarification
2. complete the current draft step
3. revisit affected generated docs

## Stability Rules For v0.1.0

For the narrow stable Claude Code integration path, depend primarily on:

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

Do not depend critically on:

- `workflow`
- prompt text formatting details
- undocumented edge ordering

## Conversation Style

Speak in task-level language.

Do:

- explain the next action in user terms
- ask focused follow-up questions
- keep progress moving quietly

Do not:

- expose command names unless the user explicitly asks
- ask the user to manage context files or change batches
- invent missing facts

## Escalate Instead Of Guessing

Pause and ask the user when:

- the repository mode is unclear
- the workflow cannot continue without a business fact
- the requested action conflicts with current workflow state

Make the escalation short:

- say what is blocked
- say why it matters
- ask for one concrete clarification
