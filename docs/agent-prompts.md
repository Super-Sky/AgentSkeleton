# Agent Prompt Guide

These prompts show how host models should consume `plan` / `next` outputs and respond using the structured envelope.

## Consumption Pattern

1. Run `agentskeleton plan` → capture the YAML/JSON output (context, recommended documents).
2. Build a prompt around the `recommended_documents` list and `missing_information` array, asking for factual answers to the `open_questions` from the CLI output.
3. Request that the host respond using the response envelope schema:

```yaml
status: <ok|invalid|unresolved>
schema: question-answer-set-v1
data:
  <field>: "<value>"
errors: []
raw_text: ""
```

4. If you receive `invalid`, map the `errors` back into a stricter prompt that focuses only on the missing fields.
5. If you reach `unresolved`, pause automation and flag the context for manual review.

## Sample Prompt

```
Based on this plan, please supply the missing answers in the requested schema.

Plan -> missing_information:
- project_summary
- deployment_shape
- ownership_model

Respond with the envelope and make sure `status` is `ok`.
```

## Chinese Prompts

For Claude Code, you can mirror the above instructions in Chinese, but the schema must stay the same.
