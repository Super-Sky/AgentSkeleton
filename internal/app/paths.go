package app

import (
	"path/filepath"
	"strings"
)

const defaultProjectRoot = "."

func resolveProjectRoot(project string) (string, error) {
	if project == "" {
		project = defaultProjectRoot
	}
	return filepath.Abs(project)
}

func resolveOutputDir(projectRoot, outputDir string) (string, error) {
	if outputDir == "" {
		return projectRoot, nil
	}
	if filepath.IsAbs(outputDir) {
		return filepath.Clean(outputDir), nil
	}
	return filepath.Abs(filepath.Join(projectRoot, outputDir))
}

func artifactDirForOutput(outputDir string) string {
	return filepath.Join(outputDir, ".agentskeleton")
}

func defaultContextPathForOutput(outputDir string) string {
	return filepath.Join(artifactDirForOutput(outputDir), "context.yaml")
}

func resolveContextPath(outputDir, explicitContext string, contextExplicitlySet bool) string {
	if contextExplicitlySet {
		return explicitContext
	}
	return defaultContextPathForOutput(outputDir)
}

func flagExplicitlySet(args []string, longName string) bool {
	long := "--" + longName
	for i := range args {
		if args[i] == long || strings.HasPrefix(args[i], long+"=") {
			return true
		}
	}
	return false
}

func resolveDocPath(baseDir, rel string) string {
	if baseDir == "" {
		return rel
	}
	return filepath.Join(baseDir, rel)
}
