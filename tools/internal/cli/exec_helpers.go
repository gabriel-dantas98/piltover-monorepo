package cli

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/configs"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/runner"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

// LanguageDefaults holds default commands for one language.
type LanguageDefaults struct {
	Lint  string `yaml:"lint"`
	Test  string `yaml:"test"`
	Build string `yaml:"build"`
}

// LoadDefaults parses the embedded defaults.yaml.
func LoadDefaults() (map[schema.Language]LanguageDefaults, error) {
	out := map[schema.Language]LanguageDefaults{}
	raw := map[string]LanguageDefaults{}
	if err := yaml.Unmarshal(configs.DefaultsYAML, &raw); err != nil {
		return nil, fmt.Errorf("parse defaults.yaml: %w", err)
	}
	for k, v := range raw {
		out[schema.Language(k)] = v
	}
	return out, nil
}

func resolveCommand(p *schema.Project, name string, defaults map[schema.Language]LanguageDefaults) string {
	switch name {
	case "lint":
		if p.Commands.Lint != "" {
			return p.Commands.Lint
		}
		return defaults[p.Language].Lint
	case "test":
		if p.Commands.Test != "" {
			return p.Commands.Test
		}
		return defaults[p.Language].Test
	case "build":
		if p.Commands.Build != "" {
			return p.Commands.Build
		}
		return defaults[p.Language].Build
	}
	return ""
}

func runOnProjects(g *Globals, name string, projects []*schema.Project) error {
	defaults, err := LoadDefaults()
	if err != nil {
		return err
	}
	r := runner.New(runner.Options{Verbose: g.Verbose, Quiet: g.Quiet, DryRun: g.DryRun})
	for _, p := range projects {
		cmdline := resolveCommand(p, name, defaults)
		if cmdline == "" {
			continue
		}
		// naïve split: defaults are simple commands; complex pipelines should be
		// declared via a wrapper script invoked from project.yaml.
		parts := strings.Fields(cmdline)
		if err := r.Run(runner.Cmd{Cwd: p.Path, Name: parts[0], Args: parts[1:]}); err != nil {
			return err
		}
	}
	return nil
}
