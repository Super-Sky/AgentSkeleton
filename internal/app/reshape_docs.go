package app

import (
	"flag"
	"fmt"
)

func runReshapeDocs(args []string) error {
	fs := flag.NewFlagSet("reshape-docs", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	host := fs.String("host", "claude-code", "host environment")
	name := fs.String("name", "", "project name")
	summary := fs.String("summary", "", "project summary")
	domain := fs.String("domain", "", "project domain")
	layout := fs.String("layout-summary", "", "current layout summary")
	format := fs.String("format", "yaml", "output format")
	if err := fs.Parse(args); err != nil {
		return err
	}

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name:         *name,
			Summary:      *summary,
			Mode:         "legacy",
			Domain:       *domain,
			PrimaryUsers: []string{},
			Host:         *host,
		},
		Documentation: Documentation{
			Phase:          "planning",
			GeneratedDocs:  []string{},
			MissingDocs:    defaultLegacyMissingDocs(),
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy:             "existing",
			RecommendedLayout:    "internal/app",
			CurrentLayoutSummary: *layout,
		},
		Conversation: Conversation{
			AnsweredQuestions: []QuestionAnswer{
				{ID: "project_mode", Value: "legacy"},
			},
			OpenQuestions: []string{
				"undocumented_directories",
				"active_release_docs_strategy",
				"ownership_model",
			},
		},
	}
	if *name != "" {
		ctx.Conversation.AnsweredQuestions = append(ctx.Conversation.AnsweredQuestions, QuestionAnswer{
			ID:    "project_name",
			Value: *name,
		})
	}
	if *layout != "" {
		ctx.Conversation.AnsweredQuestions = append(ctx.Conversation.AnsweredQuestions, QuestionAnswer{
			ID:    "current_layout_summary",
			Value: *layout,
		})
	}

	if err := writeContext(*contextPath, ctx); err != nil {
		return err
	}

	return printOutput(*format, map[string]any{
		"command":      "reshape-docs",
		"context_path": *contextPath,
		"project_mode": "legacy",
		"status":       "initialized",
		"next_hint":    fmt.Sprintf("run `agentskeleton plan --context %s`", *contextPath),
	})
}

func defaultLegacyMissingDocs() []string {
	return []string{
		"docs/domain-overview.md",
		"docs/architecture.md",
		"docs/legacy-reshape-guide.md",
		"docs/legacy-structure-inventory.md",
	}
}
