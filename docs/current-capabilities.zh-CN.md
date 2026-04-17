# 当前能力总览

本文档描述 AgentSkeleton 当前仓库真相。

它不是版本快照。

它首先要回答的问题是：

- 仓库现在到底支持什么

## 当前产品状态

AgentSkeleton 当前已经提供：

- 一个用于文档工作流编排的 Go CLI core
- 基于 `.agentskeleton` 的结构化上下文模型
- `plan`、drafting package、`response --apply` 工作流命令
- 一个用于基于仓库结构安全刷新上下文的 `update` 命令
- 面向零命令体验的宿主接入规范
- 初版 Codex 与 Claude Code skill 产物
- 验证场景、验证报告模板与种子报告
- CI、release-build 骨架以及本地 smoke 测试工具

## 当前稳定工作流核心

当前收窄后的稳定工作流核心是：

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

这些命令构成了 `v0.1.0` 发布路径下预期的公开宿主契约。

当前仓库还包含：

- `update`

`update` 的定位是已初始化上下文的刷新辅助命令。它很有用，但当前应视为辅助命令，而不是最窄的宿主契约。

## 当前宿主接入状态

仓库当前已经包含：

- `skills/agentskeleton-codex/SKILL.md`
- `skills/agentskeleton-claude/SKILL.md`
- 共享宿主规则文档：
  - `host-action-mapping.zh-CN.md`
  - `zero-command-flow.zh-CN.md`
  - `protocol-stability.zh-CN.md`

这意味着宿主接入模型已经定义并完成骨架化。

但这并不意味着真实宿主验证已经完成。

## 当前验证状态

仓库当前已经有：

- 验证场景
- 验证报告模板
- 种子报告文件
- 本地 CLI smoke 验证证据
- 两个真实本地项目的 CLI 验证证据：
  - `sast_server`
  - `sast_task_assistant`

当前仍缺：

- 真实 Codex 宿主运行证据
- 真实 Claude Code 宿主运行证据

## 当前发布状态

仓库当前已经有：

- CI workflow
- tag 触发的 release build workflow
- CLI `version` 命令
- 本地 `build` / `test` / `smoke` / `release-build` 目标
- 已知限制与 changelog 发布骨架

这意味着发布加固已经开始。

但这并不意味着 `v0.1.0` 已经可以切出。

## 当前阻塞项

当前主要发布阻塞项是：

- 获取真实 Codex 验证证据
- 获取真实 Claude Code 验证证据
- 根据真实宿主行为调整 skill
- 在 GitHub Actions 中验证 tag release 行为

## 阅读指引

接下来建议阅读：

1. `features/README.zh-CN.md`
2. `host-action-mapping.zh-CN.md`
3. `zero-command-flow.zh-CN.md`
4. `protocol-stability.zh-CN.md`

只有在需要计划细节或版本历史时，再进入版本路径文档。
