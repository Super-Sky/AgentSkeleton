# Shopping Mall Documentation Blueprint

这是一个经过脱敏处理的购物商场类文档蓝图模板。

它只保留了一个较大生产项目中可复用的文档模式，不包含任何真实业务情报：

- 项目级总览文档
- 面向 agent 的仓库说明文件
- 架构与领域说明文档
- 路线图与约束文档
- 仓库文档地图

它不会包含：

- 来源项目中的真实业务实体
- 生产密钥或真实配置
- 客户数据模型
- 内部接口命名
- 私有包名

## 适用场景

当你需要为以下类型项目生成一组领域化文档包时，可以使用这个模板：

- 商场管理平台
- 门店运营系统
- 商户协作平台
- 商品目录运营
- 订单配套平台工作

## 模板结构

```text
shopping-mall-doc-blueprint/
├── template.yaml
└── files/
    ├── AGENTS.md.tmpl
    ├── CLAUDE.md.tmpl
    ├── README.md.tmpl
    ├── docs/architecture.md.tmpl
    ├── docs/code-style-guide.md.tmpl
    ├── docs/constraints.md.tmpl
    ├── docs/doc-templates.md.tmpl
    ├── docs/document-map.md.tmpl
    ├── docs/domain-overview.md.tmpl
    ├── docs/repo-workflow-guide.md.tmpl
    ├── docs/roadmap.md.tmpl
    ├── docs/task-delivery-guide.md.tmpl
    ├── docs/task-execution-template.md.tmpl
    ├── docs/legacy-reshape-guide.md.tmpl
    ├── docs/legacy-structure-inventory.md.tmpl
    └── docs/versioned/
```

## 建模说明

这个模板只保留值得复用的文档模式：

- 清晰的仓库级项目摘要
- 显式的 agent 协作规则
- 面向新读者和模型的领域说明
- 用于解释仓库文档意图的文档地图
- 架构、路线图和约束文档
- 工作流、交付规范和执行模板文档
- 版本化 feature 文档结构
- 老项目文档化重塑流程

## 结构策略

- 这个蓝图不会写死单一代码架构。
- 对新项目，它可以描述推荐结构，例如 `internal/app`。
- 对已有项目，它应优先文档化当前已经存在的结构。

商场领域仅用于提供一个中性的示例业务上下文。
