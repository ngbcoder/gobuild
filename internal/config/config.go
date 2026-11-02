package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents a single project entry loaded from gobuild.yaml.
type ProjectConfig struct {
	Name      string            `yaml:"name"`
	Path      string            `yaml:"path"`
	Language  string            `yaml:"language"`
	BuildCmd  string            `yaml:"build_cmd"`
	TestCmd   string            `yaml:"test_cmd"`
	Env       map[string]string `yaml:"env"`
	Overrides map[string]string `yaml:"overrides"`
}

// GobuildConfig is the root configuration object processed by the app.
type GobuildConfig struct {
	Version    string          `yaml:"version"`
	PluginsDir string          `yaml:"plugins_dir"`
	Projects   []ProjectConfig `yaml:"projects"`
}

// DefaultConfig returns a baseline configuration with safe defaults.
func DefaultConfig() *GobuildConfig {
	return &GobuildConfig{
		Version:    "1",
		PluginsDir: "~/.gobuild/plugins",
		Projects:   []ProjectConfig{},
	}
}

// LoadConfig reads a gobuild.yaml file and returns a populated configuration or
// the default configuration if the file is missing.
func LoadConfig(path string) (*GobuildConfig, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	var raw GobuildConfig
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return nil, err
	}

	cfg.Version = raw.Version
	if cfg.Version == "" {
		cfg.Version = "1"
	}

	if raw.PluginsDir != "" {
		cfg.PluginsDir = expandPath(raw.PluginsDir)
	} else {
		cfg.PluginsDir = "~/.gobuild/plugins"
	}

	cfg.Projects = nil
	for _, project := range raw.Projects {
		p := project
		if p.Path != "" {
			p.Path = expandPath(p.Path)
		}
		if p.Env == nil {
			p.Env = map[string]string{}
		}
		if p.Overrides == nil {
			p.Overrides = map[string]string{}
		}
		cfg.Projects = append(cfg.Projects, p)
	}

	return cfg, nil
}

func expandPath(path string) string {
	if path == "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}

	switch {
	case path == "~":
		return home
	case strings.HasPrefix(path, "~/"):
		return filepath.Join(home, filepath.Clean(strings.TrimPrefix(path, "~/")))
	case strings.HasPrefix(path, "~\\"):
		return filepath.Join(home, filepath.Clean(strings.TrimPrefix(path, "~\\")))
	default:
		return path
	}
}
