# Real Project Validation Report: `sast_task_assistant` Legacy Stability

- Date: 2026-04-10
- Validator: Codex
- Host: local-cli
- Host integration artifact: `/tmp/agentskeleton-test`
- Repository under test: `/Users/mac/Desktop/Murphy/sast/sast_task_assistant`
- AgentSkeleton version/commit: `b178de2`
- Validation scenario: repeated legacy repository reshape via CLI on a real local project
- Result: `pass`

## Scenario Summary

- Scenario name: Repeated local legacy reshape on `sast_task_assistant`
- Validation goal: Verify that repeated AgentSkeleton CLI runs on the same real legacy repository produce stable and consistent workflow outputs across isolated output directories.
- Why this scenario matters: Stability across repeated runs is important before relying on the workflow during host-driven validation or release gating.

## Setup

- Repository state before the run:
  - git repository present
  - existing Go task orchestration service layout
  - repository already had a modified `makefile` before the validation run
- Whether AgentSkeleton context already existed: no
- Relevant files or project conditions:
  - `README.md`
  - `main.go`
  - top-level packages including `routers`, `controller`, `services`, `model`, `util`, `timer`, `common`, and `config`
- Isolated output directories used:
  - `/tmp/agentskeleton-sast-task-assistant-run1-C4efzM`
  - `/tmp/agentskeleton-sast-task-assistant-run2-hyYOYO`
- User prompt equivalent:

```text
Reshape this existing repository into an agent-friendly documented structure.
```

## Expected Behavior

- Expected host action path:
  - run `reshape-docs`
  - run `plan`
  - run `focus-doc`
  - apply structured answers for discovered repository structure
  - continue from the updated plan
- Expected user experience:
  - not applicable as a host UX test
- Expected stable protocol fields used:
  - `plan.current_priority`
  - `plan.missing_information`
  - `focus-doc.change_batch_id`
  - `response.post_apply_plan`
  - `review_candidates`

## Actual Behavior

- Host response sequence:
  - built `/tmp/agentskeleton-test`
  - executed two isolated validation runs
  - each run executed `reshape-docs`, `plan`, `focus-doc`, `response --apply`, and a final `plan`
- AgentSkeleton actions selected by the host:
  - identical across both runs
- Whether the host exposed command names to the user:
  - yes, this was an explicit CLI validation run
- Whether the host asked clarifying questions:
  - no host questions; structured answers were supplied directly from repository inspection
- Whether the host drafted or updated a document:
  - no final user-facing doc was drafted in this run
- Whether the host applied structured answers:
  - yes, in both runs
- Whether the host continued from `post_apply_plan`:
  - yes, and the same result was confirmed by a follow-up `plan`
- Whether the host refreshed stale draft packages:
  - not needed in this run

## Stability Check

- `plan-before.json` was identical across run 1 and run 2: `yes`
- `focus-before.json` was identical across run 1 and run 2: `yes`
- `response-apply.json` was identical across run 1 and run 2: `yes`
- `plan-after.json` was identical across run 1 and run 2: `yes`

## Stable Contract Check

- Used `plan` as workflow snapshot: `yes`
- Used `focus-doc` before drafting: `yes`
- Used `response --apply` for accepted structured results: `yes`
- Used `post_apply_plan` for continuation: `yes`
- Treated `review_candidates` as latest-batch convergence work: `yes`
- Respected `change_batch_id` freshness rules: `yes`

## Zero-Command Experience Check

- User had to choose commands manually: `yes`
- User had to inspect workflow internals manually: `no`
- Host progress felt task-oriented rather than command-oriented: `n/a`
- Host asked only necessary clarification questions: `n/a`

## Outcome

- Final result: `pass`
- What worked:
  - both isolated runs produced the same initial priority: `docs/legacy-structure-inventory.md`
  - both runs reported the same missing context before drafting
  - after applying `current_layout_summary` and `undocumented_directories`, both runs advanced to the same next priority: `docs/domain-overview.md`
  - both runs produced the same `review_candidates` set
  - no cross-run instability or non-deterministic output was observed in the tested path
- What failed:
  - no blocking failures in the repeated CLI legacy flow
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue:
  - none in this validation run

## Evidence

- Run 1 output directory:
  - `/tmp/agentskeleton-sast-task-assistant-run1-C4efzM`
- Run 2 output directory:
  - `/tmp/agentskeleton-sast-task-assistant-run2-hyYOYO`
- Important observed results:
  - initial current priority: `docs/legacy-structure-inventory.md`
  - post-apply current priority: `docs/domain-overview.md`
  - remaining missing information after apply: `active_release_docs_strategy`, `ownership_model`
- Notes about missing evidence:
  - this validates CLI stability on a real repository, but it is not a substitute for real Codex or Claude Code host validation

## Follow-Up

- Required fix:
  - continue with real host validation in Codex and Claude Code
- Owner: pending
- Priority: medium
- Should this block `v0.1.0`: `no`
