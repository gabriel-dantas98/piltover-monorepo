package cli

import "github.com/spf13/cobra"

func newRunCmd(g *Globals, name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.PrintErrln(name, "is wired in Task 8")
			return nil
		},
	}
}
