package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/discovery"
)

func newLsCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List every discovered subproject",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projects, err := discovery.Discover(g.Root)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "PATH\tNAME\tKIND\tLANGUAGE\tTAGS")
			for _, p := range projects {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					p.Path, p.Name, p.Kind, p.Language, strings.Join(p.Tags, ","))
			}
			return w.Flush()
		},
	}
}
