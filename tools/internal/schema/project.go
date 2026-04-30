// Package schema defines the project.yaml shape and validation rules
// used by the piltover engine to discover and operate on subprojects.
package schema

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Kind enumerates supported subproject kinds.
type Kind string

// Kind values.
const (
	KindApp         Kind = "app"
	KindPackage     Kind = "package"
	KindCLI         Kind = "cli"
	KindPlugin      Kind = "plugin"
	KindAction      Kind = "action"
	KindStack       Kind = "stack"
	KindInfraModule Kind = "infra-module"
)

// Language enumerates supported toolchains.
type Language string

// Language values.
const (
	LangGo     Language = "go"
	LangTS     Language = "ts"
	LangPython Language = "python"
	LangRust   Language = "rust"
	LangShell  Language = "shell"
	LangHCL    Language = "hcl"
	LangNone   Language = "none"
)

// ReleaseStrategy enumerates supported release pipelines.
type ReleaseStrategy string

// ReleaseStrategy values.
const (
	ReleaseChangesets    ReleaseStrategy = "changesets"
	ReleaseGoReleaser    ReleaseStrategy = "goreleaser"
	ReleasePyPI          ReleaseStrategy = "pypi-twine"
	ReleaseContainerOnly ReleaseStrategy = "container-only"
	ReleaseNone          ReleaseStrategy = "none"
)

// Commands captures per-project command overrides. Empty fields are filled
// from language defaults at runtime by the discovery layer.
type Commands struct {
	Lint  string `yaml:"lint"`
	Test  string `yaml:"test"`
	Build string `yaml:"build"`
}

// Release captures release strategy.
type Release struct {
	Strategy ReleaseStrategy `yaml:"strategy"`
}

// Project is the parsed shape of a project.yaml file.
type Project struct {
	Name     string   `yaml:"name"`
	Kind     Kind     `yaml:"kind"`
	Language Language `yaml:"language"`
	Tags     []string `yaml:"tags"`
	Commands Commands `yaml:"commands"`
	Release  Release  `yaml:"release"`

	// Path is the absolute or repo-relative path to the project directory
	// containing the project.yaml. Set by LoadFromDir; not present in YAML.
	Path string `yaml:"-"`
}

var validKinds = map[Kind]bool{
	KindApp: true, KindPackage: true, KindCLI: true, KindPlugin: true,
	KindAction: true, KindStack: true, KindInfraModule: true,
}

var validLanguages = map[Language]bool{
	LangGo: true, LangTS: true, LangPython: true, LangRust: true,
	LangShell: true, LangHCL: true, LangNone: true,
}

// ParseProject parses YAML bytes into a Project, validating required fields.
func ParseProject(data []byte) (*Project, error) {
	var p Project
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	if p.Name == "" {
		return nil, fmt.Errorf("project.yaml: name is required")
	}
	if !validKinds[p.Kind] {
		return nil, fmt.Errorf("project.yaml: unknown kind %q", p.Kind)
	}
	if !validLanguages[p.Language] {
		return nil, fmt.Errorf("project.yaml: unknown language %q", p.Language)
	}
	return &p, nil
}

// LoadFromDir reads dir/project.yaml and returns a parsed Project with Path set
// to the directory.
func LoadFromDir(dir string) (*Project, error) {
	// #nosec G304 -- project.yaml path comes from in-repo discovery, not external input
	data, err := os.ReadFile(filepath.Join(dir, "project.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read project.yaml: %w", err)
	}
	p, err := ParseProject(data)
	if err != nil {
		return nil, err
	}
	p.Path = dir
	return p, nil
}
