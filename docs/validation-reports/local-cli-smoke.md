# Local CLI Smoke Validation Report

- Date: 2026-04-09
- Validator: Codex
- Host: local-cli
- Host integration artifact: `scripts/smoke_test.sh`
- Repository under test: `AgentSkeleton`
- AgentSkeleton version/commit: `7314ae9`
- Validation scenario: Local CLI smoke path
- Result: `pass`

## Scenario Summary

- Scenario name: Local CLI smoke path
- Validation goal: Verify that the core CLI release path is runnable end-to-end without host-layer help.
- Why this scenario matters: It proves the protocol core can initialize context, plan, focus a document, apply accepted structured answers, and report version metadata.

## Setup

- Repository state before the run: current workspace checkout
- Whether AgentSkeleton context already existed: no, the smoke run creates an isolated temporary project and output directory
- Relevant files or project conditions:
  - `Makefile`
  - `scripts/smoke_test.sh`
- User prompt used to trigger the scenario:

```text
Run the local smoke path for the AgentSkeleton CLI.
```

## Expected Behavior

- Expected host action path:
  - build the CLI
  - run `init-docs`
  - run `plan`
  - run `focus-doc`
  - run `response --apply`
  - run `version`
- Expected user experience:
  - not applicable as a host UX test
- Expected stable protocol fields used:
  - `plan.current_priority`
  - `focus-doc.change_batch_id`
  - `response.post_apply_plan`
  - `version.version`

## Actual Behavior

- Host response sequence:
  - `make smoke`
  - `smoke test passed`
- AgentSkeleton actions selected by the host:
  - local smoke script invoked `init-docs`, `plan`, `focus-doc`, `response --apply`, and `version`
- Whether the host exposed command names to the user:
  - yes, this is a CLI-only validation path
- Whether the host asked clarifying questions:
  - no
- Whether the host drafted or updated a document:
  - no repository document drafting was attempted; the focus was protocol execution
- Whether the host applied structured answers:
  - yes
- Whether the host continued from `post_apply_plan`:
  - yes, the response output was checked for `post_apply_plan`
- Whether the host refreshed stale draft packages:
  - no, not part of this smoke scenario

## Stable Contract Check

- Used `plan` as workflow snapshot: `yes`
- Used `focus-doc` before drafting: `yes`
- Used `response --apply` for accepted structured results: `yes`
- Used `post_apply_plan` for continuation: `yes`
- Treated `review_candidates` as latest-batch convergence work: `n/a`
- Respected `change_batch_id` freshness rules: `n/a`

## Zero-Command Experience Check

- User had to choose commands manually: `yes`
- User had to inspect workflow internals manually: `no`
- Host progress felt task-oriented rather than command-oriented: `n/a`
- Host asked only necessary clarification questions: `n/a`

## Outcome

- Final result: `pass`
- What worked:
  - CLI build succeeded
  - smoke script completed successfully
  - `response --apply` produced `post_apply_plan`
  - `version --format json` returned valid output
- What failed:
  - this is not a host UX validation and does not prove zero-command behavior
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue:
  - no blocking failure in the local CLI path

## Evidence

- Transcript excerpt:

```text
$ make smoke
sh scripts/smoke_test.sh
smoke test passed
```

- Files or outputs produced:
  - temporary CLI binary
  - temporary context/output directories inside the smoke script
- Relevant CLI output snapshot:

```json
{
  "command": "version",
  "version": "dev"
}
```

- Notes about missing evidence:
  - this is a protocol-core validation, not a real Codex or Claude Code host validation

## Follow-Up

- Required fix:
  - continue running real host validation scenarios for release gating
- Owner: pending
- Priority: medium
- Should this block `v0.1.0`: `no`
