package app

import (
	"flag"
	"fmt"
)

type PlanOutput struct {
	Command              string           `yaml:"command" json:"command"`
	ProjectMode          string           `yaml:"project_mode" json:"project_mode"`
	DocumentationPhase   string           `yaml:"documentation_phase" json:"documentation_phase"`
	ReleaseVersion       string           `yaml:"release_version" json:"release_version"`
	KnownFacts           []Fact           `yaml:"known_facts" json:"known_facts"`
	MissingInformation   []string         `yaml:"missing_information" json:"missing_information"`
	RecommendedDocuments []DocumentAdvice `yaml:"recommended_documents" json:"recommended_documents"`
	NextActions          []string         `yaml:"next_actions" json:"next_actions"`
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

func runPlan(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("plan", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	format := fs.String("format", "yaml", "output format")
	if err := fs.Parse(args); err != nil {
		return err
	}
	projectRoot, err := resolveProjectRoot(*project)
	if err != nil {
		return err
	}
	resolvedContext := resolveContextPath(projectRoot, *contextPath, contextSet)

	ctx, err := loadContext(resolvedContext)
	if err != nil {
		return err
	}
	if err := validateMode(ctx.Project.Mode); err != nil {
		return err
	}

	out := PlanOutput{
		Command:            "plan",
		ProjectMode:        ctx.Project.Mode,
		DocumentationPhase: ctx.Documentation.Phase,
		ReleaseVersion:     ctx.Documentation.ReleaseVersion,
		KnownFacts:         buildKnownFacts(ctx),
		MissingInformation: append([]string{}, ctx.Conversation.OpenQuestions...),
		RecommendedDocuments: append(
			recommendedDocumentsForMode(ctx.Project.Mode),
			versionedDocuments(ctx.Documentation.ReleaseVersion)...),
		NextActions: nextActionsForMode(ctx.Project.Mode),
	}

	return printOutput(*format, out)
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
