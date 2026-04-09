package app

import "flag"

type VersionOutput struct {
	Command string `yaml:"command" json:"command"`
	Version string `yaml:"version" json:"version"`
	Commit  string `yaml:"commit,omitempty" json:"commit,omitempty"`
	Date    string `yaml:"date,omitempty" json:"date,omitempty"`
}

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

func runVersion(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	format := fs.String("format", "yaml", "output format")
	if err := fs.Parse(args); err != nil {
		return err
	}

	return printOutput(*format, VersionOutput{
		Command: "version",
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	})
}
