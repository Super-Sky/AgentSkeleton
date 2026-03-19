package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const defaultBlueprintPath = "templates/shopping-mall-doc-blueprint/template.yaml"

var jinjaVariablePattern = regexp.MustCompile(`{{\s*([a-zA-Z0-9_]+)\s*}}`)

type blueprintSpec struct {
	Outputs []blueprintOutput `yaml:"outputs"`
}

type blueprintOutput struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

type writeResult struct {
	Created []string `yaml:"created" json:"created"`
	Skipped []string `yaml:"skipped" json:"skipped"`
}

func renderPlannedFiles(ctx Context, plan PlanOutput, overwrite bool) (writeResult, error) {
	specPath, err := resolveBlueprintPath()
	if err != nil {
		return writeResult{}, err
	}

	spec, err := loadBlueprintSpec(specPath)
	if err != nil {
		return writeResult{}, err
	}

	planned := plannedRelativePaths(ctx, plan)
	outputMap := make(map[string]blueprintOutput, len(spec.Outputs))
	for _, out := range spec.Outputs {
		outputMap[out.Path] = out
	}

	data := blueprintData(ctx)
	result := writeResult{
		Created: []string{},
		Skipped: []string{},
	}

	for _, rel := range planned {
		out, ok := outputMap[rel]
		if !ok {
			continue
		}
		targetPath := filepath.Join(ctx.Paths.OutputDir, rel)
		if !overwrite {
			if _, err := os.Stat(targetPath); err == nil {
				result.Skipped = append(result.Skipped, targetPath)
				continue
			}
		}
		if err := renderTemplateFile(specPath, out.Template, targetPath, data); err != nil {
			return result, err
		}
		result.Created = append(result.Created, targetPath)
	}

	return result, nil
}

func resolveBlueprintPath() (string, error) {
	if _, err := os.Stat(defaultBlueprintPath); err == nil {
		return defaultBlueprintPath, nil
	}

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve blueprint path: unable to determine source location")
	}

	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	specPath := filepath.Join(repoRoot, defaultBlueprintPath)
	if _, err := os.Stat(specPath); err != nil {
		return "", fmt.Errorf("resolve blueprint path: %w", err)
	}

	return specPath, nil
}

func loadBlueprintSpec(specPath string) (blueprintSpec, error) {
	var spec blueprintSpec
	data, err := os.ReadFile(specPath)
	if err != nil {
		return spec, fmt.Errorf("read blueprint spec: %w", err)
	}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return spec, fmt.Errorf("parse blueprint spec: %w", err)
	}
	return spec, nil
}

func renderTemplateFile(specPath, templateRelPath, targetPath string, data map[string]string) error {
	specDir := filepath.Dir(specPath)
	templatePath := filepath.Join(specDir, templateRelPath)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read template file: %w", err)
	}

	normalized := normalizeTemplateSyntax(string(content))
	tpl, err := template.New(filepath.Base(templatePath)).Option("missingkey=zero").Parse(normalized)
	if err != nil {
		return fmt.Errorf("parse template file: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()

	if err := tpl.Execute(file, data); err != nil {
		return fmt.Errorf("render template file: %w", err)
	}
	return nil
}

func normalizeTemplateSyntax(content string) string {
	return jinjaVariablePattern.ReplaceAllString(content, "{{.$1}}")
}

func plannedRelativePaths(ctx Context, plan PlanOutput) []string {
	rels := make([]string, 0, len(plan.RecommendedDocuments))
	seen := map[string]struct{}{}
	for _, doc := range plan.RecommendedDocuments {
		rel := strings.TrimPrefix(doc.Path, ctx.Paths.OutputDir+string(filepath.Separator))
		rel = strings.TrimPrefix(rel, ctx.Paths.OutputDir+"/")
		if filepath.IsAbs(doc.Path) && strings.HasPrefix(doc.Path, ctx.Paths.OutputDir) {
			rel = strings.TrimPrefix(doc.Path, ctx.Paths.OutputDir)
			rel = strings.TrimPrefix(rel, string(filepath.Separator))
		} else if filepath.IsAbs(doc.Path) {
			continue
		}
		if rel == "" {
			continue
		}
		if _, ok := seen[rel]; ok {
			continue
		}
		seen[rel] = struct{}{}
		rels = append(rels, rel)
	}
	return rels
}

func blueprintData(ctx Context) map[string]string {
	return map[string]string{
		"project_name":           fallback(ctx.Project.Name, "Unnamed Project"),
		"project_summary":        fallback(ctx.Project.Summary, answeredValue(ctx, "project_summary"), "Project summary to be clarified"),
		"primary_users":          joinedOrFallback(ctx.Project.PrimaryUsers, "to be clarified"),
		"deployment_shape":       fallback(answeredValue(ctx, "deployment_shape"), "to be clarified"),
		"documentation_priority": "repository clarity and agent collaboration",
		"release_version":        fallback(ctx.Documentation.ReleaseVersion, "v0.0.0"),
	}
}

func answeredValue(ctx Context, id string) string {
	for _, a := range ctx.Conversation.AnsweredQuestions {
		if a.ID == id {
			return a.Value
		}
	}
	return ""
}

func fallback(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func joinedOrFallback(values []string, fallbackValue string) string {
	if len(values) == 0 {
		return fallbackValue
	}
	return strings.Join(values, ", ")
}
