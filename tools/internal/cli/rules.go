package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/rules"
)

const (
	kodyRulesDir = ".kody/rules"
	docsRulesDir = "apps/docs/content/rules"
)

func newRulesCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "rules",
		Short: "Manage Kody Custom Rules",
	}
	c.AddCommand(newRulesLsCmd(g))
	c.AddCommand(newRulesLintCmd(g))
	c.AddCommand(newRulesSyncDocsCmd(g))
	return c
}

func newRulesLsCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List every Kody rule",
		RunE: func(cmd *cobra.Command, _ []string) error {
			items, err := loadRules(g.Root)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SLUG\tSCOPE\tSEVERITY\tPATHS\tENABLED")
			for _, r := range items {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n",
					r.Slug, r.Scope, r.SeverityMin, strings.Join(r.Path, ","), r.Enabled)
			}
			return w.Flush()
		},
	}
}

func newRulesLintCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "Validate every Kody rule's frontmatter",
		RunE: func(cmd *cobra.Command, _ []string) error {
			items, err := loadRules(g.Root)
			if err != nil {
				return err
			}
			var problems []string
			for _, r := range items {
				if err := r.Validate(); err != nil {
					problems = append(problems, err.Error())
				}
			}
			if len(problems) > 0 {
				for _, p := range problems {
					cmd.PrintErrln(p)
				}
				return fmt.Errorf("%d rule(s) failed validation", len(problems))
			}
			cmd.Println("ok")
			return nil
		},
	}
}

func newRulesSyncDocsCmd(g *Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "sync-docs",
		Short: "Project .kody/rules/*.md to apps/docs/content/rules/*.mdx",
		RunE: func(cmd *cobra.Command, _ []string) error {
			items, err := loadRules(g.Root)
			if err != nil {
				return err
			}
			outDir := filepath.Join(g.Root, docsRulesDir)
			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}
			// Track desired filenames so we can delete stale ones.
			desired := map[string]bool{}
			for _, r := range items {
				name := r.Slug + ".mdx"
				desired[name] = true
				body := renderRuleMDX(r)
				path := filepath.Join(outDir, name)
				if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
					return err
				}
			}
			// Delete stale .mdx files.
			entries, err := os.ReadDir(outDir)
			if err != nil {
				return err
			}
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				if !strings.HasSuffix(e.Name(), ".mdx") {
					continue
				}
				if !desired[e.Name()] {
					if err := os.Remove(filepath.Join(outDir, e.Name())); err != nil {
						return err
					}
				}
			}
			fmt.Fprintln(cmd.OutOrStdout(), "synced", len(items), "rule(s) →", outDir)
			return nil
		},
	}
}

func renderRuleMDX(r *rules.Rule) string {
	desc := firstLine(r.Body)
	var b strings.Builder
	fmt.Fprintf(&b, "---\ntitle: %q\ndescription: %q\n---\n\n", r.Title, desc)
	fmt.Fprintln(&b, "import { Card } from 'fumadocs-ui/components/card';")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "<Card title=\"Kody Rule\" description=\"Lives at .kody/rules/%s.md\">\n", r.Slug)
	fmt.Fprintf(&b, "- **Scope:** %s\n", r.Scope)
	fmt.Fprintf(&b, "- **Severity ≥** %s\n", r.SeverityMin)
	fmt.Fprintf(&b, "- **Paths:** `%s`\n", strings.Join(r.Path, ", "))
	fmt.Fprintf(&b, "- **Languages:** %s\n", strings.Join(r.Languages, ", "))
	fmt.Fprintf(&b, "- **Buckets:** %s\n", strings.Join(r.Buckets, ", "))
	fmt.Fprintf(&b, "- **Enabled:** %v\n", r.Enabled)
	fmt.Fprintln(&b, "</Card>")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, r.Body)
	return b.String()
}

func firstLine(body string) string {
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		// Truncate to 160 chars max.
		if len(t) > 160 {
			t = t[:157] + "..."
		}
		return t
	}
	return ""
}

func loadRules(root string) ([]*rules.Rule, error) {
	dir := filepath.Join(root, kodyRulesDir)
	var out []*rules.Rule
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		data, err := os.ReadFile(path) // #nosec G304,G122 -- discovered path under .kody/rules; WalkDir is bounded to .kody/rules subtree
		if err != nil {
			return err
		}
		r, err := rules.Parse(data)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		r.Slug = strings.TrimSuffix(d.Name(), ".md")
		out = append(out, r)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out, nil
}
