package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type PlanOutput struct {
	Command              string            `yaml:"command" json:"command"`
	ProjectMode          string            `yaml:"project_mode" json:"project_mode"`
	DocumentationPhase   string            `yaml:"documentation_phase" json:"documentation_phase"`
	ReleaseVersion       string            `yaml:"release_version" json:"release_version"`
	KnownFacts           []Fact            `yaml:"known_facts" json:"known_facts"`
	MissingInformation   []string          `yaml:"missing_information" json:"missing_information"`
	RecommendedDocuments []DocumentAdvice  `yaml:"recommended_documents" json:"recommended_documents"`
	CurrentPriority      *PriorityDocument `yaml:"current_priority,omitempty" json:"current_priority,omitempty"`
	NextActions          []string          `yaml:"next_actions" json:"next_actions"`
}

type Fact struct {
	Key   string `yaml:"key" json:"key"`
	Value any    `yaml:"value" json:"value"`
}

type DocumentAdvice struct {
	Path    string `yaml:"path" json:"path"`
	Purpose string `yaml:"purpose" json:"purpose"`
	Status  string `yaml:"status" json:"status"`
}

type PriorityDocument struct {
	Path            string   `yaml:"path" json:"path"`
	Purpose         string   `yaml:"purpose" json:"purpose"`
	RequiredContext []string `yaml:"required_context" json:"required_context"`
	MissingContext  []string `yaml:"missing_context" json:"missing_context"`
	Ready           bool     `yaml:"ready" json:"ready"`
	Reason          string   `yaml:"reason" json:"reason"`
}

func runPlan(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("plan", flag.ContinueOnError)
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

	out := buildPlanOutput(ctx)

	return printOutput(*format, out)
}

func buildPlanOutput(ctx Context) PlanOutput {
	recommended := append(
		recommendedDocumentsForMode(ctx.Project.Mode),
		versionedDocuments(ctx.Documentation.ReleaseVersion)...)

	return PlanOutput{
		Command:              "plan",
		ProjectMode:          ctx.Project.Mode,
		DocumentationPhase:   ctx.Documentation.Phase,
		ReleaseVersion:       ctx.Documentation.ReleaseVersion,
		KnownFacts:           buildKnownFacts(ctx),
		MissingInformation:   append([]string{}, ctx.Conversation.OpenQuestions...),
		RecommendedDocuments: recommended,
		CurrentPriority:      selectPriorityDocument(ctx, recommended),
		NextActions:          nextActionsForMode(ctx.Project.Mode),
	}
}

func buildKnownFacts(ctx Context) []Fact {
	facts := []Fact{
		{Key: "project_name", Value: ctx.Project.Name},
		{Key: "project_mode", Value: ctx.Project.Mode},
		{Key: "structure_strategy", Value: ctx.Structure.Strategy},
	}
	if len(ctx.Project.PrimaryUsers) > 0 {
		facts = append(facts, Fact{Key: "primary_users", Value: ctx.Project.PrimaryUsers})
	}
	if ctx.Structure.CurrentLayoutSummary != "" {
		facts = append(facts, Fact{Key: "current_layout_summary", Value: ctx.Structure.CurrentLayoutSummary})
	}
	return facts
}

func recommendedDocumentsForMode(mode string) []DocumentAdvice {
	if mode == "legacy" {
		return []DocumentAdvice{
			{Path: "docs/legacy-structure-inventory.md", Purpose: "capture the current repository layout before reshaping docs", Status: "required"},
			{Path: "docs/domain-overview.md", Purpose: "define domain language for future documentation", Status: "required"},
			{Path: "docs/architecture.md", Purpose: "explain the current repository and system shape", Status: "required"},
			{Path: "docs/legacy-reshape-guide.md", Purpose: "define the reshape flow for this repository", Status: "required"},
		}
	}

	return []DocumentAdvice{
		{Path: "README.md", Purpose: "repository entrypoint and project summary", Status: "required"},
		{Path: "AGENTS.md", Purpose: "shared agent collaboration rules", Status: "required"},
		{Path: "CLAUDE.md", Purpose: "Claude Code-specific adaptation notes", Status: "required"},
		{Path: "docs/domain-overview.md", Purpose: "domain language for humans and models", Status: "required"},
		{Path: "docs/architecture.md", Purpose: "describe the intended repository and system shape", Status: "required"},
	}
}

