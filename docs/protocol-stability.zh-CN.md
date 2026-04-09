# 协议稳定性

## 目的

本文档用于定义 AgentSkeleton CLI 契约中哪些部分应在 `v0.1.0` 中被视为稳定能力。

目标是让宿主接入能够依赖一个可预测的协议表面，同时又不必冻结所有内部实现细节。

## 稳定性模型

对于 `v0.1.0`，协议应被划分为三类：

- stable
- provisional
- experimental

这些分类适用于：

- commands
- structured output fields
- workflow semantics

## Stable 的含义

当某个协议元素被标记为 stable 时，宿主接入可以把它当作 `v0.1.0` 的公开契约一部分来依赖。

在 `v0.1.x` 范围内，stable 元素不应发生不兼容变化，除非：

- 当前行为本身存在缺陷
- 变更被清楚记录
- 迁移影响低且明确

## Provisional 的含义

当某个协议元素被标记为 provisional 时，表示它会出现在 `v0.1.0` 中，但后续 minor 版本仍可能被进一步打磨。

宿主可以使用这类元素，但不应依赖那些未文档化的边界行为。

## Experimental 的含义

当某个协议元素被标记为 experimental 时，它不应被视为耐久的公开契约一部分。

宿主可以把它用于内部验证，但面向公开用户的接入不应把它作为关键依赖。

## `v0.1.0` 的稳定命令

以下命令应被视为 stable：

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

这些命令构成零命令主链路的核心。

## `v0.1.0` 的 provisional 命令

以下命令应被视为 provisional：

- `next`
- `prompt`
- `response`（不带 `--apply`）
- `workflow`

它们是有用能力，但当前还不应被视为最窄、最稳定的宿主契约。

## 稳定工作流语义

以下工作流语义应被视为 stable：

- 宿主使用 `plan` 作为工作流状态快照
- 宿主使用 `focus-doc` 作为起草包事实来源
- 被接受的结构化答案通过 `response --apply` 写回
- `post_apply_plan` 是 apply 之后优先使用的继续推进来源
- `review_candidates` 只绑定最近一次 change batch
- `change_batch_id` 用于识别过期 drafting package
- 当 required context 缺失时，除非显式使用占位，否则不能可靠起草

## Provisional 工作流语义

以下语义应视为 provisional：

- 宿主 repair flow 采用的具体 prompting 策略
- 在边界场景下 convergence 与 override drafting 的精确优先级
- 打包式 `workflow` 输出的详细布局
- 面向宿主的引导文本块的具体措辞

## Experimental 工作流语义

以下内容应视为 experimental：

- 超出当前有界重试模型之外的高级自动修复循环
- 超出当前 plan-writing 能力之外的广泛文件物化策略
- 未在共享规范中定义的宿主专属工作流捷径

## 稳定输出字段

### `plan`

以下 `plan` 字段应被视为 stable：

- `command`
- `project_mode`
- `documentation_phase`
- `known_facts`
- `missing_information`
- `recommended_documents`
- `current_priority`
- `review_candidates`
- `next_actions`

在 `current_priority` 内，以下字段应被视为 stable：

- `path`
- `purpose`
- `required_context`
- `missing_context`
- `ready`
- `reason`

### `focus-doc`

以下 `focus-doc` 字段应被视为 stable：

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

当使用 `--apply` 时，以下 `response` 输出字段应被视为 stable：

- `command`
- `result`
- `context_updated`
- `post_apply_plan`

在 `result` 内，以下字段应被视为 stable：

- `decision`
- `validation_errors`

## Provisional 输出字段

以下结构化元素应被视为 provisional：

- `prompt` 返回的详细 prompt body
- `workflow` 中那些未在其他核心命令中重复出现的详细响应结构
- trace 文件布局与命名细节
- 那些重复结构化语义的人类可读文本块的具体格式

## 稳定宿主假设

宿主接入可以稳定假设：

- `plan` 可作为当前工作流快照安全使用
- 在起草文档前，`focus-doc` 是必需的 drafting package
- `response --apply` 是已接受结构化结果的正确写回路径
- `post_apply_plan` 可替代立即再次执行 `plan`
- 当 `change_batch_id` 变化时，过期 draft package 必须刷新

## 宿主注意事项

宿主接入不应假设：

- 每个便利命令都属于最窄稳定核心
- 每个自由文本解释块都具备协议稳定性
- 未文档化的边界执行顺序是有保证的
- 未被共享文档明确说明的宿主 prompting 行为可以跨接入复用

## `v0.1.x` 变更策略

在 `v0.1.x` 范围内：

- stable 元素可以做兼容扩展
- provisional 元素可以被继续收敛或缩小
- experimental 元素可以在没有强兼容承诺的情况下调整

如果 stable 元素必须发生不兼容变化，仓库应：

- 更新本文档
- 更新宿主接入文档
- 在 changelog 中记录该变更

## 给实现者的建议

如果一个宿主接入只想依赖最窄、最可靠的公开契约，它应围绕以下内容构建：

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`
- 稳定的 change-batch 与 continuation 语义

这就是 `v0.1.0` 预期的公开协议核心。

## 结论

协议稳定性的作用，是防止宿主接入发生意外漂移，同时又为产品改进保留空间。

对 `v0.1.0` 来说，稳定性应集中在零命令主链路核心，而不是覆盖所有便利表面。
