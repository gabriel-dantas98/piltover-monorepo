package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/discovery"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

type matrixEntry struct {
	Path     string   `json:"path"`
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
}

func newAffectedCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "affected",
		Short: "Emit JSON matrix of projects touched since --base",
		RunE: func(cmd *cobra.Command, _ []string) error {
			base, _ := cmd.Flags().GetString("base")
			projects, err := discovery.Discover(g.Root)
			if err != nil {
				return err
			}
			changed, err := changedFiles(g, base)
			if err != nil {
				return err
			}
			affected := projectsForChangedFiles(projects, changed)

			entries := make([]matrixEntry, 0, len(affected))
			for _, p := range affected {
				entries = append(entries, matrixEntry{
					Path:     p.Path,
					Name:     p.Name,
					Kind:     string(p.Kind),
					Language: string(p.Language),
					Tags:     p.Tags,
				})
			}
			payload := map[string]any{"include": entries}
			b, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
	c.Flags().String("base", "origin/main", "git ref to diff against")
	return c
}

func projectsForChangedFiles(projects []*schema.Project, changed []string) []*schema.Project {
	out := make([]*schema.Project, 0)
	seen := map[string]bool{}
	for _, p := range projects {
		prefix := p.Path + "/"
		for _, f := range changed {
			if strings.HasPrefix(f, prefix) || f == p.Path {
				if !seen[p.Path] {
					out = append(out, p)
					seen[p.Path] = true
				}
				break
			}
		}
	}
	return out
}

func changedFiles(g *Globals, base string) ([]string, error) {
	out, err := gitDiffNames(g, base)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}
	return lines, nil
}

func gitDiffNames(g *Globals, base string) (string, error) {
	if !g.Quiet {
		fmt.Fprintf(stderrSink(), "→ [.] $ git diff --name-only %s...HEAD\n", base)
	}
	if g.DryRun {
		return "", nil
	}
	cmd := exec.Command("git", "diff", "--name-only", base+"...HEAD")
	cmd.Dir = g.Root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff failed: %s", stderr.String())
	}
	return stdout.String(), nil
}
