package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/runner"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/tf"
)

func newTfCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "tf <action> <target> [-- extra args]",
		Short: "Run an OpenTofu action against a root module (init|plan|apply|destroy|validate|fmt)",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			actionStr := args[0]
			targetArg := args[1]
			extra := args[2:]

			if !tf.ValidAction(actionStr) {
				return fmt.Errorf("unsupported action %q (valid: init|plan|apply|destroy|validate|fmt)", actionStr)
			}
			action := tf.Action(actionStr)

			target, err := tf.Resolve(g.Root, targetArg)
			if err != nil {
				return err
			}

			r := runner.New(runner.Options{
				Quiet:   g.Quiet,
				Verbose: g.Verbose,
				DryRun:  g.DryRun,
				Stdout:  cmd.OutOrStdout(),
				Stderr:  cmd.ErrOrStderr(),
			})
			return r.Run(runner.Cmd{
				Name: "tofu",
				Args: action.Args(target, extra),
			})
		},
	}
	return c
}
