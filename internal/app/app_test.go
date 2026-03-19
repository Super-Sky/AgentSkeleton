package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestRunInitDocsResolvesProjectDefaultContextAndOutputDir(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectDir := filepath.Join(root, "project")

	err := runInitDocs([]string{
		"--project", projectDir,
		"--output-dir", "generated-docs",
		"--name", "MallHub",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runInitDocs() error = %v", err)
	}

	contextPath := filepath.Join(projectDir, "generated-docs", ".agentskeleton", "context.yaml")
	ctx, err := loadContext(contextPath)
	if err != nil {
		t.Fatalf("loadContext() error = %v", err)
	}

	expectedOutput := filepath.Join(projectDir, "generated-docs")
	if ctx.Paths.ProjectRoot != projectDir {
		t.Fatalf("project_root = %q, want %q", ctx.Paths.ProjectRoot, projectDir)
	}
	if ctx.Paths.OutputDir != expectedOutput {
		t.Fatalf("output_dir = %q, want %q", ctx.Paths.OutputDir, expectedOutput)
	}
	expectedArtifact := filepath.Join(expectedOutput, ".agentskeleton")
	if ctx.Paths.ArtifactDir != expectedArtifact {
		t.Fatalf("artifact_dir = %q, want %q", ctx.Paths.ArtifactDir, expectedArtifact)
	}
	if ctx.Paths.ContextPath != contextPath {
		t.Fatalf("context_path = %q, want %q", ctx.Paths.ContextPath, contextPath)
	}
	for _, p := range ctx.Documentation.MissingDocs {
		if !strings.HasPrefix(p, expectedOutput) {
			t.Fatalf("missing doc path not rooted in output dir: %q", p)
		}
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

func TestRunResponseApplyAcceptUpdatesMultipleAnswers(t *testing.T) {
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
			MissingDocs:   []string{"README.md", "docs/architecture.md"},
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary", "deployment_shape"},
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
		"  deployment_shape: web platform\n" +
		"errors: []\n" +
		"raw_text: \"\"\n"
	if err := os.WriteFile(respPath, []byte(resp), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := runResponse([]string{
		"--file", respPath,
		"--context", contextPath,
		"--apply",
		"--docs", "README.md,docs/architecture.md",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runResponse() error = %v", err)
	}

	updated, err := loadContext(contextPath)
	if err != nil {
		t.Fatalf("loadContext() error = %v", err)
	}

	gotAnswers := map[string]string{}
	for _, a := range updated.Conversation.AnsweredQuestions {
		gotAnswers[a.ID] = a.Value
	}
	if gotAnswers["project_summary"] != "summary text" {
		t.Fatalf("project_summary answer mismatch: %q", gotAnswers["project_summary"])
	}
	if gotAnswers["deployment_shape"] != "web platform" {
		t.Fatalf("deployment_shape answer mismatch: %q", gotAnswers["deployment_shape"])
	}
	if strings.Contains(strings.Join(updated.Conversation.OpenQuestions, ","), "project_summary") {
		t.Fatalf("project_summary should be removed from open_questions")
	}
	if strings.Contains(strings.Join(updated.Conversation.OpenQuestions, ","), "deployment_shape") {
		t.Fatalf("deployment_shape should be removed from open_questions")
	}
	if len(updated.Documentation.GeneratedDocs) != 2 {
		t.Fatalf("generated_docs len = %d, want 2", len(updated.Documentation.GeneratedDocs))
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

func TestRunWorkflowInitial(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary", "deployment_shape"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}

	err := runWorkflow([]string{
		"--context", contextPath,
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runWorkflow() error = %v", err)
	}
}

func TestRunWorkflowApplyAcceptedResponse(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	contextPath := filepath.Join(root, ".agentskeleton", "context.yaml")
	respPath := filepath.Join(root, "response.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:         "discovery",
			GeneratedDocs: []string{},
			MissingDocs:   []string{"README.md"},
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
		"  project_summary: workflow summary\n" +
		"errors: []\n" +
		"raw_text: \"\"\n"
	if err := os.WriteFile(respPath, []byte(resp), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := runWorkflow([]string{
		"--context", contextPath,
		"--response-file", respPath,
		"--apply",
		"--question", "project_summary",
		"--docs", "README.md",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runWorkflow() error = %v", err)
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
	if len(updated.Documentation.GeneratedDocs) != 1 || updated.Documentation.GeneratedDocs[0] != "README.md" {
		t.Fatalf("generated_docs not updated: %#v", updated.Documentation.GeneratedDocs)
	}
}

func TestRunPlanResolvesContextFromOutputDir(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	outputDir := filepath.Join(projectDir, "docs-generated")
	contextPath := filepath.Join(outputDir, ".agentskeleton", "context.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectDir,
			OutputDir:   outputDir,
			ArtifactDir: filepath.Join(outputDir, ".agentskeleton"),
			ContextPath: contextPath,
		},
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy: "recommended",
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}

	err := runPlan([]string{
		"--project", projectDir,
		"--output-dir", "docs-generated",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runPlan() error = %v", err)
	}
}

func TestRunWorkflowWritePlanFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	outputDir := filepath.Join(root, "output")
	contextPath := filepath.Join(outputDir, ".agentskeleton", "context.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectDir,
			OutputDir:   outputDir,
			ArtifactDir: filepath.Join(outputDir, ".agentskeleton"),
			ContextPath: contextPath,
		},
		Project: Project{
			Name:    "MallHub",
			Summary: "AI-friendly shopping mall platform",
			Mode:    "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy: "recommended",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"ownership_model"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}

	err := runWorkflow([]string{
		"--project", projectDir,
		"--output-dir", outputDir,
		"--write-plan-files",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runWorkflow() error = %v", err)
	}

	for _, path := range []string{
		filepath.Join(outputDir, "README.md"),
		filepath.Join(outputDir, "AGENTS.md"),
		filepath.Join(outputDir, "CLAUDE.md"),
		filepath.Join(outputDir, "docs", "architecture.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected generated file missing: %s", path)
		}
	}
}

func TestRunWorkflowWritePlanFilesDoesNotOverwriteByDefault(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	outputDir := filepath.Join(root, "output")
	contextPath := filepath.Join(outputDir, ".agentskeleton", "context.yaml")
	readmePath := filepath.Join(outputDir, "README.md")

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectDir,
			OutputDir:   outputDir,
			ArtifactDir: filepath.Join(outputDir, ".agentskeleton"),
			ContextPath: contextPath,
		},
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy: "recommended",
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(readmePath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(readmePath, []byte("keep me"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := runWorkflow([]string{
		"--project", projectDir,
		"--output-dir", outputDir,
		"--write-plan-files",
		"--format", "yaml",
	})
	if err != nil {
		t.Fatalf("runWorkflow() error = %v", err)
	}

	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "keep me" {
		t.Fatalf("README.md should not be overwritten by default")
	}
}

func TestExecuteWorkflowAutoRepair(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	outputDir := filepath.Join(root, "output")
	contextPath := filepath.Join(outputDir, ".agentskeleton", "context.yaml")
	responsePath := filepath.Join(root, "invalid-response.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectDir,
			OutputDir:   outputDir,
			ArtifactDir: filepath.Join(outputDir, ".agentskeleton"),
			ContextPath: contextPath,
		},
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy: "recommended",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary", "deployment_shape"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}
	if err := os.WriteFile(responsePath, []byte("status: invalid\nerrors:\n  - missing project_summary\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	out, err := executeWorkflow(workflowConfig{
		ContextPath:  contextPath,
		ProjectRoot:  projectDir,
		OutputDir:    outputDir,
		Schema:       "question-answer-set-v1",
		ResponseFile: responsePath,
		Attempt:      0,
		AutoRepair:   true,
	})
	if err != nil {
		t.Fatalf("executeWorkflow() error = %v", err)
	}

	if out.ResponseEval == nil || out.ResponseEval.Decision != RetryDecisionRetry {
		t.Fatalf("response decision = %+v, want retry", out.ResponseEval)
	}
	if out.AutoRepair == nil {
		t.Fatalf("auto repair output should be present")
	}
	if out.AutoRepair.NextAttempt != 1 {
		t.Fatalf("next_attempt = %d, want 1", out.AutoRepair.NextAttempt)
	}
	if out.AutoRepair.Prompt.Mode != "repair" {
		t.Fatalf("auto repair prompt mode = %q, want repair", out.AutoRepair.Prompt.Mode)
	}
	if !strings.Contains(out.AutoRepair.Prompt.PromptText, "Repair structure only") {
		t.Fatalf("auto repair prompt missing repair instruction: %s", out.AutoRepair.Prompt.PromptText)
	}
}

func TestRunWorkflowAutoRepairJSONOutput(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	outputDir := filepath.Join(root, "output")
	contextPath := filepath.Join(outputDir, ".agentskeleton", "context.yaml")
	responsePath := filepath.Join(root, "invalid-response.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectDir,
			OutputDir:   outputDir,
			ArtifactDir: filepath.Join(outputDir, ".agentskeleton"),
			ContextPath: contextPath,
		},
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy: "recommended",
		},
		Conversation: Conversation{
			OpenQuestions: []string{"project_summary"},
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}
	if err := os.WriteFile(responsePath, []byte("status: invalid\nerrors:\n  - missing project_summary\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	output := captureStdout(t, func() {
		err := runWorkflow([]string{
			"--project", projectDir,
			"--output-dir", outputDir,
			"--response-file", responsePath,
			"--attempt", "0",
			"--auto-repair",
			"--format", "json",
		})
		if err != nil {
			t.Fatalf("runWorkflow() error = %v", err)
		}
	})

	if !strings.Contains(output, "\"auto_repair\"") {
		t.Fatalf("workflow output missing auto_repair block: %s", output)
	}
	if !strings.Contains(output, "\"next_attempt\": 1") {
		t.Fatalf("workflow output missing next_attempt: %s", output)
	}
}

func TestRunWorkflowPersistTrace(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	outputDir := filepath.Join(root, "output")
	contextPath := filepath.Join(outputDir, ".agentskeleton", "context.yaml")

	ctx := Context{
		Version: "v0.0.0",
		Paths: Paths{
			ProjectRoot: projectDir,
			OutputDir:   outputDir,
			ArtifactDir: filepath.Join(outputDir, ".agentskeleton"),
			ContextPath: contextPath,
		},
		Project: Project{
			Name: "MallHub",
			Mode: "new",
		},
		Documentation: Documentation{
			Phase:          "discovery",
			ReleaseVersion: "v0.0.0",
		},
		Structure: Structure{
			Strategy: "recommended",
		},
	}
	if err := writeContext(contextPath, ctx); err != nil {
		t.Fatalf("writeContext() error = %v", err)
	}

	originalNow := nowFunc
	nowFunc = func() time.Time {
		return time.Date(2026, 3, 19, 10, 11, 12, 123000000, time.UTC)
	}
	defer func() { nowFunc = originalNow }()

	output := captureStdout(t, func() {
		err := runWorkflow([]string{
			"--project", projectDir,
			"--output-dir", outputDir,
			"--persist-trace",
			"--format", "yaml",
		})
		if err != nil {
			t.Fatalf("runWorkflow() error = %v", err)
		}
	})

	tracePath := filepath.Join(outputDir, ".agentskeleton", "traces", "workflow-discovery-20260319T101112.123000000Z.yaml")
	data, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), "command: workflow") {
		t.Fatalf("trace file missing workflow payload: %s", string(data))
	}
	if !strings.Contains(output, "trace_path: "+tracePath) {
		t.Fatalf("workflow output missing trace_path: %s", output)
	}
}

func TestPersistWorkflowTraceUsesRequestedFormat(t *testing.T) {
	root := t.TempDir()
	artifactDir := filepath.Join(root, ".agentskeleton")

	originalNow := nowFunc
	nowFunc = func() time.Time {
		return time.Date(2026, 3, 19, 10, 11, 12, 456000000, time.UTC)
	}
	defer func() { nowFunc = originalNow }()

	tracePath, err := persistWorkflowTrace(artifactDir, "json", WorkflowOutput{
		Command:     "workflow",
		ContextPath: filepath.Join(root, ".agentskeleton", "context.yaml"),
		Plan: PlanOutput{
			DocumentationPhase: "planning",
		},
	})
	if err != nil {
		t.Fatalf("persistWorkflowTrace() error = %v", err)
	}

	if filepath.Ext(tracePath) != ".json" {
		t.Fatalf("trace extension = %q, want .json", filepath.Ext(tracePath))
	}
	if !strings.Contains(filepath.Base(tracePath), "workflow-planning-") {
		t.Fatalf("trace file name = %q, want planning phase segment", filepath.Base(tracePath))
	}
	data, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), "\"command\": \"workflow\"") {
		t.Fatalf("trace file missing json payload: %s", string(data))
	}
}

func TestSanitizeTraceSegment(t *testing.T) {
	got := sanitizeTraceSegment("Drafting / Review")
	if got != "drafting-review" {
		t.Fatalf("sanitizeTraceSegment() = %q, want %q", got, "drafting-review")
	}
}

func TestRunResponseApplyRejectsExamplePathByDefault(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	respPath := filepath.Join(root, "resp.yaml")
	resp := "" +
		"status: ok\n" +
		"schema: question-answer-set-v1\n" +
		"data:\n" +
		"  project_summary: blocked write\n" +
		"errors: []\n" +
		"raw_text: \"\"\n"
	if err := os.WriteFile(respPath, []byte(resp), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	err := runResponse([]string{
		"--file", respPath,
		"--context", "examples/cli/new-project/context.yaml",
		"--apply",
		"--format", "yaml",
	})
	if err == nil {
		t.Fatalf("runResponse() expected error for example context path")
	}
	if !strings.Contains(err.Error(), "refusing to write example context path") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = writer

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(reader)
		done <- buf.String()
	}()

	fn()

	_ = writer.Close()
	os.Stdout = original
	return <-done
}
