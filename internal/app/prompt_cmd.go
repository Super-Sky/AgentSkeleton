package app

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type PromptOutput struct {
	Command    string   `yaml:"command" json:"command"`
	Mode       string   `yaml:"mode" json:"mode"`
	Schema     string   `yaml:"schema" json:"schema"`
	Questions  []string `yaml:"questions" json:"questions"`
	PromptText string   `yaml:"prompt_text" json:"prompt_text"`
}

func runPrompt(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("prompt", flag.ContinueOnError)
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	mode := fs.String("mode", "initial", "prompt mode: initial|repair")
	schema := fs.String("schema", "question-answer-set-v1", "response schema name")
	errorsRaw := fs.String("errors", "", "comma-separated validation errors for repair mode")
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

	out := PromptOutput{
		Command:   "prompt",
		Mode:      *mode,
		Schema:    *schema,
		Questions: append([]string{}, ctx.Conversation.OpenQuestions...),
	}

	switch *mode {
	case "initial":
		out.PromptText = buildInitialPrompt(ctx, *schema)
	case "repair":
		errs := parseDocs(*errorsRaw)
		if len(errs) == 0 {
			return errors.New("repair mode requires --errors")
		}
		out.PromptText = buildRepairPrompt(ctx, *schema, errs)
	default:
		return fmt.Errorf("unsupported prompt mode %q", *mode)
	}

	return printOutput(*format, out)
}

func buildInitialPrompt(ctx Context, schema string) string {
	var b strings.Builder
	b.WriteString("Based on this context, answer only the missing fields.\n\n")
	b.WriteString("Missing fields:\n")
	for _, q := range ctx.Conversation.OpenQuestions {
		b.WriteString("- " + q + "\n")
	}
	b.WriteString("\nRespond with this envelope:\n")
	b.WriteString("status: ok\n")
	b.WriteString("schema: " + schema + "\n")
	b.WriteString("data:\n")
	for _, q := range ctx.Conversation.OpenQuestions {
		b.WriteString("  " + q + ": \"\"\n")
	}
	b.WriteString("errors: []\n")
	b.WriteString("raw_text: \"\"\n")
	return strings.TrimSpace(b.String())
}

func buildRepairPrompt(ctx Context, schema string, errs []string) string {
	var b strings.Builder
	b.WriteString("Your last response failed schema validation. Repair structure only.\n\n")
	b.WriteString("Validation errors:\n")
	for _, e := range errs {
		b.WriteString("- " + e + "\n")
	}
	b.WriteString("\nOpen questions to repair:\n")
	for _, q := range ctx.Conversation.OpenQuestions {
		b.WriteString("- " + q + "\n")
	}
	b.WriteString("\nRespond again with this envelope:\n")
	b.WriteString("status: ok\n")
	b.WriteString("schema: " + schema + "\n")
	b.WriteString("data:\n")
	for _, q := range ctx.Conversation.OpenQuestions {
		b.WriteString("  " + q + ": \"\"\n")
	}
	b.WriteString("errors: []\n")
	b.WriteString("raw_text: \"\"\n")
	return strings.TrimSpace(b.String())
}
