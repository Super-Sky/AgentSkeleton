package app

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

func runResponse(args []string) error {
	contextSet := flagExplicitlySet(args, "context")
	fs := flag.NewFlagSet("response", flag.ContinueOnError)
	file := fs.String("file", "", "host-model response file (yaml|json)")
	attempt := fs.Int("attempt", 0, "current retry attempt (0-based)")
	contextPath := fs.String("context", defaultContextPath, "context file path")
	project := fs.String("project", defaultProjectRoot, "project root path")
	outputDir := fs.String("output-dir", "", "documentation output directory (default: project root)")
	apply := fs.Bool("apply", false, "apply accepted response into context")
	allowExampleWrite := fs.Bool("allow-example-write", false, "allow writing context under examples/")
	question := fs.String("question", "", "question id to update in context")
	docs := fs.String("docs", "", "comma-separated docs to mark as generated when accepted")
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

	if *file == "" {
		return errors.New("response file is required")
	}

	data, err := os.ReadFile(*file)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	envelope, err := parseEnvelope(data)
	if err != nil {
		return err
	}

	result := EvaluateResponse(DefaultRetryPolicy(), *attempt, envelope)

	output := map[string]any{
		"command": "response",
		"result":  result,
	}

	if *apply {
		updated, err := applyAcceptedResponse(resolvedContext, *question, *docs, *allowExampleWrite, result, envelope)
		if err != nil {
			return err
		}
		output["context_updated"] = updated
	}

	return printOutput(*format, output)
}

func parseEnvelope(data []byte) (ResponseEnvelope, error) {
	var envelope ResponseEnvelope
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return envelope, errors.New("empty response payload")
	}

	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		if err := json.Unmarshal([]byte(trimmed), &envelope); err != nil {
			return envelope, fmt.Errorf("parse response as json: %w", err)
		}
		return envelope, nil
	}

	if err := yaml.Unmarshal([]byte(trimmed), &envelope); err != nil {
		return envelope, fmt.Errorf("parse response as yaml: %w", err)
	}
	return envelope, nil
}

func applyAcceptedResponse(contextPath, question, docs string, allowExampleWrite bool, result RetryResult, envelope ResponseEnvelope) (bool, error) {
	if result.Decision != RetryDecisionAccept {
		return false, nil
	}
	if !allowExampleWrite && isExamplePath(contextPath) {
		return false, fmt.Errorf("refusing to write example context path %q (use --allow-example-write to override)", contextPath)
	}

	ctx, err := loadContext(contextPath)
	if err != nil {
		return false, err
	}

	questionIDs, err := selectQuestionIDs(question, envelope.Data)
	if err != nil {
		return false, err
	}
	for _, questionID := range questionIDs {
		value, ok := envelope.Data[questionID]
		if !ok {
			return false, fmt.Errorf("response data does not include question key %q", questionID)
		}
		ctx.applyAnswer(questionID, fmt.Sprint(value))
	}
	ctx.markGenerated(parseDocs(docs))

	if err := writeContext(contextPath, ctx); err != nil {
		return false, err
	}
	return true, nil
}

func isExamplePath(path string) bool {
	normalized := strings.ReplaceAll(path, "\\", "/")
	return strings.Contains(normalized, "/examples/") || strings.HasPrefix(normalized, "examples/")
}

func selectQuestionIDs(explicitQuestion string, data map[string]any) ([]string, error) {
	if explicitQuestion != "" {
		return []string{explicitQuestion}, nil
	}

	keys, err := inferQuestionIDs(data)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func inferQuestionIDs(data map[string]any) ([]string, error) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return nil, errors.New("cannot infer question id from empty response data")
	}
	slices.Sort(keys)
	return keys, nil
}

func parseDocs(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
