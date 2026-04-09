# Protocol Stability

## Purpose

This document defines which parts of the AgentSkeleton CLI contract should be treated as stable for `v0.1.0`.

The goal is to let host integrations depend on a predictable protocol surface without freezing every internal implementation detail.

## Stability Model

For `v0.1.0`, the protocol should be divided into three categories:

- stable
- provisional
- experimental

These categories apply to:

- commands
- structured output fields
- workflow semantics

## Stable Means

When a protocol element is marked stable, host integrations may rely on it as part of the public `v0.1.0` contract.

Stable elements should not change incompatibly within `v0.1.x` unless:

- the current behavior is broken
- the change is clearly documented
- migration impact is low and explicit

## Provisional Means

When a protocol element is marked provisional, it is expected to exist in `v0.1.0` but may still be refined in later minor releases.

Hosts may use provisional elements, but should avoid depending on undocumented edge behavior.

## Experimental Means

When a protocol element is marked experimental, it should not be treated as part of the durable public contract.

Hosts may use it for internal testing, but public-facing integrations should not make it a critical dependency.

## Stable Commands For v0.1.0

The following commands should be considered stable:

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

These commands define the core zero-command workflow.

## Provisional Commands For v0.1.0

The following commands should be considered provisional:

- `next`
- `prompt`
- `response` without `--apply`
- `workflow`

These are useful, but they should not be treated as the narrowest stable host contract yet.

## Stable Workflow Semantics

The following workflow semantics should be considered stable:

- the host uses `plan` as the workflow state snapshot
- the host uses `focus-doc` as the drafting package source of truth
- accepted structured answers are written back through `response --apply`
- `post_apply_plan` is the preferred continuation source after apply
- `review_candidates` are scoped to the latest change batch only
- `change_batch_id` is used to detect stale drafting packages
- missing required context should block reliable drafting unless placeholders are explicit

## Provisional Workflow Semantics

The following semantics should be treated as provisional:

- the exact prompting strategy used by host repair flows
- the exact priority between convergence and override drafting in edge cases
- the detailed layout of bundled `workflow` output
- the exact wording style of host-facing guidance blocks

## Experimental Workflow Semantics

The following should be treated as experimental:

- advanced autonomous repair loops beyond the current bounded retry model
- broad file materialization policies outside the current supported plan-writing behavior
- any host-specific workflow shortcuts that are not defined in shared specs

## Stable Output Fields

### `plan`

The following `plan` fields should be considered stable:

- `command`
- `project_mode`
- `documentation_phase`
- `known_facts`
- `missing_information`
- `recommended_documents`
- `current_priority`
- `review_candidates`
- `next_actions`

Within `current_priority`, the following should be considered stable:

- `path`
- `purpose`
- `required_context`
- `missing_context`
- `ready`
- `reason`

### `focus-doc`

The following `focus-doc` fields should be considered stable:

- `command`
- `path`
- `purpose`
- `ready`
- `change_batch_id`
- `change_batch_inputs`
- `required_context`
- `missing_context`
- `available_context`
- `suggested_sections`
- `review_after_draft`
- `next_actions`

### `response --apply`

The following `response` output fields should be considered stable when `--apply` is used:

- `command`
- `result`
- `context_updated`
- `post_apply_plan`

Within `result`, the following should be considered stable:

- `decision`
- `validation_errors`

## Provisional Output Fields

The following structured elements should be treated as provisional:

- detailed prompt bodies from `prompt`
- bundled `workflow` response structure beyond core fields already mirrored elsewhere
- trace file layout and naming details
- exact formatting of human-readable sections that duplicate structured meaning

## Stable Host Assumptions

Host integrations may assume that:

- `plan` is safe to use as a current workflow snapshot
- `focus-doc` is the required drafting package before drafting a document
- `response --apply` is the correct write-back path for accepted structured results
- `post_apply_plan` can be used instead of immediately re-running `plan`
- stale draft packages must be refreshed when `change_batch_id` moves

## Host Cautions

Host integrations should not assume that:

- every convenience command is part of the narrow stable core
- every free-text explanation block is protocol-stable
- undocumented edge ordering is guaranteed
- host-specific prompting behavior is portable across integrations unless shared docs say so

## Change Policy For v0.1.x

During the `v0.1.x` series:

- stable elements may be extended compatibly
- provisional elements may be refined or narrowed
- experimental elements may change without strong compatibility guarantees

If a stable element must change incompatibly, the repository should:

- update this document
- update host integration docs
- document the change in the changelog

## Guidance For Implementers

If a host integration needs the narrowest dependable contract, it should center its behavior on:

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`
- the stable change-batch and continuation semantics

That is the intended public protocol core for `v0.1.0`.

## Conclusion

Protocol stability should protect host integrations from accidental drift while still leaving room to improve the product.

For `v0.1.0`, stability should concentrate on the zero-command workflow core rather than on every convenience surface.
