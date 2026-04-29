package cli

import "github.com/spf13/cobra"

func newAffectedCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "affected",
		Short: "Emit JSON matrix of projects touched since --base",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.PrintErrln("affected is wired in Task 9")
			return nil
		},
	}
	c.Flags().String("base", "origin/main", "git ref to diff against")
	return c
}
