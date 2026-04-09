#!/bin/sh

set -eu

if [ "$#" -lt 3 ]; then
  echo "usage: $0 <host> <scenario-file-name> <title>" >&2
  echo "example: $0 codex codex-scenario-3 \"Codex Validation Report: Scenario 3\"" >&2
  exit 1
fi

HOST="$1"
FILE_NAME="$2"
TITLE="$3"

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
REPORT_DIR="$ROOT_DIR/docs/validation-reports"
TARGET="$REPORT_DIR/$FILE_NAME.md"

if [ -e "$TARGET" ]; then
  echo "report already exists: $TARGET" >&2
  exit 1
fi

cat > "$TARGET" <<EOF
# $TITLE

- Date: $(date -u +%Y-%m-%d)
- Validator: pending
- Host: $HOST
- Host integration artifact: pending
- Repository under test: pending
- AgentSkeleton version/commit: $(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || echo pending)
- Validation scenario: pending
- Result: \`partial\`

## Scenario Summary

- Scenario name: pending
- Validation goal: pending
- Why this scenario matters: pending

## Setup

- Repository state before the run: pending
- Whether AgentSkeleton context already existed: pending
- Relevant files or project conditions: pending
- User prompt used to trigger the scenario:

\`\`\`text
pending
\`\`\`

## Expected Behavior

- Expected host action path: pending
- Expected user experience: pending
- Expected stable protocol fields used: pending

## Actual Behavior

- Host response sequence: pending real run
- AgentSkeleton actions selected by the host: pending real run
- Whether the host exposed command names to the user: pending
- Whether the host asked clarifying questions: pending
- Whether the host drafted or updated a document: pending
- Whether the host applied structured answers: pending
- Whether the host continued from \`post_apply_plan\`: pending
- Whether the host refreshed stale draft packages: pending

## Stable Contract Check

- Used \`plan\` as workflow snapshot: \`pending\`
- Used \`focus-doc\` before drafting: \`pending\`
- Used \`response --apply\` for accepted structured results: \`pending\`
- Used \`post_apply_plan\` for continuation: \`pending\`
- Treated \`review_candidates\` as latest-batch convergence work: \`pending\`
- Respected \`change_batch_id\` freshness rules: \`pending\`

## Zero-Command Experience Check

- User had to choose commands manually: \`pending\`
- User had to inspect workflow internals manually: \`pending\`
- Host progress felt task-oriented rather than command-oriented: \`pending\`
- Host asked only necessary clarification questions: \`pending\`

## Outcome

- Final result: \`partial\`
- What worked: pending
- What failed: pending
- Was the failure a routing issue, protocol issue, UX issue, or documentation issue: pending

## Evidence

- Transcript excerpt: pending
- Files or outputs produced: pending
- Relevant CLI output snapshot: pending
- Notes about missing evidence: pending

## Follow-Up

- Required fix: pending
- Owner: pending
- Priority: pending
- Should this block \`v0.1.0\`: \`pending\`
EOF

echo "created $TARGET"
