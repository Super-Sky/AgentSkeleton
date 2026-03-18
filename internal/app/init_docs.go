package app

import (
	"flag"
	"fmt"
)

func runInitDocs(args []string) error {
	fs := flag.NewFlagSet("init-docs", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	host := fs.String("host", "codex", "host environment")
	name := fs.String("name", "", "project name")
	summary := fs.String("summary", "", "project summary")
	domain := fs.String("domain", "", "project domain")
	format := fs.String("format", "yaml", "output format")
	if err := fs.Parse(args); err != nil {
		return err
	}

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name:         *name,
			Summary:      *summary,
			Mode:         "new",
			Domain:       *domain,
			PrimaryUsers: []string{},
			Host:         *host,
		},
		Documentation: Documentation{
			Phase:          "discovery",
			GeneratedDocs:  []string{},
			MissingDocs:    defaultNewMissingDocs(),
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy:             "recommended",
			RecommendedLayout:    "internal/app",
			CurrentLayoutSummary: "",
		},
		Conversation: Conversation{
			AnsweredQuestions: []QuestionAnswer{
				{ID: "project_mode", Value: "new"},
			},
			OpenQuestions: []string{
				"project_summary",
				"deployment_shape",
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

	if err := writeContext(*contextPath, ctx); err != nil {
		return err
	}

	return printOutput(*format, map[string]any{
		"command":      "init-docs",
		"context_path": *contextPath,
		"project_mode": "new",
		"status":       "initialized",
		"next_hint":    fmt.Sprintf("run `agentskeleton plan --context %s`", *contextPath),
	})
}

func defaultNewMissingDocs() []string {
	return []string{
		"README.md",
		"AGENTS.md",
		"CLAUDE.md",
		"docs/domain-overview.md",
		"docs/architecture.md",
	}
}
