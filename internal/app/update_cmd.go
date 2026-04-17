package app

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type UpdateOutput struct {
	Command        string            `yaml:"command" json:"command"`
	ContextPath    string            `yaml:"context_path" json:"context_path"`
	ContextUpdated bool              `yaml:"context_updated" json:"context_updated"`
	Inferred       map[string]string `yaml:"inferred,omitempty" json:"inferred,omitempty"`
	PostUpdatePlan PlanOutput        `yaml:"post_update_plan" json:"post_update_plan"`
}

func runUpdate(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	outputDir := fs.String("output-dir", "", "documentation output directory (default: project root)")
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
	resolvedContext := resolveContextPath(outputRoot, *contextPath, contextSet)

	ctx, err := loadContext(resolvedContext)
	if err != nil {
		return err
	}
	if err := validateMode(ctx.Project.Mode); err != nil {
		return err
	}

	inferred := inferUpdateAnswers(ctx, projectRoot)
	updated := applyInferredAnswers(&ctx, inferred)
	if updated {
		resolved := make([]string, 0, len(inferred))
		for key := range inferred {
			resolved = append(resolved, key)
		}
		sort.Strings(resolved)
		ctx.recordChangeBatch(resolved, nil)
		if err := writeContext(resolvedContext, ctx); err != nil {
			return err
		}
	}

	return printOutput(*format, UpdateOutput{
		Command:        "update",
		ContextPath:    resolvedContext,
		ContextUpdated: updated,
		Inferred:       inferred,
		PostUpdatePlan: buildPlanOutput(ctx),
	})
}

func inferUpdateAnswers(ctx Context, projectRoot string) map[string]string {
	out := map[string]string{}

	switch ctx.Project.Mode {
	case "legacy":
		dirs, files := scanTopLevelProjectEntries(projectRoot)

		if answeredValue(ctx, "undocumented_directories") == "" {
			if len(dirs) > 0 {
				out["undocumented_directories"] = strings.Join(dirs, ",")
			}
		}
		if answeredValue(ctx, "current_layout_summary") == "" && ctx.Structure.CurrentLayoutSummary == "" {
			if summary := inferLegacyLayoutSummary(projectRoot, dirs, files); summary != "" {
				out["current_layout_summary"] = summary
			}
		}
	}

	return out
}

func applyInferredAnswers(ctx *Context, inferred map[string]string) bool {
	if len(inferred) == 0 {
		return false
	}

	updated := false
	keys := make([]string, 0, len(inferred))
	for key := range inferred {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := strings.TrimSpace(inferred[key])
		if value == "" {
			continue
		}
		if current := answeredValue(*ctx, key); current == value {
			continue
		}
		ctx.applyAnswer(key, value)
		updated = true
	}
	return updated
}

func scanTopLevelProjectEntries(projectRoot string) (dirs []string, files []string) {
	entries, err := os.ReadDir(projectRoot)
	if err != nil {
		return nil, nil
	}

	for _, entry := range entries {
		name := entry.Name()
		if shouldSkipTopLevelEntry(name) {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, name)
			continue
		}
		files = append(files, name)
	}

	sort.Strings(dirs)
	sort.Strings(files)
	return dirs, files
}

func shouldSkipTopLevelEntry(name string) bool {
	switch name {
	case ".DS_Store", ".agentskeleton":
		return true
	default:
		return false
	}
}

func inferLegacyLayoutSummary(projectRoot string, dirs []string, files []string) string {
	repoKind := inferRepositoryKind(files)

	parts := []string{repoKind}
	if containsString(files, "main.go") {
		parts = append(parts, "with main.go at repo root")
	}
	if containsString(files, "go.mod") {
		parts = append(parts, "using Go module layout")
	}

	if len(dirs) > 0 {
		displayDirs := dirs
		if len(displayDirs) > 8 {
			displayDirs = displayDirs[:8]
		}
		parts = append(parts, fmt.Sprintf("and top-level directories: %s", strings.Join(displayDirs, ", ")))
	}

	summary := strings.Join(parts, " ")
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return ""
	}
	if !strings.HasSuffix(summary, ".") {
		summary += "."
	}
	return summary
}

func inferRepositoryKind(files []string) string {
	switch {
	case containsString(files, "go.mod"):
		return "Go repository"
	case containsString(files, "package.json"):
		return "Node.js repository"
	case containsString(files, "pyproject.toml"), containsString(files, "requirements.txt"):
		return "Python repository"
	default:
		return "Repository"
	}
}
