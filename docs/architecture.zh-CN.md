# Architecture Notes

## 核心模型

AgentSkeleton 有两个核心层：

- 存储在仓库中的文档蓝图资产
- 负责引导、规划和校验文档工作的 CLI

这个仓库预期通过 AI 辅助的文档工作持续演进，人类更像产品负责人、评审者和约束定义者。

## 支持模型

项目应通过以下方式同时支持 Codex 与 Claude Code：

- 一套共享结构
- 一套共享文档模型
- 仅在需要时补充少量 agent 专属文件

## 交付模型

- 用户在 Codex 或 Claude Code 的对话环境中工作。
- AgentSkeleton 以 CLI 形式运行在这些环境中。
- CLI 输出结构化提示、计划和文档期望。
- 宿主模型将这些输出转化为实际仓库文档。

## 仓库结构处理策略

- 对新项目，AgentSkeleton 可以推荐默认结构，以便提升文档化和协作效率。
- 对于应用型新项目，当前优先推荐 `internal/app` 导向的布局。
- 对已有项目，AgentSkeleton 应优先文档化现有结构，而不是默认强行替换架构。

## 初始目录意图

- `cmd/`：CLI 入口
- `internal/`：内部实现包
- `templates/`：文档蓝图与引导资产
- `docs/`：产品与设计文档
