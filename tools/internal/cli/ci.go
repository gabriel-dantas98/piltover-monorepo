package cli

import (
	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/discovery"
)

func newCiCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "ci",
		Short: "Run lint + test + build for every discovered project",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projects, err := discovery.Discover(g.Root)
			if err != nil {
				return err
			}
			for _, name := range []string{"lint", "test", "build"} {
				if err := runOnProjects(g, name, projects); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
