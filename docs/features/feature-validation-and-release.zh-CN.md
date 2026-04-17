# 功能：验证与发布

## 目的

本功能域覆盖当前验证与发布加固资产。

## 当前范围

仓库当前提供：

- 宿主验证场景
- 验证报告模板
- 种子验证报告
- 本地 smoke 测试
- CI 与 release build workflow
- 已知限制与发布路径 changelog 更新

## 当前产物

- `host-validation-scenarios.zh-CN.md`
- `host-validation-report-template.zh-CN.md`
- `validation-reports/README.zh-CN.md`
- `known-limitations.zh-CN.md`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `Makefile`
- `scripts/smoke_test.sh`

## 当前限制

发布门槛仍依赖真实 Codex 与 Claude Code 验证证据，以及一次真实 tag release 运行。
