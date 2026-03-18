package app

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func printOutput(format string, v any) error {
	switch format {
	case "yaml", "":
		data, err := yaml.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal yaml: %w", err)
		}
		_, err = os.Stdout.Write(data)
		return err
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}
