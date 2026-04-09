# Known Limitations

## Purpose

This document lists the known limitations for the current public AgentSkeleton release track.

For `v0.1.0`, these limitations should be treated as explicit scope boundaries rather than accidental omissions.

## Current Limitations

- The official host integrations are skill-based first, not full plugin products.
- Zero-command behavior is defined and scaffolded, but still depends on host validation runs for proof.
- The narrow stable host contract is intentionally small and centered on:
  - `init-docs`
  - `reshape-docs`
  - `plan`
  - `focus-doc`
  - `response --apply`
- `workflow`, prompt formatting details, and broader convenience surfaces are not yet the recommended stable host dependency path.
- Release packaging currently focuses on CLI build artifacts, not marketplace-grade host distribution.
- The repository does not yet claim broad cross-platform support beyond the launch build matrix.
- Host validation evidence must still be recorded before `v0.1.0` can be treated as release-ready.
- The project remains focused on repository documentation orchestration, not business code generation.

## Non-Goals At This Stage

- heavy plugin ecosystems
- rich UI management surfaces
- advanced autonomous repair loops beyond the current bounded retry path
- deep project-language scaffolding
- full autonomous repository completion without user clarification

## Guidance

If a workflow depends on behavior outside the stable protocol core, it should be treated as provisional or experimental until explicitly documented otherwise.
