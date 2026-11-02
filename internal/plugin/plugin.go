package plugin

// Manifest describes a language plugin advertised by a plugin.yaml file.
type Manifest struct {
	Name          string
	Language      string
	DetectFiles   []string
	FileExtensions []string
	Commands      map[string]string
}

// Plugin represents a discovered plugin installation.
type Plugin struct {
	Name     string
	Path     string
	Manifest Manifest
}

// DefaultPluginDir is the default filesystem location for installed plugins.
const DefaultPluginDir = "~/.gobuild/plugins"
