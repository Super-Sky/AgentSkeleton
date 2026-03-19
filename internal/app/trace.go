package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var nowFunc = time.Now

func persistWorkflowTrace(artifactDir, format string, output WorkflowOutput) (string, error) {
	if artifactDir == "" {
		return "", fmt.Errorf("persist workflow trace: artifact dir is required")
	}

	traceDir := filepath.Join(artifactDir, "traces")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return "", fmt.Errorf("create trace dir: %w", err)
	}

	ext := "yaml"
	if format == "json" {
		ext = "json"
	}

	phase := sanitizeTraceSegment(output.Plan.DocumentationPhase)
	if phase == "" {
		phase = "unknown"
	}
	filename := fmt.Sprintf("workflow-%s-%s.%s", phase, nowFunc().UTC().Format("20060102T150405.000000000Z"), ext)
	tracePath := filepath.Join(traceDir, filename)

	output.TracePath = tracePath
	data, err := marshalOutput(format, output)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(tracePath, data, 0o644); err != nil {
		return "", fmt.Errorf("write workflow trace: %w", err)
	}

	return tracePath, nil
}

func sanitizeTraceSegment(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return ""
	}
	var b strings.Builder
	lastDash := false
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastDash = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
