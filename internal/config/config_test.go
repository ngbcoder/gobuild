package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigMissingFileUsesDefaults(t *testing.T) {
	cfg, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("LoadConfig() returned error for missing file: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}
	if cfg.Version != "1" {
		t.Fatalf("Version = %q, want %q", cfg.Version, "1")
	}
	if cfg.PluginsDir != "~/.gobuild/plugins" {
		t.Fatalf("PluginsDir = %q, want %q", cfg.PluginsDir, "~/.gobuild/plugins")
	}
	if len(cfg.Projects) != 0 {
		t.Fatalf("len(Projects) = %d, want 0", len(cfg.Projects))
	}
}

func TestLoadConfigParsesAndExpandsPaths(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() returned error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "gobuild.yaml")
	content := "version: \"2\"\nplugins_dir: \"~/.gobuild/plugins\"\nprojects:\n  - name: demo\n    path: \"~/workspace/demo\"\n    language: go\n    build_cmd: \"go build\"\n    test_cmd: \"go test ./...\"\n    env:\n      APP_MODE: debug\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}
	if cfg.Version != "2" {
		t.Fatalf("Version = %q, want %q", cfg.Version, "2")
	}
	if cfg.PluginsDir != filepath.Join(home, ".gobuild", "plugins") {
		t.Fatalf("PluginsDir = %q, want %q", cfg.PluginsDir, filepath.Join(home, ".gobuild", "plugins"))
	}
	if len(cfg.Projects) != 1 {
		t.Fatalf("len(Projects) = %d, want 1", len(cfg.Projects))
	}
	if cfg.Projects[0].Path != filepath.Join(home, "workspace", "demo") {
		t.Fatalf("Project path = %q, want %q", cfg.Projects[0].Path, filepath.Join(home, "workspace", "demo"))
	}
	if cfg.Projects[0].Env["APP_MODE"] != "debug" {
		t.Fatalf("Env[APP_MODE] = %q, want %q", cfg.Projects[0].Env["APP_MODE"], "debug")
	}
}

func TestLoadConfigInvalidYAMLReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.yaml")
	if err := os.WriteFile(path, []byte("projects: ["), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := LoadConfig(path); err == nil {
		t.Fatal("LoadConfig() expected error for invalid YAML, got nil")
	}
}
