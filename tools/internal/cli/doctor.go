package cli

import "github.com/spf13/cobra"

func newDoctorCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "doctor",
		Short: "Verify required toolchains",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.PrintErrln("doctor is wired in Task 10")
			return nil
		},
	}
	c.Flags().Bool("json", false, "emit JSON")
	return c
}
