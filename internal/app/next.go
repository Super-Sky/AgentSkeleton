package app

import "flag"

type NextOutput struct {
	Command            string         `yaml:"command" json:"command"`
	DocumentationPhase string         `yaml:"documentation_phase" json:"documentation_phase"`
	ConversationGoal   string         `yaml:"conversation_goal" json:"conversation_goal"`
	Questions          []NextQuestion `yaml:"questions" json:"questions"`
}

type NextQuestion struct {
	ID      string   `yaml:"id" json:"id"`
	Prompt  string   `yaml:"prompt" json:"prompt"`
	Reason  string   `yaml:"reason" json:"reason"`
	Affects []string `yaml:"affects" json:"affects"`
}

func runNext(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("next", flag.ContinueOnError)
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

	out := NextOutput{
		Command:            "next",
		DocumentationPhase: ctx.Documentation.Phase,
		ConversationGoal:   conversationGoalForMode(ctx.Project.Mode),
		Questions:          questionsForMode(ctx.Project.Mode),
	}

	return printOutput(*format, out)
}

func conversationGoalForMode(mode string) string {
	if mode == "legacy" {
		return "document the current repository before proposing reshape actions"
	}
	return "clarify the starting documentation scope for a new project"
}

func questionsForMode(mode string) []NextQuestion {
	if mode == "legacy" {
		return []NextQuestion{
			{
				ID:     "undocumented_directories",
				Prompt: "Which top-level directories or major modules are currently undocumented?",
				Reason: "Legacy reshaping should begin by explaining what already exists.",
				Affects: []string{
					"docs/legacy-structure-inventory.md",
					"docs/architecture.md",
				},
			},
			{
				ID:     "active_release_docs_strategy",
				Prompt: "Does this repository need versioned feature documents for current active work?",
				Reason: "This determines whether docs/vX.Y.Z/features should be introduced now.",
				Affects: []string{
					"docs/README.md",
					"docs/v0.0.0/README.md",
				},
			},
			{
				ID:     "ownership_model",
				Prompt: "Who will maintain the reshaped documentation going forward?",
				Reason: "This affects governance and task-delivery expectations.",
				Affects: []string{
					"docs/repo-workflow-guide.md",
					"docs/task-delivery-guide.md",
				},
			},
		}
	}

	return []NextQuestion{
		{
			ID:     "project_summary",
			Prompt: "What is the one-sentence summary of this project?",
			Reason: "The summary anchors README and domain overview drafts.",
			Affects: []string{
				"README.md",
				"docs/domain-overview.md",
			},
		},
		{
			ID:     "deployment_shape",
			Prompt: "What kind of system shape is this project expected to have, such as web platform, CLI, backend service, or mixed system?",
			Reason: "This determines how architecture and repository structure should be described.",
			Affects: []string{
				"docs/architecture.md",
				"README.md",
			},
		},
		{
			ID:     "ownership_model",
			Prompt: "Who is expected to maintain the documentation over time?",
			Reason: "This affects workflow guidance and governance expectations.",
			Affects: []string{
				"docs/repo-workflow-guide.md",
				"docs/task-delivery-guide.md",
			},
		},
	}
}
