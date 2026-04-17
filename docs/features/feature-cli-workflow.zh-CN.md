# 功能：CLI 工作流

## 目的

本功能域覆盖当前 CLI 工作流职责。

## 当前范围

CLI 当前支持：

- 初始化新项目文档上下文
- 重塑老项目文档上下文
- 输出工作流计划
- 输出 drafting package
- 校验并写回结构化答案
- 基于现有仓库结构安全刷新上下文
- 输出版本元信息

## 当前稳定命令面

当前稳定命令面为：

- `init-docs`
- `reshape-docs`
- `plan`
- `focus-doc`
- `response --apply`
- `version`

当前辅助命令面还包括：

- `update`

## 相关文档

- `cli-contract.zh-CN.md`
- `cli-runbook.zh-CN.md`
- `protocol-stability.zh-CN.md`
