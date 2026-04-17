# Claude Code 接入

## 目的

本文档用于定义 AgentSkeleton 的第一个官方 Claude Code 接入产物。

对于 `v0.1.0`，推荐的 Claude Code 接入形态是一个随仓库分发的 skill，用来教会 Claude Code 如何提供 AgentSkeleton 的零命令工作流体验。

## 接入产物

当前官方 Claude Code 产物位于：

- `skills/agentskeleton-claude/SKILL.md`

它的职责是：

- 在合适时机激活 AgentSkeleton
- 在 `init-docs`、`reshape-docs`、`plan`、`focus-doc` 与 `response --apply` 之间做路由
- 让零命令主链路与共享宿主规范保持一致
- 不把 CLI 工作流管理暴露给用户

## 为什么优先做 Skill

对于 `v0.1.0`，skill 是当前最薄、最实用的接入形态，因为它：

- 让宿主层保持轻量
- 保持 CLI 仍是工作流引擎
- 避免过早把产品锁死到更重的 plugin 形态
- 让 Claude Code 可以立即复用共享 host adapter 规范

这并不排斥未来继续演进成 plugin。

它只是意味着：第一版官方 Claude Code 接入应优先追求清晰、对齐和验证速度。

## Skill 职责

Claude Code skill 应负责：

- 判断是否应激活 AgentSkeleton
- 选择 `init-docs` 还是 `reshape-docs`
- 当任务是基于当前仓库事实刷新已有文档时，选择 `update`
- 使用 `plan` 作为工作流快照
- 使用 `focus-doc` 作为起草包来源
- 使用 `response --apply` 做写回
- 从 `post_apply_plan` 继续推进
- 刷新过期的 drafting package
- 只提出有针对性的澄清问题

Claude Code skill 不应负责：

- 变成第二套工作流引擎
- 引入与共享文档冲突的宿主专属语义
- 关键依赖实验性 CLI 表面

## 依赖顺序

skill 应指导 Claude Code 优先读取以下文档：

1. `CLAUDE.md`
2. `AGENTS.md`
3. `docs/host-action-mapping.md`
4. `docs/zero-command-flow.md`
5. `docs/protocol-stability.md`

这些文档共同定义了 `v0.1.0` 的稳定宿主行为。

## 启用模型

预期使用方式是：

- 仓库随代码一起分发 Claude Code skill
- 当任务适合 AgentSkeleton 时，Claude Code 加载这个 skill
- 用户仍然只停留在自然对话中
- AgentSkeleton 命令隐藏在宿主层背后

## `v0.1.0` 成功标准

当满足以下条件时，Claude Code 接入可以被视为满足 `v0.1.0` 的最小要求：

- Claude Code 无需用户选命令即可激活 AgentSkeleton
- Claude Code 可以初始化或恢复工作流状态
- Claude Code 能在澄清与起草之间稳定路由
- Claude Code 能通过 `response --apply` 写回被接受的结构化结果
- Claude Code 能从 `post_apply_plan` 继续推进
- Claude Code 能刷新过期 drafting package

## 下一步

当这份初版产物存在后，建议继续做：

1. 用一条新项目链路验证 skill
2. 用一条老项目 reshape 链路验证 skill
3. 在 skill 路径被证明有效后，再决定 Claude Code 是否还需要 plugin 或 wrapper 形态

## 结论

对 `v0.1.0` 来说，官方 Claude Code 接入应保持足够轻。

skill 是宿主行为层。

AgentSkeleton 继续作为工作流和协议核心存在。
