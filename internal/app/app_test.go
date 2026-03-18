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

func TestRunResponseApplyAcceptUpdatesContext(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")
	respPath := filepath.Join(root, "resp.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:         "discovery",
			GeneratedDocs: []string{},
			MissingDocs:   []string{"docs/domain-overview.md"},
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}
	resp := "" +
		"status: ok\n" +
		"schema: question-answer-set-v1\n" +
		"data:\n" +
		"  project_summary: summary text\n" +
		"errors: []\n" +
		"raw_text: \"\"\n"
	if err := os.WriteFile(respPath, []byte(resp), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := runResponse([]string{
		"--file", respPath,
		"--context", contextPath,
		"--apply",
		"--question", "project_summary",
		"--docs", "docs/domain-overview.md",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runResponse() error = %v", err)
	}

	updated, err := loadContext(contextPath)
	if err != nil {
		t.Fatalf("loadContext() error = %v", err)
	}
	if len(updated.Conversation.AnsweredQuestions) != 1 {
		t.Fatalf("answered_questions len = %d, want 1", len(updated.Conversation.AnsweredQuestions))
	}
	if updated.Conversation.AnsweredQuestions[0].ID != "project_summary" {
		t.Fatalf("answered question id = %q", updated.Conversation.AnsweredQuestions[0].ID)
	}
	if strings.Contains(strings.Join(updated.Conversation.OpenQuestions, ","), "project_summary") {
		t.Fatalf("project_summary should be removed from open_questions")
	}
	if len(updated.Documentation.GeneratedDocs) != 1 || updated.Documentation.GeneratedDocs[0] != "docs/domain-overview.md" {
		t.Fatalf("generated_docs not updated: %#v", updated.Documentation.GeneratedDocs)
	}
	if strings.Contains(strings.Join(updated.Documentation.MissingDocs, ","), "docs/domain-overview.md") {
		t.Fatalf("docs/domain-overview.md should be removed from missing_docs")
	}
}

func TestRunResponseApplyInvalidDoesNotUpdateContext(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")
	respPath := filepath.Join(root, "resp.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}
	before, err := os.ReadFile(contextPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	resp := "" +
		"status: invalid\n" +
		"schema: question-answer-set-v1\n" +
		"data: {}\n" +
		"errors:\n" +
		"  - \"missing required field: project_summary\"\n" +
		"raw_text: \"\"\n"
	if err := os.WriteFile(respPath, []byte(resp), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err = runResponse([]string{
		"--file", respPath,
		"--context", contextPath,
		"--apply",
		"--question", "project_summary",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runResponse() error = %v", err)
	}

	after, err := os.ReadFile(contextPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("context should not change on invalid response")
	}
}

func TestRunPromptInitial(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary", "deployment_shape"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}

	err := runPrompt([]string{
		"--context", contextPath,
		"--mode", "initial",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runPrompt() error = %v", err)
	}
}

func TestRunPromptRepairRequiresErrors(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}

	err := runPrompt([]string{
		"--context", contextPath,
		"--mode", "repair",
		"--format", "yaml",
	})
	if err == nil {
		t.Fatalf("runPrompt() expected error in repair mode without --errors")
	}
}
