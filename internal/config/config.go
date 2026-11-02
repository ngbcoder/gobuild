package config

// ProjectConfig represents a single project entry loaded from gobuild.yaml.
type ProjectConfig struct {
	Name      string
	Path      string
	Language  string
	BuildCmd  string
	TestCmd   string
	Env       map[string]string
	Overrides map[string]string
}

// GobuildConfig is the root configuration object processed by the app.
type GobuildConfig struct {
	Version    string
	PluginsDir string
	Projects   []ProjectConfig
}

// DefaultConfig returns a baseline configuration with safe defaults.
func DefaultConfig() *GobuildConfig {
	return &GobuildConfig{
		Version:    "1",
		PluginsDir: "~/.gobuild/plugins",
		Projects:   []ProjectConfig{},
	}
}