func nextActionsForMode(mode string) []string {
	if mode == "legacy" {
		return []string{
			"list undocumented top-level directories",
			"decide whether versioned feature docs are needed for active work",
			"clarify documentation ownership",
			"draft docs/legacy-structure-inventory.md",
		}
	}

	return []string{
		"ask for a one-sentence project summary",
		"ask for the deployment shape",
		"ask who owns documentation maintenance",
		"draft README.md",
	}
}

func versionedDocuments(release string) []DocumentAdvice {
	if release == "" {
		return nil
	}
	return []DocumentAdvice{
		{Path: fmt.Sprintf("docs/%s/README.md", release), Purpose: "versioned work index and entrypoint", Status: "optional"},
		{Path: fmt.Sprintf("docs/%s/features/README.md", release), Purpose: "feature documentation index for the release", Status: "optional"},
		{Path: fmt.Sprintf("docs/%s/features/feature-template.md", release), Purpose: "feature note template for the release", Status: "optional"},
		{Path: fmt.Sprintf("docs/%s/features/review-checklist-template.md", release), Purpose: "review checklist template for the release", Status: "optional"},
	}
}

func selectPriorityDocument(ctx Context, docs []DocumentAdvice) *PriorityDocument {
	for _, doc := range docs {
		if documentMaterialized(ctx, doc.Path) {
			continue
		}
		required := requiredContextForDocument(ctx.Project.Mode, doc.Path)
		missing := missingRequiredContext(ctx, required)
		reason := "next recommended document that has not been generated yet"
		ready := len(missing) == 0
		if !ready {
			reason = "waiting for missing context before reliable drafting"
		}
		return &PriorityDocument{
			Path:            doc.Path,
			Purpose:         doc.Purpose,
			RequiredContext: required,
			MissingContext:  missing,
			Ready:           ready,
			Reason:          reason,
		}
	}
	return nil
}

func requiredContextForDocument(mode, path string) []string {
	switch path {
	case "README.md":
		return []string{"project_summary", "deployment_shape"}
	case "AGENTS.md", "CLAUDE.md":
		return []string{"ownership_model"}
	case "docs/domain-overview.md":
		return []string{"project_summary"}
	case "docs/architecture.md":
		if mode == "legacy" {
			return []string{"current_layout_summary"}
		}
		return []string{"deployment_shape"}
	case "docs/legacy-structure-inventory.md":
		return []string{"undocumented_directories", "current_layout_summary"}
	case "docs/legacy-reshape-guide.md":
		return []string{"active_release_docs_strategy", "ownership_model"}
	default:
		if isVersionedReleaseReadme(path) {
			return []string{"active_release_docs_strategy", "ownership_model"}
		}
		if isVersionedFeaturePath(path) {
			return []string{"active_release_docs_strategy"}
		}
		return nil
	}
}

func isVersionedReleaseReadme(path string) bool {
	segments := strings.Split(filepath.ToSlash(path), "/")
	return len(segments) == 3 && segments[0] == "docs" && segments[2] == "README.md"
}

func isVersionedFeaturePath(path string) bool {
	segments := strings.Split(filepath.ToSlash(path), "/")
	return len(segments) >= 4 && segments[0] == "docs" && segments[2] == "features"
}

func missingRequiredContext(ctx Context, required []string) []string {
	if len(required) == 0 {
		return nil
	}
	var missing []string
	for _, id := range required {
		if slices.Contains(ctx.Conversation.OpenQuestions, id) {
			missing = append(missing, id)
			continue
		}
		if answeredValue(ctx, id) == "" && !factPresent(ctx, id) {
			missing = append(missing, id)
		}
	}
	return missing
}

func factPresent(ctx Context, id string) bool {
	switch id {
	case "project_summary":
		return ctx.Project.Summary != ""
	case "current_layout_summary":
		return ctx.Structure.CurrentLayoutSummary != ""
	default:
		return false
	}
}

func documentMaterialized(ctx Context, path string) bool {
	rel, abs := ctx.normalizeDocPath(path)
	if rel != "" && slices.Contains(ctx.Documentation.GeneratedDocs, rel) {
		return true
	}
	if abs != "" && slices.Contains(ctx.Documentation.GeneratedDocs, abs) {
		return true
	}
	if abs != "" {
		if _, err := os.Stat(abs); err == nil {
			return true
		}
	}
	return false
}
