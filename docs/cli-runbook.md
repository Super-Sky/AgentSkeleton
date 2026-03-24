# CLI Runbook

This runbook describes the minimum end-to-end loop for documentation guidance.

## New Project Flow

1. Initialize context:

```bash
agentskeleton init-docs --project /path/to/project --output-dir /path/to/output --name MallHub
```

This creates two kinds of files:

```text
/path/to/output/.agentskeleton/    # AgentSkeleton process state
/path/to/output/README.md          # final docs and skills belong here later
/path/to/output/docs/...
```

2. Generate plan:

```bash
agentskeleton plan --project /path/to/project --output-dir /path/to/output --format yaml
```

If `--context` is omitted, it resolves to:

```text
/path/to/output/.agentskeleton/context.yaml
```

3. Generate initial host-model prompt:

```bash
agentskeleton prompt --project /path/to/project --output-dir /path/to/output --mode initial --format yaml
```

4. Validate and apply a host-model response:

```bash
agentskeleton response \
  --file /path/to/host-response.yaml \
  --project /path/to/project \
  --output-dir /path/to/output \
  --attempt 0 \
  --apply \
  --docs /path/to/output/README.md,/path/to/output/docs/domain-overview.md
```

If the response `data` contains multiple fields, all of them are applied by default. Use `--question <id>` only when you want to apply a single field.

5. Continue with next questions:

```bash
agentskeleton next --project /path/to/project --output-dir /path/to/output --format yaml
```

## One-Command Flow

You can run one bundled step with:

```bash
agentskeleton workflow --project /path/to/project --output-dir /path/to/output --format yaml
```

If you also want the currently supported planned documents to be materialized into the output directory:

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --write-plan-files \
  --format yaml
```

By default this writes missing files only. Add `--overwrite` only when you intentionally want to replace existing generated docs.
When files are created or already present, the workflow also writes their generated state back into `<output-dir>/.agentskeleton/context.yaml`.

If you want to keep a structured snapshot of the whole step for auditing or later replay:

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --persist-trace \
  --format yaml
```

This writes a trace file under `<output-dir>/.agentskeleton/traces/` using the current documentation phase in the filename, and returns `trace_path` in the CLI output.

If you already have a host-model response:

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --apply \
  --write-plan-files \
  --docs /path/to/output/README.md,/path/to/output/docs/domain-overview.md \
  --format yaml
```

## Cleanup Model

- `<output-dir>/.agentskeleton` contains only AgentSkeleton process artifacts.
- Final docs and skills under `<output-dir>/...` are user-facing deliverables.
- If the user deletes `<output-dir>/.agentskeleton`, the final docs remain intact and AgentSkeleton no longer retains process state for that run.

## Retry Loop

If `agentskeleton response` returns `decision: retry`:

1. Build repair prompt:

```bash
agentskeleton prompt \
  --project /path/to/project \
  --output-dir /path/to/output \
  --mode repair \
  --errors "missing required field: project_summary" \
  --format yaml
```

2. Ask host model to repair structure only.
3. Validate again with incremented `--attempt`.
4. If decision becomes `unresolved`, stop applying and escalate to manual review.

You can bundle this into `workflow` as well:

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --auto-repair \
  --format yaml
```

When the response is retryable, the output includes an `auto_repair` block with the next attempt number, validation errors, a repair prompt, and instructions for the host-model loop.
