// Package discovery walks a repository root and returns every subproject
// declared by a project.yaml file.
package discovery

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

// skipDirs are folder names whose subtree we never descend into.
var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	".next":        true,
	".turbo":       true,
	"dist":         true,
	"build":        true,
	".terraform":   true,
	".venv":        true,
	"__pycache__":  true,
	"bin":          true,
}

// Discover returns every project found under root, sorted by relative path.
func Discover(root string) ([]*schema.Project, error) {
	var projects []*schema.Project

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != "project.yaml" {
			return nil
		}
		dir := filepath.Dir(path)
		p, err := schema.LoadFromDir(dir)
		if err != nil {
			rel, _ := filepath.Rel(root, dir)
			return fmt.Errorf("invalid project at %s: %w", rel, err)
		}
		rel, _ := filepath.Rel(root, dir)
		p.Path = rel
		projects = append(projects, p)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(projects, func(i, j int) bool {
		return strings.Compare(projects[i].Path, projects[j].Path) < 0
	})
	return projects, nil
}
