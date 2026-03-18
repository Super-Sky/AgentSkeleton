package app

import "fmt"

// ResponseEnvelope is the normalized host-model response shape.
type ResponseEnvelope struct {
	Status  string         `yaml:"status" json:"status"`
	Schema  string         `yaml:"schema" json:"schema"`
	Data    map[string]any `yaml:"data" json:"data"`
	Errors  []string       `yaml:"errors" json:"errors"`
	RawText string         `yaml:"raw_text" json:"raw_text"`
}

// Validate checks whether the response envelope can be accepted into the next step.
func (r ResponseEnvelope) Validate() error {
	switch r.Status {
	case "ok":
		if r.Schema == "" {
			return fmt.Errorf("schema is required when status is ok")
		}
		if len(r.Errors) > 0 {
			return fmt.Errorf("errors must be empty when status is ok")
		}
	case "invalid":
		if len(r.Errors) == 0 {
			return fmt.Errorf("errors are required when status is invalid")
		}
	case "unresolved":
		if r.RawText == "" && len(r.Errors) == 0 {
			return fmt.Errorf("raw_text or errors are required when status is unresolved")
		}
	default:
		return fmt.Errorf("unsupported status %q", r.Status)
	}

	return nil
}
