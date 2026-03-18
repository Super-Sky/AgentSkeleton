# CLI 运行手册

本手册描述文档引导流程的最小端到端闭环。

## 新项目流程

1. 初始化上下文：

```bash
agentskeleton init-docs --name MallHub --context .agentskeleton/context.yaml
```

2. 生成计划：

```bash
agentskeleton plan --context .agentskeleton/context.yaml --format yaml
```

3. 生成宿主模型初始提示词：

```bash
agentskeleton prompt --context .agentskeleton/context.yaml --mode initial --format yaml
```

4. 校验并写回宿主模型返回：

```bash
agentskeleton response \
  --file /path/to/host-response.yaml \
  --context .agentskeleton/context.yaml \
  --attempt 0 \
  --apply \
  --docs README.md,docs/domain-overview.md
```

如果返回 `data` 中包含多个字段，默认会批量写回。只有在你希望单字段写回时才需要加 `--question <id>`。

5. 继续下一轮问题：

```bash
agentskeleton next --context .agentskeleton/context.yaml --format yaml
```

## 单命令流程

你可以用一个命令执行打包步骤：

```bash
agentskeleton workflow --context .agentskeleton/context.yaml --format yaml
```

如果已经拿到宿主模型返回：

```bash
agentskeleton workflow \
  --context .agentskeleton/context.yaml \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --apply \
  --question project_summary \
  --docs README.md,docs/domain-overview.md \
  --format yaml
```

## 重试循环

如果 `agentskeleton response` 返回 `decision: retry`：

1. 生成修复提示词：

```bash
agentskeleton prompt \
  --context .agentskeleton/context.yaml \
  --mode repair \
  --errors "missing required field: project_summary" \
  --format yaml
```

2. 让宿主模型只修复结构，不重写全部内容。
3. 使用递增的 `--attempt` 再次校验。
4. 如果变成 `unresolved`，停止自动写回并转人工处理。
