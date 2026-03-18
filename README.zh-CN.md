# AgentSkeleton

语言版本：

- [English](README.md)
- [中文](README.zh-CN.md)

AgentSkeleton 是一个 AI-first 的文档架构引导工具，用于通过对话帮助团队建立 AI 友好的项目结构。

它围绕一个简单模型展开：

- 文档蓝图是基础资产
- CLI 是产品核心
- Codex 与 Claude Code 是主要协作宿主

它主要服务于两类场景：

- 引导新项目建立清晰的 AI 友好文档结构
- 将历史项目重塑为更适合 AI 工具与人类协作维护的文档结构

CLI 不替代大模型。它负责组织引导流程、沉淀结构化上下文，并告诉模型下一步应当生成什么文档。

## 它是什么

AgentSkeleton 不是业务系统，而是一套文档架构引导系统。

它的目标，是在不改动业务代码的前提下，通过引导式对话帮助用户建立 AI 友好的仓库文档。

## 项目目标

- 提供一套可复用的文档蓝图，作为基础资产。
- 交付一个 CLI，作为产品主要入口。
- 同时支持 Codex 风格 agent 工作流与 Claude Code 工作流。
- 在两种 agent 模式之间尽量共享核心结构，只保留最小必要差异。
- 同时支持新项目文档初始化和历史项目文档化重塑。
- 保持产品聚焦在引导和文档，而不是业务代码生成。

## 核心原则

- AI-first 协作方式
- 在不同 agent 模式之间共享结构
- 仓库规则显式化
- 先有稳定默认值，再谈重度定制
- 既支持新项目，也支持历史项目文档重塑
- 新项目可以采用推荐结构；老项目应先文档化并尊重现有结构

## 第一阶段的非目标

- 直接生成业务代码。
- 深度语言级项目脚手架。
- 在核心流程尚未验证前就引入复杂插件体系。

## MVP 范围

第一阶段重点是先把定义和结构建立起来。

- 引导新仓库建立文档结构。
- 为已有仓库提供文档化重塑指导。
- 生成协作文件，例如 `README.md`、`AGENTS.md` 和 `CLAUDE.md`。
- 输出结构化问题流，告诉模型下一步应补哪些文档。

## 支持的 Agent 模式

AgentSkeleton 默认支持：

- Codex / agent mode
- Claude / Claude Code mode

支持方式采用：

- 一套共享核心结构
- 一套共享文档蓝图基础
- 在确有必要时补充少量 agent 专属说明文件

这样可以控制维护成本，避免演化成两套互相分裂的项目体系。

## 产品模型

AgentSkeleton 的设计目标，是与大模型协作，而不是替代大模型。

默认前提是：

- CLI 负责引导对话并组织上下文
- Codex 或 Claude Code 负责真正撰写文档草案
- 人类负责定义目标、约束和验收标准

这条原则会直接影响产品设计和仓库协作方式。完整基线见 `docs/principles.zh-CN.md`。

## 仓库结构

```text
.
├── AGENTS.md
├── CLAUDE.md
├── README.md
├── cmd/
├── docs/
├── internal/
└── templates/
```

## 规划中的 CLI 方向

CLI 预期会成为用户的主要入口。首批命令方向包括：

- `init-docs`：引导新项目建立 AI 友好的文档结构
- `reshape-docs`：引导已有项目进行文档化重塑
- `plan`：总结当前应有哪些文档
- `next`：输出下一轮对话应追问的问题
- `response`：校验/评估模型输出，并可将合法答案写回上下文
- `prompt`：基于上下文生成初始提示或修复提示

请参考 `docs/agent-prompts.zh-CN.md`，了解宿主模型如何消费这些输出并进行重试。

当前仓库已经包含第一版最小 CLI 骨架，位于 `cmd/agentskeleton`，输出协议定义见 `docs/cli-contract.zh-CN.md`。

## 当前状态

仓库目前处于初始定义阶段。第一次推送应当至少建立：

- 核心文档
- 基础目录结构
- 命名约定
- agent 支持策略
- 面向文档引导的 CLI 方向
- 初版 CLI 协议与可运行命令骨架

## 贡献方向

在早期阶段，所有决策都应优先选择清晰、显式、默认稳定的方案，而不是过早追求高度灵活性。

## 提交约定

提交信息应带上 Jira 风格编号。

推荐格式：

- `docs [AG-001]: align product positioning and blueprint strategy`
- `feat [AG-001]: add initial documentation guidance flow`
- `fix [AG-001]: correct document planning output`
