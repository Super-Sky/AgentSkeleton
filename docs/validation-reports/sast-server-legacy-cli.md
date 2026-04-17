# Real Project Validation Report: `sast_server` Legacy CLI Flow

- Date: 2026-04-10
- Validator: Codex
- Host: local-cli
- Host integration artifact: `/tmp/agentskeleton-test`
- Repository under test: `/Users/mac/Desktop/Murphy/sast/sast_server`
- AgentSkeleton version/commit: `b178de2`
- Validation scenario: legacy repository reshape via CLI on a real local project
- Result: `pass`

## Scenario Summary

- Scenario name: Real local legacy reshape on `sast_server`
- Validation goal: Verify that AgentSkeleton can initialize legacy context, compute a sensible plan, focus the correct first document, accept structured answers, and advance to the next priority on a real repository.
- Why this scenario matters: It is the closest in-repo validation to the real `Claude Code` / `Codex` legacy-host flow without claiming host-side UX validation that did not happen.

## Setup

- Repository state before the run:
  - git repository present
  - existing Go service layout
  - untracked `.idea/` and `sast_server` binary already present in the target repo before the run
- Whether AgentSkeleton context already existed: no
- Relevant files or project conditions:
  - `README.md`
  - `main.go`
  - top-level packages such as `controller`, `services`, `router`, `utils`, `timer`, `common`, `config`, and `test`
- Output directory used:
  - `/tmp/agentskeleton-sast-server-nq47Rx`
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
  - continue from updated plan
- Expected user experience:
  - not applicable as a host UX test
- Expected stable protocol fields used:
  - `plan.current_priority`
  - `plan.missing_information`
  - `focus-doc.change_batch_id`
  - `response.post_apply_plan`

## Actual Behavior

- Host response sequence:
  - built `/tmp/agentskeleton-test`
  - ran `reshape-docs`
  - ran `plan`
  - ran `focus-doc`
  - applied structured answers with `response --apply`
  - re-ran `plan`
- AgentSkeleton actions selected by the host:
  - `reshape-docs`
  - `plan`
  - `focus-doc`
  - `response --apply`
  - `plan`
- Whether the host exposed command names to the user:
  - yes, this was an explicit CLI validation run
- Whether the host asked clarifying questions:
  - no host questions; structured answers were supplied directly from repository inspection
- Whether the host drafted or updated a document:
  - no final user-facing doc was drafted in this run
- Whether the host applied structured answers:
  - yes
- Whether the host continued from `post_apply_plan`:
  - yes, and the same result was confirmed by a follow-up `plan`
- Whether the host refreshed stale draft packages:
  - not needed in this run

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
  - legacy context initialized successfully
  - first priority selected as `docs/legacy-structure-inventory.md`
  - `focus-doc` correctly reported missing context before drafting
  - after applying `current_layout_summary` and `undocumented_directories`, the plan advanced to `docs/domain-overview.md`
  - `review_candidates` correctly flagged `docs/legacy-structure-inventory.md` for convergence after newly resolved structure context
- What failed:
  - no blocking failures in the CLI legacy flow
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue:
  - none in this validation run

## Evidence

- Relevant commands executed:

```bash
/tmp/agentskeleton-test reshape-docs --project /Users/mac/Desktop/Murphy/sast/sast_server --output-dir /tmp/agentskeleton-sast-server-nq47Rx --name sast_server --summary "Static application security testing service" --format json
/tmp/agentskeleton-test plan --project /Users/mac/Desktop/Murphy/sast/sast_server --output-dir /tmp/agentskeleton-sast-server-nq47Rx --format json
/tmp/agentskeleton-test focus-doc --project /Users/mac/Desktop/Murphy/sast/sast_server --output-dir /tmp/agentskeleton-sast-server-nq47Rx --format json
/tmp/agentskeleton-test response --file /tmp/agentskeleton-sast-server-nq47Rx/response.yaml --project /Users/mac/Desktop/Murphy/sast/sast_server --output-dir /tmp/agentskeleton-sast-server-nq47Rx --apply --docs docs/legacy-structure-inventory.md --format json
```

- Key observed outputs:
  - initial current priority: `docs/legacy-structure-inventory.md`
  - post-apply current priority: `docs/domain-overview.md`
  - post-apply missing information: `active_release_docs_strategy`, `ownership_model`
- Notes about missing evidence:
  - this validates the CLI against a real repository, but it is not a substitute for real Codex or Claude Code host validation

## Follow-Up

- Required fix:
  - continue with real host validation on Codex and Claude Code
- Owner: pending
- Priority: medium
- Should this block `v0.1.0`: `no`
