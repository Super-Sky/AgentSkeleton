# Host Validation Report Template

## Purpose

This template is used to record validation results for AgentSkeleton host integrations in `Codex` and `Claude Code`.

Use one report per validation run or one section per scenario run, depending on how much detail is needed.

## Report Metadata

- Date:
- Validator:
- Host:
- Host integration artifact:
- Repository under test:
- AgentSkeleton version/commit:
- Validation scenario:
- Result: `pass` | `fail` | `partial`

## Scenario Summary

- Scenario name:
- Validation goal:
- Why this scenario matters:

## Setup

- Repository state before the run:
- Whether AgentSkeleton context already existed:
- Relevant files or project conditions:
- User prompt used to trigger the scenario:

## Expected Behavior

- Expected host action path:
- Expected user experience:
- Expected stable protocol fields used:

## Actual Behavior

- Host response sequence:
- AgentSkeleton actions selected by the host:
- Whether the host exposed command names to the user:
- Whether the host asked clarifying questions:
- Whether the host drafted or updated a document:
- Whether the host applied structured answers:
- Whether the host continued from `post_apply_plan`:
- Whether the host refreshed stale draft packages:

## Stable Contract Check

- Used `plan` as workflow snapshot: `yes` | `no` | `n/a`
- Used `focus-doc` before drafting: `yes` | `no` | `n/a`
- Used `response --apply` for accepted structured results: `yes` | `no` | `n/a`
- Used `post_apply_plan` for continuation: `yes` | `no` | `n/a`
- Treated `review_candidates` as latest-batch convergence work: `yes` | `no` | `n/a`
- Respected `change_batch_id` freshness rules: `yes` | `no` | `n/a`

## Zero-Command Experience Check

- User had to choose commands manually: `yes` | `no`
- User had to inspect workflow internals manually: `yes` | `no`
- Host progress felt task-oriented rather than command-oriented: `yes` | `no`
- Host asked only necessary clarification questions: `yes` | `no`

## Outcome

- Final result:
- What worked:
- What failed:
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue:

## Evidence

- Transcript excerpt:
- Files or outputs produced:
- Relevant CLI output snapshot:
- Notes about missing evidence:

## Follow-Up

- Required fix:
- Owner:
- Priority:
- Should this block `v0.1.0`: `yes` | `no`

## Short Example

```markdown
- Date: 2026-04-09
- Validator: internal
- Host: Codex
- Host integration artifact: skills/agentskeleton-codex/SKILL.md
- Repository under test: sample legacy repo
- AgentSkeleton version/commit: e4b0d27
- Validation scenario: Scenario 2 - Legacy Repository Reshape
- Result: pass

## Outcome

- Final result: pass
- What worked: Codex selected reshape-docs, created context, ran plan, and asked one focused follow-up question.
- What failed: none
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue: n/a
```
