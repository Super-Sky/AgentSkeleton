# Codex Validation Report: Scenario 1

- Date: 2026-04-09
- Validator: pending
- Host: Codex
- Host integration artifact: `skills/agentskeleton-codex/SKILL.md`
- Repository under test: pending real repository
- AgentSkeleton version/commit: `4018373`
- Validation scenario: Scenario 1 - New Project Initialization
- Result: `partial`

## Scenario Summary

- Scenario name: New Project Initialization
- Validation goal: Verify that Codex can activate AgentSkeleton and choose `init-docs` for a clearly new project without requiring command-level user action.
- Why this scenario matters: This is the first zero-command entrypoint for greenfield repositories and one of the most important release-gating flows.

## Setup

- Repository state before the run: pending real run
- Whether AgentSkeleton context already existed: expected `no`
- Relevant files or project conditions: empty or near-empty repository
- User prompt used to trigger the scenario:

```text
Help me start this project with an AI-friendly documentation structure.
```

## Expected Behavior

- Expected host action path:
  - detect AgentSkeleton should be activated
  - choose `init-docs`
  - establish context
  - run `plan`
  - ask one focused follow-up question only if drafting is blocked
- Expected user experience:
  - no manual command selection
  - no need to inspect `.agentskeleton` files
  - progress explained in task language
- Expected stable protocol fields used:
  - `plan.current_priority`
  - `plan.missing_information`
  - `response.post_apply_plan` if answers are applied during the flow

## Actual Behavior

- Host response sequence: pending real run
- AgentSkeleton actions selected by the host: pending real run
- Whether the host exposed command names to the user: pending real run
- Whether the host asked clarifying questions: pending real run
- Whether the host drafted or updated a document: pending real run
- Whether the host applied structured answers: pending real run
- Whether the host continued from `post_apply_plan`: pending real run
- Whether the host refreshed stale draft packages: `n/a` for initial run unless context changes mid-flow

## Stable Contract Check

- Used `plan` as workflow snapshot: `pending`
- Used `focus-doc` before drafting: `pending`
- Used `response --apply` for accepted structured results: `pending`
- Used `post_apply_plan` for continuation: `pending`
- Treated `review_candidates` as latest-batch convergence work: `n/a`
- Respected `change_batch_id` freshness rules: `n/a`

## Zero-Command Experience Check

- User had to choose commands manually: `pending`
- User had to inspect workflow internals manually: `pending`
- Host progress felt task-oriented rather than command-oriented: `pending`
- Host asked only necessary clarification questions: `pending`

## Outcome

- Final result: `partial`
- What worked:
  - the repository now contains a dedicated Codex integration skill and shared host specs to guide this flow
  - the CLI core commands required by this scenario already exist
  - local smoke validation proves the CLI path itself is runnable
- What failed:
  - no real Codex host execution evidence has been captured yet
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue:
  - validation evidence missing

## Evidence

- Transcript excerpt: pending real run
- Files or outputs produced: pending real run
- Relevant CLI output snapshot:
  - local CLI smoke path exists via `make smoke`
- Notes about missing evidence:
  - this report is pre-created as a release-tracking artifact and still requires a real Codex-hosted validation session

## Follow-Up

- Required fix:
  - run this scenario in a real Codex environment
  - fill in transcript evidence
  - mark pass/fail with concrete routing observations
- Owner: pending
- Priority: high
- Should this block `v0.1.0`: `yes`
