package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// DiscoveredProject represents a project detected from files in the workspace.
type DiscoveredProject struct {
	Name       string
	Path       string
	Language   string
	DetectFile string
}

// DetectProjectType reports whether a directory contains a recognized project file.
func DetectProjectType(path string) bool {
	_, _, ok := detectProjectLanguage(path)
	return ok
}

func detectProjectLanguage(path string) (string, string, bool) {
	if fileExists(filepath.Join(path, "go.mod")) {
		return "go", "go.mod", true
	}
	if fileExists(filepath.Join(path, "Cargo.toml")) {
		return "rust", "Cargo.toml", true
	}
	if fileExists(filepath.Join(path, "package.json")) {
		return "node", "package.json", true
	}
	if fileExists(filepath.Join(path, "Makefile")) {
		return "make", "Makefile", true
	}
	return "", "", false
}

// ScanProjects recursively scans a workspace root and returns detected projects.
func ScanProjects(root string) ([]DiscoveredProject, error) {
	if root == "" {
		root = "."
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", root)
	}

	projects := []DiscoveredProject{}
	seen := map[string]struct{}{}

	if err := filepath.WalkDir(root, func(current string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d == nil || !d.IsDir() {
			return nil
		}
		if shouldSkipDir(d.Name()) {
			return filepath.SkipDir
		}

		lang, detectFile, ok := detectProjectLanguage(current)
		if !ok {
			return nil
		}

		key := filepath.Clean(current)
		if _, exists := seen[key]; exists {
			return nil
		}
		seen[key] = struct{}{}

		projects = append(projects, DiscoveredProject{
			Name:       filepath.Base(current),
			Path:       current,
			Language:   lang,
			DetectFile: detectFile,
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return projects, nil
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", ".hg", ".svn":
		return true
	case ".gobuild", ".idea", ".vscode", "node_modules", "vendor", "dist", "build", "bin", "target", "coverage":
		return true
	default:
		if len(name) > 0 && name[0] == '.' {
			return true
		}
		return false
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
