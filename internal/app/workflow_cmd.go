package app

import (
	"flag"
	"fmt"
	"os"
)

type AutoRepairOutput struct {
	NextAttempt      int          `yaml:"next_attempt" json:"next_attempt"`
	ValidationErrors []string     `yaml:"validation_errors" json:"validation_errors"`
	Prompt           PromptOutput `yaml:"prompt" json:"prompt"`
	Instructions     []string     `yaml:"instructions" json:"instructions"`
}

type WorkflowOutput struct {
	Command      string            `yaml:"command" json:"command"`
	ContextPath  string            `yaml:"context_path" json:"context_path"`
	Plan         PlanOutput        `yaml:"plan" json:"plan"`
	Prompt       PromptOutput      `yaml:"prompt" json:"prompt"`
	Next         NextOutput        `yaml:"next" json:"next"`
	ResponseEval *RetryResult      `yaml:"response_eval,omitempty" json:"response_eval,omitempty"`
	WriteResult  *writeResult      `yaml:"write_result,omitempty" json:"write_result,omitempty"`
	AutoRepair   *AutoRepairOutput `yaml:"auto_repair,omitempty" json:"auto_repair,omitempty"`
}

type workflowConfig struct {
	ContextPath       string
	ProjectRoot       string
	OutputDir         string
	Schema            string
	ResponseFile      string
	Attempt           int
	Apply             bool
	AllowExampleWrite bool
	WritePlanFiles    bool
	Overwrite         bool
	Question          string
	Docs              string
	AutoRepair        bool
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
	writePlanFiles := fs.Bool("write-plan-files", false, "write planned document skeleton files into output-dir")
	overwrite := fs.Bool("overwrite", false, "overwrite existing generated plan files")
	autoRepair := fs.Bool("auto-repair", false, "when decision is retry, emit a repair package for the next host-model turn")
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

	out, err := executeWorkflow(workflowConfig{
		ContextPath:       resolvedContext,
		ProjectRoot:       projectRoot,
		OutputDir:         outputRoot,
		Schema:            *schema,
		ResponseFile:      *responseFile,
		Attempt:           *attempt,
		Apply:             *apply,
		AllowExampleWrite: *allowExampleWrite,
		WritePlanFiles:    *writePlanFiles,
		Overwrite:         *overwrite,
		Question:          *question,
		Docs:              *docs,
		AutoRepair:        *autoRepair,
	})
	if err != nil {
		return err
	}

	return printOutput(*format, out)
}

func executeWorkflow(cfg workflowConfig) (WorkflowOutput, error) {
	var responseEval *RetryResult
	if cfg.ResponseFile != "" {
		data, err := os.ReadFile(cfg.ResponseFile)
		if err != nil {
			return WorkflowOutput{}, fmt.Errorf("read response file: %w", err)
		}
		envelope, err := parseEnvelope(data)
		if err != nil {
			return WorkflowOutput{}, err
		}
		result := EvaluateResponse(DefaultRetryPolicy(), cfg.Attempt, envelope)
		responseEval = &result
		if cfg.Apply {
			if _, err := applyAcceptedResponse(cfg.ContextPath, cfg.Question, cfg.Docs, cfg.AllowExampleWrite, result, envelope); err != nil {
				return WorkflowOutput{}, err
			}
		}
	}

	ctx, err := loadContext(cfg.ContextPath)
	if err != nil {
		return WorkflowOutput{}, err
	}
	if err := validateMode(ctx.Project.Mode); err != nil {
		return WorkflowOutput{}, err
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
	promptText := buildInitialPrompt(ctx, cfg.Schema)
	if responseEval != nil && responseEval.Decision == RetryDecisionRetry {
		promptMode = "repair"
		promptText = buildRepairPrompt(ctx, cfg.Schema, responseEval.ValidationErrs)
	}
	prompt := PromptOutput{
		Command:    "prompt",
		Mode:       promptMode,
		Schema:     cfg.Schema,
		Questions:  append([]string{}, ctx.Conversation.OpenQuestions...),
		PromptText: promptText,
	}

	next := NextOutput{
		Command:            "next",
		DocumentationPhase: ctx.Documentation.Phase,
		ConversationGoal:   conversationGoalForMode(ctx.Project.Mode),
		Questions:          questionsForMode(ctx.Project.Mode),
	}

	var fileWriteResult *writeResult
	if cfg.WritePlanFiles {
		res, err := renderPlannedFiles(ctx, plan, cfg.Overwrite)
		if err != nil {
			return WorkflowOutput{}, err
		}
		fileWriteResult = &res
	}

	var autoRepairOut *AutoRepairOutput
	if cfg.AutoRepair && responseEval != nil && responseEval.Decision == RetryDecisionRetry {
		autoRepairOut = buildAutoRepairOutput(ctx, cfg.Schema, cfg.Attempt+1, responseEval.ValidationErrs)
	}

	return WorkflowOutput{
		Command:      "workflow",
		ContextPath:  cfg.ContextPath,
		Plan:         plan,
		Prompt:       prompt,
		Next:         next,
		ResponseEval: responseEval,
		WriteResult:  fileWriteResult,
		AutoRepair:   autoRepairOut,
	}, nil
}

func buildAutoRepairOutput(ctx Context, schema string, nextAttempt int, validationErrs []string) *AutoRepairOutput {
	return &AutoRepairOutput{
		NextAttempt:      nextAttempt,
		ValidationErrors: append([]string{}, validationErrs...),
		Prompt: PromptOutput{
			Command:    "prompt",
			Mode:       "repair",
			Schema:     schema,
			Questions:  append([]string{}, ctx.Conversation.OpenQuestions...),
			PromptText: buildRepairPrompt(ctx, schema, validationErrs),
		},
		Instructions: []string{
			"ask the host model to repair structure only and keep semantic intent unchanged",
			"validate the repaired response with the incremented attempt value before applying it",
			"only write accepted responses back into <output-dir>/.agentskeleton/context.yaml",
		},
	}
}
