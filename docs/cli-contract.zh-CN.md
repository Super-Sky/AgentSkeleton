# CLI Contract

## 目的

本文档定义 AgentSkeleton CLI 的初始协议。

CLI 不负责独立写完整仓库文档。它的职责是：

- 组织项目上下文
- 判断项目应该具备哪些文档
- 输出结构化的下一步引导
- 支持 Codex 与 Claude Code 作为宿主环境

## 核心设计

- CLI 是产品核心。
- Codex 与 Claude Code 这类宿主模型负责消费 CLI 输出。
- CLI 应返回结构化、稳定、机器友好的结果。
- 可以提供适合人类阅读的格式，但结构化输出才是基础协议。

## 上下文文件

CLI 应维护一个项目上下文文件。

推荐位置：

- `.agentskeleton/context.yaml`

这个文件用于存储当前文档引导过程的状态。

## 上下文结构

初始字段建议如下：

```yaml
version: v0.0.0
project:
  name: ""
  summary: ""
  mode: new | legacy
  domain: ""
  primary_users: []
  host: codex | claude-code
documentation:
  phase: discovery | planning | drafting | refining
  generated_docs: []
  missing_docs: []
  release_version: v0.0.0
structure:
  strategy: recommended | existing
  recommended_layout: internal/app
  current_layout_summary: ""
conversation:
  answered_questions: []
  open_questions: []
```

## 命令集合

第一批命令区域：

- `plan`
- `next`
- `init-docs`
- `reshape-docs`

这个协议优先定义 `plan` 和 `next`。

## `plan`

### 目的

汇总当前项目状态，并输出文档规划。

### 输入

- 上下文文件
- 可选 CLI 参数，例如模式、宿主类型、版本号

### 输出要求

`plan` 应返回：

- 当前项目模式
- 当前文档阶段
- 已知项目信息
- 未解决的信息缺口
- 建议文档列表
- 每个文档存在的目的
- 推荐下一步动作

### 结构化输出形态

```yaml
command: plan
project_mode: new | legacy
documentation_phase: discovery | planning | drafting | refining
known_facts:
  - key: project_name
    value: MallHub
  - key: primary_users
    value:
      - mall operators
      - merchants
missing_information:
  - deployment_shape
  - ownership_model
recommended_documents:
  - path: README.md
    purpose: repository entrypoint and summary
    status: required
  - path: docs/domain-overview.md
    purpose: shared domain language for humans and models
    status: required
next_actions:
  - ask about deployment shape
  - ask about document ownership
  - draft README.md
```

## `next`

### 目的

输出下一轮对话中应追问的结构化问题。

### 输入

- 上下文文件
- 可选阶段覆盖参数

### 输出要求

`next` 应返回：

- 当前文档阶段
- 当前对话目标
- 有顺序的下一批问题
- 每个问题的重要性说明
- 每个答案会影响哪些文档

### 结构化输出形态

```yaml
command: next
documentation_phase: discovery
conversation_goal: clarify repository documentation scope
questions:
  - id: project_summary
    prompt: What is the one-sentence summary of this project?
    reason: The summary anchors README and domain overview drafts.
    affects:
      - README.md
      - docs/domain-overview.md
  - id: project_mode
    prompt: Is this a new project or an existing repository being reshaped?
    reason: This determines whether the CLI recommends structure or documents an existing one.
    affects:
      - docs/architecture.md
      - docs/legacy-reshape-guide.md
```

## `init-docs`

### 目的

为新项目初始化一轮文档引导会话。

### 预期行为

- 若不存在，则创建 `.agentskeleton/context.yaml`
- 将 `project.mode` 设置为 `new`
- 将 `structure.strategy` 设置为 `recommended`
- 初始化一份基础文档计划

## `reshape-docs`

### 目的

为已有仓库初始化一轮文档化重塑会话。

### 预期行为

- 若不存在，则创建 `.agentskeleton/context.yaml`
- 将 `project.mode` 设置为 `legacy`
- 将 `structure.strategy` 设置为 `existing`
- 初始化一份“先盘点结构”的文档计划

## 输出格式策略

CLI 应同时支持：

- 供宿主模型消费的结构化输出
- 供人类直接查看的可读输出

但结构化输出是稳定协议。

推荐格式：

- `yaml`
- `json`

## 非目标

- 生成最终业务代码
- 强迫老项目采用单一仓库结构
- 用硬编码文案替代宿主模型的文档推理能力
