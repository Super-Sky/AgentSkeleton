package app

import (
	"errors"
	"fmt"
)

// Run dispatches CLI commands.
func Run(args []string) error {
	if len(args) == 0 {
		return usageError()
	}

	switch args[0] {
	case "init-docs":
		return runInitDocs(args[1:])
	case "reshape-docs":
		return runReshapeDocs(args[1:])
	case "plan":
		return runPlan(args[1:])
	case "next":
		return runNext(args[1:])
	case "response":
		return runResponse(args[1:])
	case "prompt":
		return runPrompt(args[1:])
	case "focus-doc":
		return runFocusDoc(args[1:])
	case "workflow":
		return runWorkflow(args[1:])
	case "update":
		return runUpdate(args[1:])
	case "version", "-v", "--version":
		return runVersion(args[1:])
	case "help", "-h", "--help":
		return printUsage()
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usageText)
	}
}

func usageError() error {
	return errors.New(usageText)
}

func printUsage() error {
	fmt.Print(usageText)
	return nil
}

const usageText = `AgentSkeleton CLI

Usage:
  agentskeleton init-docs [flags]
  agentskeleton reshape-docs [flags]
  agentskeleton plan [flags]
  agentskeleton next [flags]
  agentskeleton response [flags]
  agentskeleton prompt [flags]
  agentskeleton focus-doc [flags]
  agentskeleton workflow [flags]
  agentskeleton update [flags]
  agentskeleton version [flags]

Flags:
  --context <path>   Override context file path
  --project <path>   Project root (default: current directory)
  --output-dir <dir> Documentation output directory
  --format <type>    Output format: yaml or json
`
