#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TMP_DIR=$(mktemp -d)
BIN_PATH="$TMP_DIR/agentskeleton"
PROJECT_DIR="$TMP_DIR/project"
OUTPUT_DIR="$TMP_DIR/output"
RESPONSE_FILE="$TMP_DIR/response.yaml"

cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cd "$ROOT_DIR"
go build -o "$BIN_PATH" ./cmd/agentskeleton

mkdir -p "$PROJECT_DIR" "$OUTPUT_DIR"

"$BIN_PATH" init-docs \
  --project "$PROJECT_DIR" \
  --output-dir "$OUTPUT_DIR" \
  --name "MallHub" \
  --format json > "$TMP_DIR/init.json"

grep -q '"command": "init-docs"' "$TMP_DIR/init.json"

"$BIN_PATH" plan \
  --project "$PROJECT_DIR" \
  --output-dir "$OUTPUT_DIR" \
  --format json > "$TMP_DIR/plan.json"

grep -q '"command": "plan"' "$TMP_DIR/plan.json"
grep -q '"current_priority"' "$TMP_DIR/plan.json"

"$BIN_PATH" focus-doc \
  --project "$PROJECT_DIR" \
  --output-dir "$OUTPUT_DIR" \
  --format json > "$TMP_DIR/focus.json"

grep -q '"command": "focus-doc"' "$TMP_DIR/focus.json"

cat > "$RESPONSE_FILE" <<'EOF'
status: ok
schema: question-answer-set-v1
data:
  project_summary: MallHub is an AI-friendly shopping mall operations platform.
errors: []
raw_text: ""
EOF

"$BIN_PATH" response \
  --file "$RESPONSE_FILE" \
  --project "$PROJECT_DIR" \
  --output-dir "$OUTPUT_DIR" \
  --apply \
  --question project_summary \
  --docs README.md \
  --format json > "$TMP_DIR/response.json"

grep -q '"context_updated": true' "$TMP_DIR/response.json"
grep -q '"post_apply_plan"' "$TMP_DIR/response.json"

"$BIN_PATH" version --format json > "$TMP_DIR/version.json"
grep -q '"command": "version"' "$TMP_DIR/version.json"
grep -q '"version": "dev"' "$TMP_DIR/version.json"

echo "smoke test passed"
