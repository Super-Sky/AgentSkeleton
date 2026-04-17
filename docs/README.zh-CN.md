# 文档索引

本目录维护 AgentSkeleton 的产品文档、工作流契约、宿主接入规则与发布路径资产。

仓库现在采用四层文档模型：

- 当前真相
- 当前功能文档
- 长期项目正文
- 版本快照

AI 与人默认不应无序阅读这些层级。

## 默认阅读顺序

除非任务明确需要历史版本上下文，否则按下面顺序读取：

1. `current-capabilities.zh-CN.md`
2. `features/README.zh-CN.md`
3. 长期项目正文
4. 只有在需要时再读版本快照

## 分层定义

### 1. 当前真相

使用：

- `current-capabilities.zh-CN.md`

目的：

- 描述仓库当前具体支持什么
- 汇总当前发布路径状态
- 避免读者从旧版本快照里反推当前状态

### 2. 当前功能文档

使用：

- `features/README.zh-CN.md`

目的：

- 描述当前主要能力域
- 作为当前仓库真相的 feature 层入口
- 避免把当前行为分散到多个版本目录

### 3. 长期项目正文

使用：

- `principles.zh-CN.md`
- `architecture.zh-CN.md`
- `cli-contract.zh-CN.md`
- `host-adapter-boundary.zh-CN.md`
- `host-action-mapping.zh-CN.md`
- `zero-command-flow.zh-CN.md`
- `protocol-stability.zh-CN.md`

目的：

- 定义长期有效的产品定位、工作流语义和宿主契约
- 不依赖单一版本存在

### 4. 版本快照

使用：

- `v0.1.0-gap-analysis.zh-CN.md`
- `v0.1.0-implementation-plan.zh-CN.md`

目的：

- 承载某个发布路径的计划与范围
- 保留某一轮版本快照
- 不替代当前真相文档

## 当前主要入口文件

- `current-capabilities.zh-CN.md`
- `features/README.zh-CN.md`
- `principles.zh-CN.md`
- `architecture.zh-CN.md`
- `cli-contract.zh-CN.md`
- `host-integration.zh-CN.md`
- `host-action-mapping.zh-CN.md`
- `zero-command-flow.zh-CN.md`
- `protocol-stability.zh-CN.md`

## 验证与发布资产

- `host-validation-scenarios.zh-CN.md`
- `host-validation-report-template.zh-CN.md`
- `validation-reports/README.zh-CN.md`
- `known-limitations.zh-CN.md`

## 技能与宿主接入文档

- `codex-integration.zh-CN.md`
- `claude-integration.zh-CN.md`
- `skill-sync.zh-CN.md`

## 规则

如果某份文档描述的是当前产品真相，它应放在：

- `current-capabilities*.md`
或
- `features/*.md`

如果某份文档只是解释某一轮发布计划，则不应成为当前真相的主入口。
