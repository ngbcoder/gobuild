package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectType(t *testing.T) {
	t.Run("go project", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if !DetectProjectType(dir) {
			t.Fatal("DetectProjectType() = false, want true for Go project")
		}
	})

	t.Run("rust project", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[package]\nname = \"demo\"\n"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if !DetectProjectType(dir) {
			t.Fatal("DetectProjectType() = false, want true for Rust project")
		}
	})

	t.Run("node project", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{\n  \"name\": \"demo\"\n}\n"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if !DetectProjectType(dir) {
			t.Fatal("DetectProjectType() = false, want true for Node project")
		}
	})

	t.Run("make project", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "Makefile"), []byte("build:\n\t@echo ok\n"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if !DetectProjectType(dir) {
			t.Fatal("DetectProjectType() = false, want true for Make project")
		}
	})
}

func TestScanProjects(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "backend"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/root\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "backend", "Cargo.toml"), []byte("[package]\nname = \"backend\"\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "frontend"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{\n  \"name\": \"frontend\"\n}\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "Makefile"), []byte("build:\n\t@echo root\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	projects, err := ScanProjects(root)
	if err != nil {
		t.Fatalf("ScanProjects() returned error: %v", err)
	}
	if len(projects) != 3 {
		t.Fatalf("len(projects) = %d, want 3", len(projects))
	}

	seen := map[string]bool{}
	for _, project := range projects {
		seen[project.Language] = true
	}
	if !seen["go"] || !seen["rust"] || !seen["node"] {
		t.Fatalf("project languages = %#v, want go, rust, and node", seen)
	}
}
