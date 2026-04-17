# 技能同步规则

## 目的

本文档用于定义共享正文与宿主 skill 之间的 source-of-truth 关系。

## 当前模型

仓库当前采用：

- `docs/*` 中的共享正文
- `skills/*` 中的宿主入口 skill

skill 应保持足够轻。

## Source Of Truth 规则

- 共享工作流规则统一维护在 `docs/*`
- skill 应指向共享正文，而不是重复长段规则正文
- 当前宿主行为真相应先更新文档，再更新 skill
- 文档变更后，应检查两个宿主 skill 是否需要同步

## 当前技能

- `skills/agentskeleton-codex/SKILL.md`
- `skills/agentskeleton-claude/SKILL.md`

## 同步要求

当以下内容发生变化时：

- 宿主路由规则
- 零命令主链路
- 稳定协议边界
- 升级/澄清策略

维护者应同步检查两个宿主 skill。

## 反模式

不要让 Codex 与 Claude skill 在没有显式说明的情况下，演化出两套不同的工作流语义。
