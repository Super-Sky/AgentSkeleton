# Shopping Mall Documentation Blueprint

This template is a sanitized documentation blueprint for shopping mall style platforms.

It intentionally mirrors only reusable documentation patterns from a larger production project:

- project-level overview documents
- agent-facing repository documents
- architecture and domain explanation documents
- roadmap and constraints documents
- repository document map

It does not include:

- real business entities from the source project
- production secrets or config values
- customer data models
- internal route names
- proprietary package names

## Intended Use

Use this template when you want a domain-specific documentation pack for projects such as:

- shopping mall management
- store operations
- merchant collaboration
- product catalog operations
- order-supporting platform work

## Template Shape

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

## Modeling Notes

This template keeps only the documentation patterns worth reusing:

- a clear repository-level summary
- explicit agent collaboration rules
- domain overview for new readers and models
- document map to explain repository documentation intent
- architecture, roadmap, and constraints documents
- workflow, delivery, and execution guidance documents
- versioned feature-document structure
- legacy-project documentation reshaping flow

## Structure Policy

- This blueprint does not fix a single code architecture.
- For new projects, it may describe a recommended structure such as `internal/app`.
- For existing projects, it should document the structure that is already present.

The mall domain is used only as a neutral example domain.
