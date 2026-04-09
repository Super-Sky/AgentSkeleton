# Claude Code Validation Report: Scenario 2

- Date: 2026-04-09
- Validator: pending
- Host: Claude Code
- Host integration artifact: `skills/agentskeleton-claude/SKILL.md`
- Repository under test: pending real repository
- AgentSkeleton version/commit: `4018373`
- Validation scenario: Scenario 2 - Legacy Repository Reshape
- Result: `partial`

## Scenario Summary

- Scenario name: Legacy Repository Reshape
- Validation goal: Verify that Claude Code can activate AgentSkeleton and choose `reshape-docs` for an existing repository without exposing command-level workflow.
- Why this scenario matters: This is the critical zero-command entrypoint for legacy repositories and a core product promise for `v0.1.0`.

## Setup

- Repository state before the run: pending real run
- Whether AgentSkeleton context already existed: expected `no`
- Relevant files or project conditions: existing repository with code and uneven documentation
- User prompt used to trigger the scenario:

```text
Help me reorganize and document this existing repository so it works better with agents.
```

## Expected Behavior

- Expected host action path:
  - detect AgentSkeleton should be activated
  - choose `reshape-docs`
  - establish context
  - run `plan`
  - move into clarification or drafting according to returned workflow state
- Expected user experience:
  - no manual command selection
  - no requirement to inspect context internals
  - clear, host-driven progression
- Expected stable protocol fields used:
  - `plan.current_priority`
  - `plan.review_candidates`
  - `focus-doc.change_batch_id` if drafting begins

## Actual Behavior

- Host response sequence: pending real run
- AgentSkeleton actions selected by the host: pending real run
- Whether the host exposed command names to the user: pending real run
- Whether the host asked clarifying questions: pending real run
- Whether the host drafted or updated a document: pending real run
- Whether the host applied structured answers: pending real run
- Whether the host continued from `post_apply_plan`: pending real run
- Whether the host refreshed stale draft packages: pending real run if context changed during validation

## Stable Contract Check

- Used `plan` as workflow snapshot: `pending`
- Used `focus-doc` before drafting: `pending`
- Used `response --apply` for accepted structured results: `pending`
- Used `post_apply_plan` for continuation: `pending`
- Treated `review_candidates` as latest-batch convergence work: `pending`
- Respected `change_batch_id` freshness rules: `pending`

## Zero-Command Experience Check

- User had to choose commands manually: `pending`
- User had to inspect workflow internals manually: `pending`
- Host progress felt task-oriented rather than command-oriented: `pending`
- Host asked only necessary clarification questions: `pending`

## Outcome

- Final result: `partial`
- What worked:
  - the repository now includes a Claude Code skill aligned with the shared host adapter spec
  - the legacy reshape workflow is explicitly defined in the shared validation scenarios
  - the CLI core path required by this scenario already exists
- What failed:
  - no real Claude Code host execution evidence has been captured yet
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue:
  - validation evidence missing

## Evidence

- Transcript excerpt: pending real run
- Files or outputs produced: pending real run
- Relevant CLI output snapshot:
  - local CLI smoke coverage exists for the core command path, but not for host behavior
- Notes about missing evidence:
  - this report is pre-created as a release-tracking artifact and still requires a real Claude Code validation session

## Follow-Up

- Required fix:
  - run this scenario in a real Claude Code environment
  - capture the transcript and selected AgentSkeleton actions
  - mark pass/fail based on observed host routing behavior
- Owner: pending
- Priority: high
- Should this block `v0.1.0`: `yes`
