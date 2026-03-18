# Agent 提示指南

这份文档说明宿主模型如何消费 `plan` / `next` 输出，并按结构化包络返回答案。

## 使用流程

1. 先运行 `agentskeleton plan`，拿到 YAML/JSON 输出，里面有 `recommended_documents` 和 `missing_information`。
2. 构建 prompt，围绕推荐文档和缺失信息提问，明确需要补哪些内容。
3. 要求宿主模型按照以下返回包络答复：

```yaml
status: <0k|invalid|unresolved>
schema: question-answer-set-v1
data:
  <字段>: "<值>"
errors: []
raw_text: ""
```

4. 如果得到 `invalid`，根据 `errors` 里的提示进一步收紧问题，要求只回答缺失的字段。
5. 收到 `unresolved` 即停止自动流，转人工干预。

## 示例提示

```
根据这个 plan，按给定的 schema 补全缺失的字段。

Plan -> missing_information:
- project_summary
- deployment_shape
- ownership_model

返回时确保 status 为 ok。
```

## 同步提示

Claude Code 的提示可以用中文描述，但是返回的结构必须保持不变。
