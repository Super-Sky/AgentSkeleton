# 宿主动作映射

## 目的

本文档用于定义宿主接入层如何把用户对话映射成 AgentSkeleton 动作。

它面向 `Codex` 与 `Claude Code` 的接入实现，目标是在保持与 AgentSkeleton core 工作流语义一致的前提下，提供零命令感知体验。

## 范围

本文档定义：

- 宿主何时应激活 AgentSkeleton
- 宿主下一步应选择哪类 CLI 动作
- 何时应向用户追问澄清
- 何时应回看已生成文档
- 何时应刷新已经过期的 drafting package

本文档不定义：

- 宿主侧的具体打包方式
- plugin 安装细节
- 宿主 skill/plugin 的内部实现方式

## 核心规则

宿主负责理解意图。

AgentSkeleton 负责解释工作流状态。

落到实际行为上就是：

- 宿主决定“用户想做什么”
- AgentSkeleton 决定“当前文档工作流下一步应做什么”

## 仓库激活规则

当满足以下任一条件时，宿主应认为 AgentSkeleton 已激活或应被激活：

- 仓库已存在 `.agentskeleton/` 或 `<output-dir>/.agentskeleton/` 上下文
- 仓库中的 `AGENTS.md`、`CLAUDE.md` 或相关文档显式提到了 AgentSkeleton
- 用户明确要求使用 AgentSkeleton 初始化、reshape 或组织仓库文档

即使以上条件都不满足，宿主在以下场景中也可以主动选择 AgentSkeleton：

- 用户要求建立 AI-friendly 的仓库文档结构
- 用户要求对仓库文档进行重塑
- 用户要求生成 agent 协作文档和仓库引导文档

## 初始路由规则

当 AgentSkeleton 尚未初始化时，宿主应在以下入口动作中二选一：

- 新项目使用 `init-docs`
- 现有仓库使用 `reshape-docs`

默认判断规则是：

- 当用户明显在描述一个新项目、早期初始化或 greenfield setup 时，使用 `init-docs`
- 当用户明显在描述一个已有项目、legacy 仓库、缺少结构或需要文档清理时，使用 `reshape-docs`

如果仓库状态与用户表达不一致，宿主应先追问一个澄清问题，再继续推进。

示例：

- 空仓库或几乎空仓库，且用户说“开始搭这个项目” -> `init-docs`
- 已有代码仓库，且用户说“帮我整理并文档化这个 repo” -> `reshape-docs`

## 主链路路由规则

一旦上下文已经存在，宿主应按以下优先顺序路由动作：

1. 如果起草不可靠，先解决阻塞性的缺失上下文
2. 起草或更新当前 priority 文档
3. 回看最近 change batch 触发的 review candidates
4. 从 `post_apply_plan` 继续推进

如果 AgentSkeleton 已经表明可以可靠起草，宿主不应继续停留在开放式追问阶段。

## 动作选择表

### 情况 1：上下文尚不存在

宿主动作：

- 选择 `init-docs` 或 `reshape-docs`
- 然后执行 `plan`

原因：

- 宿主需要先拿到第一份结构化状态快照，后续才能管理起草流程

### 情况 2：用户问“下一步做什么”

宿主动作：

- 执行 `plan`

原因：

- `plan` 是工作流状态快照，也是正确的优先级入口

### 情况 3：用户要开始起草或继续推进文档

宿主动作：

- 执行 `focus-doc`
- 使用返回的 drafting package

原因：

- `focus-doc` 是当前 priority 或指定文档的权威起草包

### 情况 3A：用户要求基于当前仓库事实刷新已有文档

宿主动作：

- 执行 `update`
- 然后从返回的 `post_update_plan` 继续推进

原因：

- 大面积文档刷新应先吸收那些无需再次确认、可以从仓库中安全推断的事实
- 这样可以避免对同一轮结构信息重复追问

### 情况 4：宿主从对话中拿到了结构化答案

宿主动作：

- 把答案规范化成 response envelope
- 执行 `response --apply`
- 从 `post_apply_plan` 继续推进

原因：

- 只有先把已接受答案写回上下文，工作流状态才能正确前进

### 情况 5：宿主拿到的结构化输出无效

宿主动作：

- 走 response 校验路径
- 如果结果是 `retry`，进入 repair prompt 流
- 如果结果是 `unresolved`，停止自动推进并升级为向用户澄清

原因：

- 不能让无效结构化写入污染工作流状态

### 情况 6：已有文档可能需要回看收敛

宿主动作：

- 检查 `plan` 或 `post_apply_plan` 中的 `review_candidates`
- 只有当最新 change batch 实质影响了已生成文档时，才选择回看

原因：

- review work 是临时的、绑定 change batch 的收敛动作，不是永久 backlog

## 澄清规则

当满足以下任一条件时，宿主应向用户追问澄清：

