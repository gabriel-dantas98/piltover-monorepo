package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/discovery"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

func newRunCmd(g *Globals, name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name + " [paths...]",
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := discovery.Discover(g.Root)
			if err != nil {
				return err
			}
			projects = filterByPaths(projects, args)
			return runOnProjects(g, name, projects)
		},
	}
}

func filterByPaths(projects []*schema.Project, paths []string) []*schema.Project {
	if len(paths) == 0 {
		return projects
	}
	out := make([]*schema.Project, 0, len(projects))
	for _, p := range projects {
		for _, want := range paths {
			if p.Path == strings.TrimRight(want, "/") {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
