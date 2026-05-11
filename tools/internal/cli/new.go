package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/scaffold"
)

// kindRoute maps a user-supplied kind to (template subdir, top-level folder).
type kindRoute struct {
	Template string
	Parent   string
}

var newRoutes = map[string]kindRoute{
	"cli":     {Template: "cli-go", Parent: "clis"},
	"package": {Template: "package-ts", Parent: "packages"},
	"app":     {Template: "app-ts", Parent: "apps"},
	"stack":   {Template: "stack", Parent: "docker-stacks"},
}

func newScaffoldCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "new <kind> <name>",
		Short: "Scaffold a new subproject (cli | package | app | stack)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := args[0]
			name := args[1]

			route, ok := newRoutes[kind]
			if !ok {
				return fmt.Errorf("kind %q not implemented for v1; supported: cli, package, app, stack", kind)
			}

			kebab := toKebab(name)
			dest := filepath.Join(g.Root, route.Parent, kebab)
			if _, err := os.Stat(dest); err == nil {
				return fmt.Errorf("destination already exists: %s", dest)
			}

			vars := scaffold.Vars{Name: name, KebabName: kebab}
			if err := scaffold.Render(scaffold.Templates(), route.Template, dest, vars); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Scaffolded %s/%s at %s\n", route.Parent, kebab, dest)
			fmt.Fprintln(cmd.OutOrStdout(), "Next steps:")
			fmt.Fprintln(cmd.OutOrStdout(), "  piltover ls")
			fmt.Fprintf(cmd.OutOrStdout(), "  piltover lint %s/%s\n", route.Parent, kebab)
			return nil
		},
	}
}

// toKebab lowercases and replaces underscores/spaces with hyphens.
func toKebab(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
