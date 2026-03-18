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
	fs := flag.NewFlagSet("response", flag.ContinueOnError)
	file := fs.String("file", "", "host-model response file (yaml|json)")
	attempt := fs.Int("attempt", 0, "current retry attempt (0-based)")
	contextPath := fs.String("context", defaultContextPath, "context file path")
	apply := fs.Bool("apply", false, "apply accepted response into context")
	question := fs.String("question", "", "question id to update in context")
	docs := fs.String("docs", "", "comma-separated docs to mark as generated when accepted")
	format := fs.String("format", "yaml", "output format")
	if err := fs.Parse(args); err != nil {
		return err
	}

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
		updated, err := applyAcceptedResponse(*contextPath, *question, *docs, result, envelope)
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

func applyAcceptedResponse(contextPath, question, docs string, result RetryResult, envelope ResponseEnvelope) (bool, error) {
	if result.Decision != RetryDecisionAccept {
		return false, nil
	}

	ctx, err := loadContext(contextPath)
	if err != nil {
		return false, err
	}

	questionID := question
	if questionID == "" {
		questionID, err = inferQuestionID(envelope.Data)
		if err != nil {
			return false, err
		}
	}

	value, ok := envelope.Data[questionID]
	if !ok {
		return false, fmt.Errorf("response data does not include question key %q", questionID)
	}
	ctx.applyAnswer(questionID, fmt.Sprint(value))
	ctx.markGenerated(parseDocs(docs))

	if err := writeContext(contextPath, ctx); err != nil {
		return false, err
	}
	return true, nil
}

func inferQuestionID(data map[string]any) (string, error) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	if len(keys) == 1 {
		return keys[0], nil
	}
	if len(keys) == 0 {
		return "", errors.New("cannot infer question id from empty response data")
	}
	slices.Sort(keys)
	return "", fmt.Errorf("cannot infer question id from multiple keys: %v", keys)
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
