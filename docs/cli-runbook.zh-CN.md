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

6. 让 CLI 明确告诉宿主当前该起草哪份文档，以及可用上下文：

```bash
agentskeleton focus-doc --project /path/to/project --output-dir /path/to/output --format yaml
```

如果你想强制聚焦某一份文档，而不是当前优先文档：

```bash
agentskeleton focus-doc \
  --project /path/to/project \
  --output-dir /path/to/output \
  --path docs/architecture.md \
  --format yaml
```

## 单命令流程

你可以用一个命令执行打包步骤：

```bash
agentskeleton workflow --project /path/to/project --output-dir /path/to/output --format yaml
```

如果你希望把当前已支持的计划文档直接落到输出目录：

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --write-plan-files \
  --format yaml
```

默认只写缺失文件；只有你明确想覆盖已有生成文档时，才传 `--overwrite`。
当文件被创建，或检测到已经存在时，workflow 也会把它们的 generated 状态回写到 `<output-dir>/.agentskeleton/context.yaml`。

如果你希望为这一轮保留一份完整结构化快照，便于审计或后续回放：

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --persist-trace \
  --format yaml
```

这会把 trace 文件写到 `<output-dir>/.agentskeleton/traces/`，文件名会自动带上当前文档阶段，并在 CLI 输出中返回 `trace_path`。

如果已经拿到宿主模型返回：

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --apply \
  --write-plan-files \
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

## 回溯收敛

当新的答案会影响已经存在的文档时，`plan` 和 `workflow` 会输出 `review_candidates`。

你应该用它来回看和收敛这些已生成文档，尤其是在以下情况之后：

- `project_summary`、`deployment_shape`、`ownership_model` 这类核心答案被补齐
- 老项目盘点发现新的目录或模块信息
- 某次版本文档策略被明确

你也可以直接在 `workflow` 中打包这个过程：

```bash
agentskeleton workflow \
  --project /path/to/project \
  --output-dir /path/to/output \
  --response-file /path/to/host-response.yaml \
  --attempt 0 \
  --auto-repair \
  --format yaml
```

当返回结果可重试时，输出里会包含 `auto_repair` 块，其中带有下一次尝试编号、校验错误、修复提示词和宿主模型循环说明。
