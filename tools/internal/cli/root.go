// Package cli defines the Cobra command tree for the piltover engine.
package cli

import (
	"github.com/spf13/cobra"
)

// Globals holds flags shared by every subcommand.
type Globals struct {
	Root    string
	Verbose bool
	Quiet   bool
	DryRun  bool
}

// NewRootCmd builds the root cobra.Command. Each subcommand is attached here.
func NewRootCmd() *cobra.Command {
	g := &Globals{}
	root := &cobra.Command{
		Use:           "piltover",
		Short:         "Thin orchestrator for the Piltover monorepo",
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.PersistentFlags().StringVar(&g.Root, "root", ".", "repository root")
	root.PersistentFlags().BoolVarP(&g.Verbose, "verbose", "v", false, "verbose logging")
	root.PersistentFlags().BoolVar(&g.Quiet, "quiet", false, "suppress command logs")
	root.PersistentFlags().BoolVar(&g.DryRun, "dry-run", false, "print commands without executing")

	root.AddCommand(newLsCmd(g))
	root.AddCommand(newDoctorCmd(g))
	root.AddCommand(newAffectedCmd(g))
	root.AddCommand(newRunCmd(g, "lint", "Run lint for affected (or specified) projects"))
	root.AddCommand(newRunCmd(g, "test", "Run tests"))
	root.AddCommand(newRunCmd(g, "build", "Run build"))
	root.AddCommand(newCiCmd(g))
	root.AddCommand(newScaffoldCmd(g))
	root.AddCommand(newStubCmd("tf", "Wrap OpenTofu (Plan 4)"))
	root.AddCommand(newStacksCmd(g))
	root.AddCommand(newRulesCmd(g))

	return root
}

func newStubCmd(name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short + " — not yet implemented",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.PrintErrln(name, "is not implemented in v0; see plan 4/5")
			return nil
		},
	}
}
