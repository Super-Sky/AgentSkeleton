# CLI Runbook

This runbook describes the minimum end-to-end loop for documentation guidance.

## New Project Flow

1. Initialize context:

```bash
agentskeleton init-docs --name MallHub --context .agentskeleton/context.yaml
```

2. Generate plan:

```bash
agentskeleton plan --context .agentskeleton/context.yaml --format yaml
```

3. Generate initial host-model prompt:

```bash
agentskeleton prompt --context .agentskeleton/context.yaml --mode initial --format yaml
```

4. Validate and apply a host-model response:

```bash
agentskeleton response \
  --file /path/to/host-response.yaml \
  --context .agentskeleton/context.yaml \
  --attempt 0 \
  --apply \
  --question project_summary \
  --docs README.md,docs/domain-overview.md
```

5. Continue with next questions:

```bash
agentskeleton next --context .agentskeleton/context.yaml --format yaml
```

## One-Command Flow

You can run one bundled step with:

```bash
agentskeleton workflow --context .agentskeleton/context.yaml --format yaml
```

If you already have a host-model response:

```bash
agentskeleton workflow \
  --context .agentskeleton/context.yaml \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --apply \
  --question project_summary \
  --docs README.md,docs/domain-overview.md \
  --format yaml
```

## Retry Loop

If `agentskeleton response` returns `decision: retry`:

1. Build repair prompt:

```bash
agentskeleton prompt \
  --context .agentskeleton/context.yaml \
  --mode repair \
  --errors "missing required field: project_summary" \
  --format yaml
```

2. Ask host model to repair structure only.
3. Validate again with incremented `--attempt`.
4. If decision becomes `unresolved`, stop applying and escalate to manual review.
