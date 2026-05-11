package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/runner"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/stacks"
)

func newStacksCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "stacks",
		Short: "Manage local-only docker-compose stacks under docker-stacks/",
	}
	c.AddCommand(newStacksLsCmd(g))
	c.AddCommand(newStacksUpCmd(g))
	c.AddCommand(newStacksDownCmd(g, false))
	c.AddCommand(newStacksNukeCmd(g))
	return c
}

func newStacksLsCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List discovered stacks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			items, err := stacks.List(g.Root)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPATH")
			for _, s := range items {
				fmt.Fprintf(w, "%s\t%s\n", s.Name, s.Path)
			}
			return w.Flush()
		},
	}
}

func newStacksUpCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "up <name>",
		Short: "Start a stack via `docker compose up -d`",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := stacks.Resolve(g.Root, args[0])
			if err != nil {
				return err
			}
			r := runner.New(runner.Options{
				Quiet:   g.Quiet,
				Verbose: g.Verbose,
				DryRun:  g.DryRun,
				Stderr:  cmd.ErrOrStderr(),
				Stdout:  cmd.OutOrStdout(),
			})
			return r.Run(runner.Cmd{
				Name: "docker",
				Args: []string{"compose", "-f", s.ComposeFile, "up", "-d"},
			})
		},
	}
}

func newStacksDownCmd(g *Globals, withVolumes bool) *cobra.Command {
	return &cobra.Command{
		Use:   "down <name>",
		Short: "Stop a stack via `docker compose down`",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := stacks.Resolve(g.Root, args[0])
			if err != nil {
				return err
			}
			downArgs := []string{"compose", "-f", s.ComposeFile, "down"}
			if withVolumes {
				downArgs = append(downArgs, "-v")
			}
			r := runner.New(runner.Options{
				Quiet:   g.Quiet,
				Verbose: g.Verbose,
				DryRun:  g.DryRun,
				Stderr:  cmd.ErrOrStderr(),
				Stdout:  cmd.OutOrStdout(),
			})
			return r.Run(runner.Cmd{Name: "docker", Args: downArgs})
		},
	}
}

func newStacksNukeCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "nuke <name>",
		Short: "Stop a stack and remove its volumes (`docker compose down -v`)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := stacks.Resolve(g.Root, args[0])
			if err != nil {
				return err
			}
			r := runner.New(runner.Options{
				Quiet:   g.Quiet,
				Verbose: g.Verbose,
				DryRun:  g.DryRun,
				Stderr:  cmd.ErrOrStderr(),
				Stdout:  cmd.OutOrStdout(),
			})
			return r.Run(runner.Cmd{
				Name: "docker",
				Args: []string{"compose", "-f", s.ComposeFile, "down", "-v"},
			})
		},
	}
}
