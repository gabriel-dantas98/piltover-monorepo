// Package tf provides helpers to resolve OpenTofu root-module paths
// and build the argv passed to `tofu`.
package tf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Target is a resolved OpenTofu root module.
type Target struct {
	Path string // repo-relative path to the target directory
}

// Action enumerates supported tofu actions.
type Action string

// Supported OpenTofu subcommands.
const (
	ActionInit     Action = "init"
	ActionValidate Action = "validate"
	ActionPlan     Action = "plan"
	ActionApply    Action = "apply"
	ActionDestroy  Action = "destroy"
	ActionFmt      Action = "fmt"
)

var validActions = map[string]bool{
	string(ActionInit):     true,
	string(ActionValidate): true,
	string(ActionPlan):     true,
	string(ActionApply):    true,
	string(ActionDestroy):  true,
	string(ActionFmt):      true,
}

// ValidAction reports whether the given string is a supported Action.
func ValidAction(s string) bool { return validActions[s] }

// Resolve verifies that <root>/<target> exists and contains at least one *.tf file.
// Returns a Target whose Path is the repo-relative target.
func Resolve(root, target string) (Target, error) {
	clean := filepath.Clean(target)
	abs := filepath.Join(root, clean)
	info, err := os.Stat(abs)
	if err != nil {
		return Target{}, fmt.Errorf("target %q: %w", target, err)
	}
	if !info.IsDir() {
		return Target{}, fmt.Errorf("target %q is not a directory", target)
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return Target{}, fmt.Errorf("read %q: %w", target, err)
	}
	hasTf := false
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".tf") {
			hasTf = true
			break
		}
	}
	if !hasTf {
		return Target{}, fmt.Errorf("target %q contains no *.tf files", target)
	}
	return Target{Path: clean}, nil
}

// Args returns the argv to pass to `tofu` for the given Action against the Target,
// appending any caller-supplied extra args after the subcommand.
func (a Action) Args(target Target, extra []string) []string {
	out := []string{"-chdir=" + target.Path, string(a)}
	out = append(out, extra...)
	return out
}
