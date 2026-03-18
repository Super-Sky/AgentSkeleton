package app

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func runResponse(args []string) error {
	fs := flag.NewFlagSet("response", flag.ContinueOnError)
	file := fs.String("file", "", "host-model response file (yaml|json)")
	attempt := fs.Int("attempt", 0, "current retry attempt (0-based)")
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

	var envelope ResponseEnvelope
	if err := yaml.Unmarshal(data, &envelope); err != nil {
		if err := json.Unmarshal(data, &envelope); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
	}

	result := EvaluateResponse(DefaultRetryPolicy(), *attempt, envelope)
	return printOutput(*format, result)
}
