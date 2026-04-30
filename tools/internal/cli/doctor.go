package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// Check defines one toolchain probe.
type Check struct {
	Name string
	Cmd  string   // executable
	Args []string // arguments
}

// CheckResult is the outcome of a single Check.
type CheckResult struct {
	Name   string `json:"name"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

var defaultChecks = []Check{
	{Name: "go", Cmd: "go", Args: []string{"version"}},
	{Name: "node", Cmd: "node", Args: []string{"--version"}},
	{Name: "bun", Cmd: "bun", Args: []string{"--version"}},
	{Name: "python", Cmd: "python3", Args: []string{"--version"}},
	{Name: "uv", Cmd: "uv", Args: []string{"--version"}},
	{Name: "tofu", Cmd: "tofu", Args: []string{"--version"}},
	{Name: "docker", Cmd: "docker", Args: []string{"--version"}},
	{Name: "git", Cmd: "git", Args: []string{"--version"}},
	{Name: "lefthook", Cmd: "lefthook", Args: []string{"version"}},
	{Name: "aws", Cmd: "aws", Args: []string{"--version"}},
	{Name: "gh", Cmd: "gh", Args: []string{"--version"}},
}

func newDoctorCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "doctor",
		Short: "Verify required toolchains",
		RunE: func(cmd *cobra.Command, _ []string) error {
			asJSON, _ := cmd.Flags().GetBool("json")
			results := runChecks(defaultChecks, g.Verbose)
			if asJSON {
				cmd.Println(renderJSON(results))
				return nil
			}
			cmd.Println(renderText(results))
			return nil
		},
	}
	c.Flags().Bool("json", false, "emit JSON")
	return c
}

func runChecks(checks []Check, verbose bool) []CheckResult {
	out := make([]CheckResult, 0, len(checks))
	for _, ch := range checks {
		if verbose {
			parts := append([]string{ch.Cmd}, ch.Args...)
			fmt.Fprintf(stderrSink(), "→ [.] $ %s\n", strings.Join(parts, " "))
		}
		// #nosec G204 -- doctor probes fixed well-known toolchain commands; not user-supplied input
		c := exec.Command(ch.Cmd, ch.Args...)
		bytesOut, err := c.CombinedOutput()
		detail := strings.TrimSpace(string(bytesOut))
		if err != nil {
			detail = err.Error()
		}
		out = append(out, CheckResult{
			Name:   ch.Name,
			OK:     err == nil,
			Detail: detail,
		})
	}
	return out
}

func renderJSON(results []CheckResult) string {
	b, _ := json.Marshal(results)
	return string(b)
}

func renderText(results []CheckResult) string {
	var b strings.Builder
	for _, r := range results {
		mark := "✓"
		if !r.OK {
			mark = "✗"
		}
		fmt.Fprintf(&b, "  %s %-10s %s\n", mark, r.Name, r.Detail)
	}
	return b.String()
}
