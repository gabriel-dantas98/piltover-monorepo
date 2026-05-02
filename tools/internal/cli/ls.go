package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/discovery"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

func newLsCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "ls",
		Short: "List every discovered subproject",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projects, err := discovery.Discover(g.Root)
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return printLsJSON(cmd, projects)
			}
			return printLsTable(cmd, projects)
		},
	}
	c.Flags().Bool("json", false, "emit JSON suitable for piping into jq")
	return c
}

func printLsTable(cmd *cobra.Command, projects []*schema.Project) error {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tNAME\tKIND\tLANGUAGE\tTAGS")
	for _, p := range projects {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			p.Path, p.Name, p.Kind, p.Language, strings.Join(p.Tags, ","))
	}
	return w.Flush()
}

// lsEntry is the JSON shape emitted by ls --json.
type lsEntry struct {
	Path     string   `json:"path"`
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
}

func printLsJSON(cmd *cobra.Command, projects []*schema.Project) error {
	entries := make([]lsEntry, 0, len(projects))
	for _, p := range projects {
		entries = append(entries, lsEntry{
			Path:     p.Path,
			Name:     p.Name,
			Kind:     string(p.Kind),
			Language: string(p.Language),
			Tags:     p.Tags,
		})
	}
	b, err := json.Marshal(map[string]any{"include": entries})
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(b))
	return nil
}
