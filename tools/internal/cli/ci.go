package cli

import "github.com/spf13/cobra"

func newCiCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "ci",
		Short: "Run lint + test + build with JSON-friendly output",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.PrintErrln("ci is wired in Task 8")
			return nil
		},
	}
}
