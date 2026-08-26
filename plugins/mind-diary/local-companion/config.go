package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const workspaceRootsEnvironment = "MIND_DIARY_WORKSPACE_ROOTS"

// Workspace roots are trusted process configuration, never MCP tool input.
// On macOS the value is a standard colon-separated path list.
func configuredWorkspaceRoots() ([]string, error) {
	value := os.Getenv(workspaceRootsEnvironment)
	if value == "" {
		return nil, nil
	}
	parts := filepath.SplitList(value)
	if len(parts) == 0 {
		return nil, errors.New("workspace root configuration is invalid")
	}
	for _, part := range parts {
		if part == "" || strings.TrimSpace(part) != part {
			return nil, errors.New("workspace root configuration is invalid")
		}
	}
	return canonicalWorkspaceRoots(parts)
}

func canonicalWorkspaceRoots(roots []string) ([]string, error) {
	canonical := make([]string, 0, len(roots))
	seen := make(map[string]struct{}, len(roots))
	for _, root := range roots {
		if root == "" || !filepath.IsAbs(root) {
			return nil, errors.New("workspace roots must be absolute directories")
		}
		resolved, err := filepath.EvalSymlinks(filepath.Clean(root))
		if err != nil || !filepath.IsAbs(resolved) {
			return nil, errors.New("workspace roots must resolve to directories")
		}
		info, err := os.Lstat(resolved)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, errors.New("workspace roots must resolve to directories")
		}
		if _, exists := seen[resolved]; exists {
			continue
		}
		seen[resolved] = struct{}{}
		canonical = append(canonical, resolved)
	}
	return canonical, nil
}
