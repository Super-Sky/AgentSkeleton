# 宿主接入指南

## 目的

本文档说明 `AgentSkeleton` 安装后，`Codex` 与 `Claude Code` 应如何在项目中使用它。

这份文档面向宿主 agent，而不是要求最终用户手工记忆 CLI 命令。

核心原则只有一条：

- 用户通过对话提出目标
- 宿主 agent 自动选择何时调用 `AgentSkeleton`
- CLI 负责输出结构化引导
- 大模型负责真正起草文档

## 宿主定位

`AgentSkeleton` 的关系应理解为：

- CLI 是引擎
- Codex / Claude Code 是用户实际使用入口
- 用户看到的是对话体验，不是命令细节

因此，宿主不应要求用户手工背诵：

- `init-docs`
- `plan`
- `focus-doc`
- `response --apply`

这些命令应由宿主根据上下文自动调用。

## 仓库发现方式

当宿主进入一个启用了 AgentSkeleton 的项目时，应优先检查：

1. `AGENTS.md`
2. `CLAUDE.md`
3. `<output-dir>/.agentskeleton/context.yaml`

其含义分别是：

- `AGENTS.md`：共享工作规则
- `CLAUDE.md`：Claude Code 补充说明
- `.agentskeleton/context.yaml`：当前文档引导状态

如果 `.agentskeleton/context.yaml` 不存在，宿主应判断：

- 这是新项目初始化场景
或
- 这是老项目文档化重塑场景

然后自动调用：

- `init-docs`
或
- `reshape-docs`

## 推荐主链路

当前推荐的稳定主链路是：

1. `init-docs` / `reshape-docs`
2. `plan`
3. `focus-doc`
4. 宿主模型起草文档
5. `response --apply`

这条链路是当前 MVP 的核心协议面。

## 各命令的宿主职责

### 1. `init-docs` / `reshape-docs`

宿主应在以下场景自动调用：

- 项目还没有 `.agentskeleton/context.yaml`
- 用户明确说“这是新项目”
- 用户明确说“这是老项目，需要文档化重塑”

宿主不应把这个决策交给用户自己去记命令。

### 2. `plan`

宿主应把 `plan` 当成当前项目状态快照。

重点读取：

- `recommended_documents`
- `current_priority`
- `review_candidates`

其中：

- `current_priority` 表示当前最该推进的文档
- `review_candidates` 表示最近变更批次触发的临时回看集合

### 3. `focus-doc`

宿主应在准备起草某份文档前调用 `focus-doc`。

重点读取：

- `change_batch_id`
- `change_batch_inputs`
- `required_context`
- `missing_context`
- `available_context`
- `suggested_sections`
- `review_after_draft`

宿主应把它视为“起草上下文包”。

### 4. 文档起草

宿主模型负责真正生成文档草案。

写作时必须遵守：

- 不编造缺失事实
- 若上下文不足，使用显式占位
- 尽量使用可复用、可公开的语言
- 不把项目引导系统变成业务代码生成器

### 5. `response --apply`

当宿主拿到结构化回答后，应调用 `response --apply` 写回上下文。

重点读取返回中的：

- `post_apply_plan`

宿主应直接使用它继续推进，而不是总是再额外调用一次 `plan`。

## 过期起草包处理

宿主必须处理起草包过期问题。

判断规则：

- 如果当前 `focus-doc.change_batch_id` 已落后于最新上下文批次
- 则这份起草包视为过期
- 必须重新调用 `focus-doc`

不要继续使用旧起草包起草或收敛文档。

## 回看机制

`review_candidates` 和 `review_after_draft` 都不是永久 backlog。

它们的含义是：

- 仅针对最近一次变更批次
- 当前轮里哪些已生成文档需要补充或收敛

宿主不应把它们缓存成长期“已审/未审”状态。

## Codex 建议

对于 Codex：

- 优先读取 `AGENTS.md`
- 如项目启用 AgentSkeleton，优先按主链路调用
- 把 `focus-doc` 视为首选起草输入，而不是直接从 `plan` 猜测文档结构
- 在 `response --apply` 成功后，直接消费 `post_apply_plan`

## Claude Code 建议

对于 Claude Code：

- 优先读取 `CLAUDE.md` 和 `AGENTS.md`
- 只有在确有必要时使用 Claude 专属偏好
- 默认仍遵循共享主链路
- 使用 `focus-doc` 和 `post_apply_plan` 作为主要推进信号

## 当前建议

当前阶段，宿主最应使用的能力只有：

- `init-docs` / `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`

以下能力当前更适合作为辅助或实验能力：

- `workflow`
- `--write-plan-files`
- `--persist-trace`
- `--auto-repair`

## 结论

宿主接入的目标不是让用户学会 CLI，而是让用户几乎感知不到 CLI 的存在。

如果接入是成功的，用户只会感觉：

- agent 更会问问题
- agent 更会组织文档
- agent 更能持续推进和回看收敛

而不是“多了一个要自己操作的新工具”。
