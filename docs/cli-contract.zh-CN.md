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
- 宿主模型的返回结果在被复用前，也必须经过结构化协议校验。

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
- `response`
- `prompt`
- `workflow`

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

## `response`

### 目的

校验宿主模型返回、评估重试决策，并在满足条件时把合法答案写回上下文。

### 预期行为

- 从 YAML/JSON 解析返回包络
- 进行 schema 校验
- 评估决策（`accept`、`retry`、`unresolved`）
- 仅在 `accept` 且指定 `--apply` 时：
  - 更新 `.agentskeleton/context.yaml`
  - 将已回答问题从 `open_questions` 移除
  - 可选将文档标记为已生成

默认会批量写回 `data` 中的所有字段；如果提供 `--question`，则仅写回该字段。

## `prompt`

### 目的

基于当前上下文生成宿主模型提示词，包括初始提示和修复提示。

### 预期行为

- 读取 `.agentskeleton/context.yaml`
- 在 `initial` 模式下，为当前 `open_questions` 生成 schema 约束提示词
- 在 `repair` 模式下，附带校验错误并要求模型仅修复结构
- 以结构化包裹返回 prompt 文本，方便下游使用

## `workflow`

### 目的

执行一次完整引导步骤，打包 `plan`、`prompt`、`next`，并支持可选响应评估与安全写回。

### 预期行为

- 读取上下文
- 可选评估 `--response-file`
- 若 `--apply` 且决策为 `accept`，写回上下文
- 返回：
  - `plan` 输出
  - `prompt` 输出
  - `next` 输出
  - 可选响应评估结果

## 输出格式策略

CLI 应同时支持：

- 供宿主模型消费的结构化输出
- 供人类直接查看的可读输出

但结构化输出是稳定协议。

推荐格式：

- `yaml`
- `json`

## 宿主模型返回策略

仅有 CLI 输出还不够。产品必须默认假设宿主模型的返回可能存在格式错误、字段缺失或部分不符合 schema 的情况。

### 规则

任何宿主模型返回结果，在通过 schema 校验前，都不能进入下一步流程。

### 必要流程

1. CLI 输出结构化引导和期望返回结构。
2. 宿主模型返回一份候选结构化答案。
3. 产品按预期 schema 对结果进行校验。
4. 如果校验失败，产品返回一条有针对性的重试提示。
5. 重试持续进行，直到：
   - 返回结果变为合法，或
   - 重试额度耗尽
6. 如果重试额度耗尽，产品应回退到一个安全的“未解决状态”输出。

### 重试要求

- 重试时必须明确指出失败点。
- 重试应优先要求模型修复结构，而不是无条件整段重答。
- 校验错误最好能机器可读。
- 重试循环必须有上限。

推荐初始重试额度：

- 自动重试 `2` 次
- 若仍失败，则返回 `1` 次最终 fallback 结果

### 回退行为

如果多次重试后仍不合法，产品应：

- 将结果标记为 unresolved
- 在安全前提下保留原始模型输出以供检查
- 不要把无效结构写入 `.agentskeleton/context.yaml`
- 返回需要进一步澄清或人工介入的下一步

## 预期返回包络

产品后续应统一定义宿主模型返回包络，例如：

```yaml
status: ok | invalid | unresolved
schema: question-answer-set-v1
data: {}
errors: []
raw_text: ""
```

这个包络的意义，是让重试和后续处理都具备确定性。

## 非目标

- 生成最终业务代码
- 强迫老项目采用单一仓库结构
- 用硬编码文案替代宿主模型的文档推理能力
