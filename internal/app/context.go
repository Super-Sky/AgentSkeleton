package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

const defaultContextPath = ".agentskeleton/context.yaml"

// Context stores the project guidance state.
type Context struct {
	Version       string        `yaml:"version" json:"version"`
	Project       Project       `yaml:"project" json:"project"`
	Documentation Documentation `yaml:"documentation" json:"documentation"`
	Structure     Structure     `yaml:"structure" json:"structure"`
	Conversation  Conversation  `yaml:"conversation" json:"conversation"`
}

type Project struct {
	Name         string   `yaml:"name" json:"name"`
	Summary      string   `yaml:"summary" json:"summary"`
	Mode         string   `yaml:"mode" json:"mode"`
	Domain       string   `yaml:"domain" json:"domain"`
	PrimaryUsers []string `yaml:"primary_users" json:"primary_users"`
	Host         string   `yaml:"host" json:"host"`
}

type Documentation struct {
	Phase          string   `yaml:"phase" json:"phase"`
	GeneratedDocs  []string `yaml:"generated_docs" json:"generated_docs"`
	MissingDocs    []string `yaml:"missing_docs" json:"missing_docs"`
	ReleaseVersion string   `yaml:"release_version" json:"release_version"`
}

type Structure struct {
	Strategy             string `yaml:"strategy" json:"strategy"`
	RecommendedLayout    string `yaml:"recommended_layout" json:"recommended_layout"`
	CurrentLayoutSummary string `yaml:"current_layout_summary" json:"current_layout_summary"`
}

type Conversation struct {
	AnsweredQuestions []QuestionAnswer `yaml:"answered_questions" json:"answered_questions"`
	OpenQuestions     []string         `yaml:"open_questions" json:"open_questions"`
}

type QuestionAnswer struct {
	ID    string `yaml:"id" json:"id"`
	Value string `yaml:"value" json:"value"`
}

func loadContext(path string) (Context, error) {
	var ctx Context

	data, err := os.ReadFile(path)
	if err != nil {
		return ctx, fmt.Errorf("read context: %w", err)
	}
	if err := yaml.Unmarshal(data, &ctx); err != nil {
		return ctx, fmt.Errorf("parse context: %w", err)
	}
	return ctx, nil
}

func writeContext(path string, ctx Context) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create context dir: %w", err)
	}

	data, err := yaml.Marshal(ctx)
	if err != nil {
		return fmt.Errorf("marshal context: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write context: %w", err)
	}
	return nil
}

func validateMode(mode string) error {
	switch mode {
	case "new", "legacy":
		return nil
	default:
		return errors.New("project.mode must be new or legacy")
	}
}

func (c *Context) applyAnswer(questionID, value string) {
	if questionID == "" {
		return
	}

	replaced := false
	for i := range c.Conversation.AnsweredQuestions {
		if c.Conversation.AnsweredQuestions[i].ID == questionID {
			c.Conversation.AnsweredQuestions[i].Value = value
			replaced = true
			break
		}
	}
	if !replaced {
		c.Conversation.AnsweredQuestions = append(c.Conversation.AnsweredQuestions, QuestionAnswer{
			ID:    questionID,
			Value: value,
		})
	}
	c.Conversation.OpenQuestions = removeString(c.Conversation.OpenQuestions, questionID)
}

func (c *Context) markGenerated(paths []string) {
	for _, p := range paths {
		if p == "" {
			continue
		}
		if !slices.Contains(c.Documentation.GeneratedDocs, p) {
			c.Documentation.GeneratedDocs = append(c.Documentation.GeneratedDocs, p)
		}
		c.Documentation.MissingDocs = removeString(c.Documentation.MissingDocs, p)
	}
}

func removeString(items []string, target string) []string {
	if target == "" {
		return items
	}
	out := make([]string, 0, len(items))
	for _, it := range items {
		if it != target {
			out = append(out, it)
		}
	}
	return out
}
