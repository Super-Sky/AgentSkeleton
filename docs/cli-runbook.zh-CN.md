# CLI 运行手册

本手册描述文档引导流程的最小端到端闭环。

## 新项目流程

1. 初始化上下文：

```bash
agentskeleton init-docs --project /path/to/project --output-dir /path/to/output --name MallHub
```

这会产生两类文件：

```text
/path/to/output/.agentskeleton/    # AgentSkeleton 过程产物
/path/to/output/README.md          # 最终文档与技能放在这里
/path/to/output/docs/...
```

2. 生成计划：

```bash
agentskeleton plan --project /path/to/project --output-dir /path/to/output --format yaml
```

如果不传 `--context`，默认解析为：

```text
/path/to/output/.agentskeleton/context.yaml
```

3. 生成宿主模型初始提示词：

```bash
agentskeleton prompt --project /path/to/project --output-dir /path/to/output --mode initial --format yaml
```

4. 校验并写回宿主模型返回：

```bash
agentskeleton response \
  --file /path/to/host-response.yaml \
  --project /path/to/project \
  --output-dir /path/to/output \
  --attempt 0 \
  --apply \
  --docs /path/to/output/README.md,/path/to/output/docs/domain-overview.md
```

如果返回 `data` 中包含多个字段，默认会批量写回。只有在你希望单字段写回时才需要加 `--question <id>`。

5. 继续下一轮问题：

```bash
agentskeleton next --project /path/to/project --output-dir /path/to/output --format yaml
```

## 单命令流程

你可以用一个命令执行打包步骤：

```bash
agentskeleton workflow --project /path/to/project --output-dir /path/to/output --format yaml
```

如果已经拿到宿主模型返回：

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --apply \
  --docs /path/to/output/README.md,/path/to/output/docs/domain-overview.md \
  --format yaml
```

## 清理模型

- `<output-dir>/.agentskeleton` 只保存 AgentSkeleton 的过程产物。
- `<output-dir>/...` 下的最终文档和技能属于用户可保留的交付结果。
- 如果用户删除 `<output-dir>/.agentskeleton`，最终文档不受影响，只是 AgentSkeleton 不再保留这次过程状态。

## 重试循环

如果 `agentskeleton response` 返回 `decision: retry`：

1. 生成修复提示词：

```bash
agentskeleton prompt \
  --project /path/to/project \
  --output-dir /path/to/output \
  --mode repair \
  --errors "missing required field: project_summary" \
  --format yaml
```

2. 让宿主模型只修复结构，不重写全部内容。
3. 使用递增的 `--attempt` 再次校验。
4. 如果变成 `unresolved`，停止自动写回并转人工处理。
