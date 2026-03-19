package app

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func printOutput(format string, v any) error {
	data, err := marshalOutput(format, v)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}

func marshalOutput(format string, v any) ([]byte, error) {
	switch format {
	case "yaml", "":
		data, err := yaml.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("marshal yaml: %w", err)
		}
		return data, nil
	case "json":
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal json: %w", err)
		}
		data = append(data, '\n')
		return data, nil
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
}
