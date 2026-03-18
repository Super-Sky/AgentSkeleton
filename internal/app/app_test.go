package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitDocsCreatesContext(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")

	err := runInitDocs([]string{
		"--context", contextPath,
		"--name", "MallHub",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runInitDocs() error = %v", err)
	}

	data, err := os.ReadFile(contextPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), "mode: new") {
		t.Fatalf("context file does not contain new mode:\n%s", string(data))
	}
}

func TestLoadContextFixtures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		mode string
	}{
		{
			name: "new",
			path: filepath.Join("..", "..", "examples", "cli", "new-project", "context.yaml"),
			mode: "new",
		},
		{
			name: "legacy",
			path: filepath.Join("..", "..", "examples", "cli", "legacy-project", "context.yaml"),
			mode: "legacy",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, err := loadContext(tt.path)
			if err != nil {
				t.Fatalf("loadContext() error = %v", err)
			}
			if ctx.Project.Mode != tt.mode {
				t.Fatalf("mode = %q, want %q", ctx.Project.Mode, tt.mode)
			}
		})
	}
}

func TestResponseEnvelopeValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   ResponseEnvelope
		wantErr bool
	}{
		{
			name: "ok",
			input: ResponseEnvelope{
				Status: "ok",
				Schema: "question-answer-set-v1",
				Data: map[string]any{
					"project_summary": "MallHub summary",
				},
			},
		},
		{
			name: "invalid requires errors",
			input: ResponseEnvelope{
				Status: "invalid",
			},
			wantErr: true,
		},
		{
			name: "unresolved with raw text",
			input: ResponseEnvelope{
				Status:  "unresolved",
				RawText: "free-form answer that could not be normalized",
			},
		},
		{
			name: "bad status",
			input: ResponseEnvelope{
				Status: "broken",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEvaluateResponse(t *testing.T) {
	t.Parallel()

	policy := DefaultRetryPolicy()

	tests := []struct {
		name     string
		attempt  int
		response ResponseEnvelope
		want     RetryDecision
	}{
		{
			name:    "accept valid response",
			attempt: 0,
			response: ResponseEnvelope{
				Status: "ok",
				Schema: "question-answer-set-v1",
				Data: map[string]any{
					"project_summary": "MallHub summary",
				},
			},
			want: RetryDecisionAccept,
		},
		{
			name:    "retry invalid response",
			attempt: 1,
			response: ResponseEnvelope{
				Status: "invalid",
				Errors: []string{"missing project_summary"},
			},
			want: RetryDecisionRetry,
		},
		{
			name:    "unresolved after retry budget",
			attempt: 2,
			response: ResponseEnvelope{
				Status: "invalid",
				Errors: []string{"missing project_summary"},
			},
			want: RetryDecisionUnresolved,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := EvaluateResponse(policy, tt.attempt, tt.response)
			if got.Decision != tt.want {
				t.Fatalf("Decision = %q, want %q", got.Decision, tt.want)
			}
		})
	}
}
