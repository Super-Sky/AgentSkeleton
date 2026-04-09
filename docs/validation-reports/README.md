# Validation Reports

This directory stores recorded host validation runs for AgentSkeleton.

Recommended contents:

- one report per host/scenario run
- filenames that identify host and scenario clearly
- reports based on `docs/host-validation-report-template.md`

Suggested naming:

- `codex-scenario-1.md`
- `codex-scenario-2.md`
- `claude-scenario-1.md`
- `claude-scenario-2.md`

You can scaffold a new report with:

```bash
make new-validation-report HOST=codex FILE=codex-scenario-3 TITLE="Codex Validation Report: Scenario 3"
```

These reports are expected to provide the concrete validation evidence required before cutting `v0.1.0`.
