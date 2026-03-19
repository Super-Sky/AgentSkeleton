package app

import (
	"flag"
	"fmt"
	"os"
)

type WorkflowOutput struct {
	Command      string       `yaml:"command" json:"command"`
	ContextPath  string       `yaml:"context_path" json:"context_path"`
	Plan         PlanOutput   `yaml:"plan" json:"plan"`
	Prompt       PromptOutput `yaml:"prompt" json:"prompt"`
	Next         NextOutput   `yaml:"next" json:"next"`
	ResponseEval *RetryResult `yaml:"response_eval,omitempty" json:"response_eval,omitempty"`
}

func runWorkflow(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("workflow", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	outputDir := fs.String("output-dir", "", "documentation output directory (default: project root)")
	format := fs.String("format", "yaml", "output format")
	schema := fs.String("schema", "question-answer-set-v1", "response schema name")
	responseFile := fs.String("response-file", "", "optional host-model response file (yaml|json)")
	attempt := fs.Int("attempt", 0, "current retry attempt (0-based)")
	apply := fs.Bool("apply", false, "apply accepted response into context")
	allowExampleWrite := fs.Bool("allow-example-write", false, "allow writing context under examples/")
	question := fs.String("question", "", "question id to update in context")
	docs := fs.String("docs", "", "comma-separated docs to mark as generated when accepted")
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

	var responseEval *RetryResult
	if *responseFile != "" {
		data, err := os.ReadFile(*responseFile)
		if err != nil {
			return fmt.Errorf("read response file: %w", err)
		}
		envelope, err := parseEnvelope(data)
		if err != nil {
			return err
		}
		result := EvaluateResponse(DefaultRetryPolicy(), *attempt, envelope)
		responseEval = &result
		if *apply {
			if _, err := applyAcceptedResponse(resolvedContext, *question, *docs, *allowExampleWrite, result, envelope); err != nil {
				return err
			}
		}
	}

	ctx, err := loadContext(resolvedContext)
	if err != nil {
		return err
	}
	if err := validateMode(ctx.Project.Mode); err != nil {
		return err
	}

	plan := PlanOutput{
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

	promptMode := "initial"
	promptText := buildInitialPrompt(ctx, *schema)
	if responseEval != nil && responseEval.Decision == RetryDecisionRetry {
		promptMode = "repair"
		promptText = buildRepairPrompt(ctx, *schema, responseEval.ValidationErrs)
	}
	prompt := PromptOutput{
		Command:    "prompt",
		Mode:       promptMode,
		Schema:     *schema,
		Questions:  append([]string{}, ctx.Conversation.OpenQuestions...),
		PromptText: promptText,
	}

	next := NextOutput{
		Command:            "next",
		DocumentationPhase: ctx.Documentation.Phase,
		ConversationGoal:   conversationGoalForMode(ctx.Project.Mode),
		Questions:          questionsForMode(ctx.Project.Mode),
	}

	out := WorkflowOutput{
		Command:      "workflow",
		ContextPath:  resolvedContext,
		Plan:         plan,
		Prompt:       prompt,
		Next:         next,
		ResponseEval: responseEval,
	}
	return printOutput(*format, out)
}
