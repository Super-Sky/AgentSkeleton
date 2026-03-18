package app

import (
	"flag"
	"fmt"
)

func runInitDocs(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("init-docs", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	outputDir := fs.String("output-dir", "", "documentation output directory (default: project root)")
	host := fs.String("host", "codex", "host environment")
	name := fs.String("name", "", "project name")
	summary := fs.String("summary", "", "project summary")
	domain := fs.String("domain", "", "project domain")
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
			Mode:         "new",
			Domain:       *domain,
			PrimaryUsers: []string{},
			Host:         *host,
		},
		Documentation: Documentation{
			Phase:          "discovery",
			GeneratedDocs:  []string{},
			MissingDocs:    defaultNewMissingDocs(outputRoot),
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

	if err := writeContext(resolvedContext, ctx); err != nil {
		return err
	}

	return printOutput(*format, map[string]any{
		"command":      "init-docs",
		"context_path": resolvedContext,
		"project_root": projectRoot,
		"output_dir":   outputRoot,
		"project_mode": "new",
		"status":       "initialized",
		"next_hint":    fmt.Sprintf("run `agentskeleton plan --context %s`", resolvedContext),
	})
}

func defaultNewMissingDocs(outputRoot string) []string {
	return []string{
		resolveDocPath(outputRoot, "README.md"),
		resolveDocPath(outputRoot, "AGENTS.md"),
		resolveDocPath(outputRoot, "CLAUDE.md"),
		resolveDocPath(outputRoot, "docs/domain-overview.md"),
		resolveDocPath(outputRoot, "docs/architecture.md"),
	}
}
