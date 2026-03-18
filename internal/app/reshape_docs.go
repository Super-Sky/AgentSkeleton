package app

import (
	"flag"
	"fmt"
)

func runReshapeDocs(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("reshape-docs", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	outputDir := fs.String("output-dir", "", "documentation output directory (default: project root)")
	host := fs.String("host", "claude-code", "host environment")
	name := fs.String("name", "", "project name")
	summary := fs.String("summary", "", "project summary")
	domain := fs.String("domain", "", "project domain")
	layout := fs.String("layout-summary", "", "current layout summary")
	format := fs.String("format", "yaml", "output format")
	if err := fs.Parse(args); err != nil {
		return err
	}
	projectRoot, err := resolveProjectRoot(*project)
	if err != nil {
		return err
	}
	outputRoot, err := resolveOutputDir(projectRoot, *outputDir)
	if err != nil {
		return err
	}
	resolvedContext := resolveContextPath(projectRoot, *contextPath, contextSet)

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectRoot,
			OutputDir:   outputRoot,
			ContextPath: resolvedContext,
		},
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
			MissingDocs:    defaultLegacyMissingDocs(outputRoot),
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

	if err := writeContext(resolvedContext, ctx); err != nil {
		return err
	}

	return printOutput(*format, map[string]any{
		"command":      "reshape-docs",
		"context_path": resolvedContext,
		"project_root": projectRoot,
		"output_dir":   outputRoot,
		"project_mode": "legacy",
		"status":       "initialized",
		"next_hint":    fmt.Sprintf("run `agentskeleton plan --context %s`", resolvedContext),
	})
}

func defaultLegacyMissingDocs(outputRoot string) []string {
	return []string{
		resolveDocPath(outputRoot, "docs/domain-overview.md"),
		resolveDocPath(outputRoot, "docs/architecture.md"),
		resolveDocPath(outputRoot, "docs/legacy-reshape-guide.md"),
		resolveDocPath(outputRoot, "docs/legacy-structure-inventory.md"),
	}
}
