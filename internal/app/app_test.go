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
