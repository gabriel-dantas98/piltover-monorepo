// Package stacks discovers and represents local docker-compose stacks
// located under docker-stacks/<name>/compose.yaml at the repo root.
package stacks

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Stack describes one discovered local docker-compose stack.
type Stack struct {
	Name        string // basename of the directory under docker-stacks/
	Path        string // repo-relative path to the stack directory
	ComposeFile string // repo-relative path to the compose.yaml file
}

// List returns every Stack found under <root>/docker-stacks/, sorted by Name.
// A directory without a compose.yaml is silently skipped. If docker-stacks/
// does not exist, returns (nil, nil).
func List(root string) ([]Stack, error) {
	dir := filepath.Join(root, "docker-stacks")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read docker-stacks/: %w", err)
	}
	out := make([]Stack, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		composeAbs := filepath.Join(dir, e.Name(), "compose.yaml")
		if _, err := os.Stat(composeAbs); err != nil {
			continue
		}
		out = append(out, Stack{
			Name:        e.Name(),
			Path:        filepath.Join("docker-stacks", e.Name()),
			ComposeFile: filepath.Join("docker-stacks", e.Name(), "compose.yaml"),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Resolve returns the Stack with the given name, or an error if not found.
func Resolve(root, name string) (Stack, error) {
	all, err := List(root)
	if err != nil {
		return Stack{}, err
	}
	for _, s := range all {
		if s.Name == name {
			return s, nil
		}
	}
	return Stack{}, fmt.Errorf("stack %q not found under docker-stacks/", name)
}