- 宿主无法高置信地区分 `init-docs` 与 `reshape-docs`
- 当前 priority 因 required context 缺失而尚未 ready
- 用户请求与当前工作流状态冲突
- 用户指定了某份文档，但 required context 不足
- response 路径最终进入 `unresolved`

宿主应避免在以下场景中继续追问：

- CLI 已提供足够上下文，可以继续起草
- 用户请求与当前 priority 明确一致
- 所需后续只是工作流推进，而不是业务歧义

## Priority 推进规则

默认情况下，宿主应按 `current_priority` 推进。

宿主不应覆盖 `current_priority`，除非：

- 用户明确要求另一份文档
- 当前正在处理最新 change batch 触发的 review convergence
- 这个覆盖在现有上下文下仍然安全

当用户要求某份指定文档时，宿主应：

- 执行 `focus-doc --path <document>`
- 检查 readiness 和 missing context
- 只有在可起草，或用户已接受“带显式占位起草”时才继续

当用户要求的是“大面积文档刷新”而不是单篇文档时，宿主应：

- 在已有 context 的前提下优先执行 `update`
- 使用刷新后的 plan 再决定下一份 draft 目标

## Review Candidate 规则

`review_candidates` 应被视为临时的收敛工作。

宿主应在以下场景考虑回看：

- 最新已解决上下文改变了仓库的共享含义
- 最近生成的文档很可能依赖这些新答案
- `review_after_draft` 指向了受影响的已生成文档

宿主不应把 `review_candidates` 视为：

- 永久任务列表
- 仓库级 review backlog
- 对 `current_priority` 起草的替代

默认规则是：

- 先完成阻塞性澄清
- 再完成当前 draft 动作
- 然后回看最近 change batch 里的 review candidates

## Stale Draft Package 规则

当满足以下任一条件时，宿主必须把一个 drafting package 视为过期：

- 当前仓库上下文已经推进到超出该 package 的 `change_batch_id`
- package 生成后又有新答案被 apply
- 宿主有理由相信 `required_context`、`missing_context` 或 `review_after_draft` 已经变化

当 package 过期时，宿主应：

- 停止基于旧 package 起草
- 重新执行 `focus-doc`
- 只从刷新后的 package 继续

在新上下文已经写回后，宿主不应静默继续使用旧 package。

## `response --apply` 之后的继续推进规则

在执行 `response --apply` 之后，宿主应优先消费返回的 `post_apply_plan`，而不是机械地再次执行 `plan`。

然后宿主应在以下动作中选择：

- 如果缺失上下文仍阻塞起草，就追问下一个高价值问题
- 起草新的 `current_priority`
- 如果最新 change batch 需要收敛，则回看 `review_candidates`

这样可以让宿主始终跟随最新工作流状态。

## 推荐默认流程

推荐的宿主默认行为是：

1. 判断 AgentSkeleton 是否已启用或应被启用
2. 如有需要，执行 `init-docs` 或 `reshape-docs`
3. 如果当前任务是大面积文档刷新且 context 已存在，则先执行 `update`
4. 否则执行 `plan`
5. 如果澄清阻塞，就提下一个高价值问题
6. 否则执行 `focus-doc`
7. 起草或更新目标文档
8. 规范化结构化答案并执行 `response --apply`
9. 从 `post_apply_plan` 继续推进
10. 只有当最新 change batch 需要收敛时，才回看 review candidates

## 升级规则

当满足以下任一条件时，宿主应停止自动推进并升级给用户，而不是继续猜测：

- project mode 含糊不清
- 结构化 response 在 retry budget 后仍是 unresolved
- 指定起草会依赖编造事实
- 用户意图与当前工作流冲突，且无法安全调和

升级时应简短且明确。

宿主应说明：

- 当前被什么决策阻塞
- 为什么工作流无法安全继续
- 下一步只需要用户回答哪一个澄清点

## 反模式

宿主接入应避免：

- 在用户未明确要求时暴露 AgentSkeleton 命令名
- 当已经有 `post_apply_plan` 时反复重新执行 `plan`
- 把 `review_candidates` 当作永久 backlog
- 基于过期的 `focus-doc` package 起草
- 在 missing context 已经狭窄明确时继续提出宽泛问题
- 硬编码与共享协议冲突的宿主行为

## 最小行为契约

在 `v0.1.0` 中，一个合格的宿主接入至少应做到：

- 在正确的仓库场景中激活 AgentSkeleton
- 可预测地选择 `init-docs` 或 `reshape-docs`
- 用 `plan` 做优先级判断
- 用 `focus-doc` 作为起草输入
- 用 `response --apply` 写回被接受的结构化结果
- 从 `post_apply_plan` 继续推进
- 把 `review_candidates` 视为最新 change batch 的收敛工作
- 刷新过期的 drafting package
- 在缺失事实时升级，而不是编造

## 结论

零命令感知体验能否成立，很大程度取决于宿主是否做出一致的动作路由决策。

因此，这层 action mapping 虽然位于 CLI core 之外，但它本身也属于产品表面的一部分。
