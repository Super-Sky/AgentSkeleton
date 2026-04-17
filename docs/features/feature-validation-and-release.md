# Feature: Validation And Release

## Purpose

This feature area covers the current validation and release-hardening assets.

## Current Scope

The repository currently provides:

- host validation scenarios
- validation report templates
- seeded validation reports
- local smoke testing
- CI and release build workflows
- known limitations and release-track changelog updates

## Current Artifacts

- `host-validation-scenarios.md`
- `host-validation-report-template.md`
- `validation-reports/README.md`
- `known-limitations.md`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `Makefile`
- `scripts/smoke_test.sh`

## Current Limitation

Release gating still depends on real Codex and Claude Code validation evidence plus a real tagged release run.
