package app

import (
	"flag"
	"fmt"
)

type FocusDocOutput struct {
	Command           string            `yaml:"command" json:"command"`
	Path              string            `yaml:"path" json:"path"`
	Purpose           string            `yaml:"purpose" json:"purpose"`
	Ready             bool              `yaml:"ready" json:"ready"`
	Reason            string            `yaml:"reason" json:"reason"`
	RequiredContext   []string          `yaml:"required_context" json:"required_context"`
	MissingContext    []string          `yaml:"missing_context" json:"missing_context"`
	AvailableContext  map[string]string `yaml:"available_context" json:"available_context"`
	SuggestedSections []string          `yaml:"suggested_sections" json:"suggested_sections"`
	WritingRules      []string          `yaml:"writing_rules" json:"writing_rules"`
	NextAfterDraft    []string          `yaml:"next_after_draft" json:"next_after_draft"`
}

func runFocusDoc(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("focus-doc", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	outputDir := fs.String("output-dir", "", "documentation output directory (default: project root)")
	docPath := fs.String("path", "", "explicit document path to focus")
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

	plan := buildPlanOutput(ctx)
	out, err := buildFocusDocOutput(ctx, plan, *docPath)
	if err != nil {
		return err
	}

	return printOutput(*format, out)
}

func buildFocusDocOutput(ctx Context, plan PlanOutput, explicitPath string) (FocusDocOutput, error) {
	doc, err := selectFocusDocument(plan, explicitPath)
	if err != nil {
		return FocusDocOutput{}, err
	}

	required := requiredContextForDocument(ctx.Project.Mode, doc.Path)
	missing := missingRequiredContext(ctx, required)

	return FocusDocOutput{
		Command:           "focus-doc",
		Path:              doc.Path,
		Purpose:           doc.Purpose,
		Ready:             len(missing) == 0,
		Reason:            focusReason(plan, doc.Path, missing),
		RequiredContext:   required,
		MissingContext:    missing,
		AvailableContext:  availableContextForDocument(ctx, required),
		SuggestedSections: suggestedSectionsForDocument(ctx.Project.Mode, doc.Path),
		WritingRules:      writingRulesForDocument(),
		NextAfterDraft:    nextAfterDraft(plan, doc.Path),
	}, nil
}

func selectFocusDocument(plan PlanOutput, explicitPath string) (DocumentAdvice, error) {
	if explicitPath != "" {
		for _, doc := range plan.RecommendedDocuments {
			if doc.Path == explicitPath {
				return doc, nil
			}
		}
		return DocumentAdvice{}, fmt.Errorf("document %q is not present in current recommended_documents", explicitPath)
	}
	if plan.CurrentPriority == nil {
		return DocumentAdvice{}, fmt.Errorf("no current priority document is available")
	}
	for _, doc := range plan.RecommendedDocuments {
		if doc.Path == plan.CurrentPriority.Path {
			return doc, nil
		}
	}
	return DocumentAdvice{}, fmt.Errorf("current priority document %q is not present in recommended_documents", plan.CurrentPriority.Path)
}

func focusReason(plan PlanOutput, path string, missing []string) string {
	if plan.CurrentPriority != nil && plan.CurrentPriority.Path == path {
		return plan.CurrentPriority.Reason
	}
	if len(missing) > 0 {
		return "explicitly requested document still depends on unresolved context"
	}
	return "explicitly requested document is available for drafting"
}

func availableContextForDocument(ctx Context, required []string) map[string]string {
	out := map[string]string{}
	addIfPresent := func(key, value string) {
		if value != "" {
			out[key] = value
		}
	}

	addIfPresent("project_name", ctx.Project.Name)
	addIfPresent("project_summary", firstNonEmpty(ctx.Project.Summary, answeredValue(ctx, "project_summary")))
	addIfPresent("deployment_shape", answeredValue(ctx, "deployment_shape"))
	addIfPresent("ownership_model", answeredValue(ctx, "ownership_model"))
	addIfPresent("current_layout_summary", firstNonEmpty(ctx.Structure.CurrentLayoutSummary, answeredValue(ctx, "current_layout_summary")))
	addIfPresent("undocumented_directories", answeredValue(ctx, "undocumented_directories"))
	addIfPresent("active_release_docs_strategy", answeredValue(ctx, "active_release_docs_strategy"))

	if len(required) == 0 {
		return out
	}

	filtered := map[string]string{}
	for _, key := range required {
		if value, ok := out[key]; ok {
			filtered[key] = value
		}
	}
	if value, ok := out["project_name"]; ok {
		filtered["project_name"] = value
	}
	return filtered
}

func suggestedSectionsForDocument(mode, path string) []string {
	switch path {
	case "README.md":
		return []string{"Project Summary", "Goals", "Documentation Layout", "Agent Collaboration"}
	case "AGENTS.md":
		return []string{"Working Rules", "Documentation Shape", "Delivery Expectations"}
	case "CLAUDE.md":
		return []string{"Host Notes", "Repository Expectations", "Documentation Workflow"}
	case "docs/domain-overview.md":
		return []string{"Domain Summary", "Core Terms", "Primary Users", "Boundaries"}
	case "docs/architecture.md":
		if mode == "legacy" {
			return []string{"Current Structure", "Module Responsibilities", "Known Constraints", "Reshape Notes"}
		}
		return []string{"Intended Structure", "Repository Layers", "Collaboration Rules", "Constraints"}
	case "docs/legacy-structure-inventory.md":
		return []string{"Top-Level Directories", "Current Responsibilities", "Missing Explanations", "Risks"}
	case "docs/legacy-reshape-guide.md":
		return []string{"Reshape Goals", "Current Constraints", "Phased Plan", "Documentation Ownership"}
	default:
		if isVersionedReleaseReadme(path) {
			return []string{"Release Scope", "Feature Index", "Review Expectations"}
		}
		if isVersionedFeaturePath(path) {
			return []string{"Feature Summary", "Affected Areas", "Open Questions", "Review Checklist"}
		}
		return []string{"Purpose", "Known Context", "Open Questions", "Next Actions"}
	}
}

func writingRulesForDocument() []string {
	return []string{
		"do not invent unresolved facts",
		"prefer explicit repository rules over implicit assumptions",
		"keep reusable language generic and publishable",
		"if context is missing, leave a clear placeholder instead of guessing",
	}
}

func nextAfterDraft(plan PlanOutput, path string) []string {
	next := []string{
		fmt.Sprintf("mark %s as generated in context", path),
		"rerun plan to refresh current_priority",
	}

	if plan.CurrentPriority != nil && plan.CurrentPriority.Path == path {
		for _, doc := range plan.RecommendedDocuments {
			if doc.Path == path {
				continue
			}
			next = append(next, fmt.Sprintf("continue with %s after refresh", doc.Path))
			break
		}
	}
	return next
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
