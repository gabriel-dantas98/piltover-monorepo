// Package runner runs external commands with mandatory pre-execution logging.
//
// Every Cmd dispatched through Runner is announced on Stderr in the form:
//
//	→ [<cwd>] $ <name> <args...>
//
// The logging is unconditional unless Options.Quiet is set. With Options.DryRun
// the command is logged but not executed (useful for CI debugging). With
// Options.Verbose, additional environment variables are printed on a second line.
package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Cmd describes a command the engine wants to execute.
type Cmd struct {
	Cwd  string   // project-relative path used in the log line
	Name string   // executable
	Args []string // arguments
	Env  []string // additional environment, KEY=VALUE
}

// Options configures a Runner.
type Options struct {
	Stderr  io.Writer
	Stdout  io.Writer
	Quiet   bool
	Verbose bool
	DryRun  bool
}

// Runner is the entry point for spawning external commands.
type Runner struct {
	opts Options
}

// New constructs a Runner. If Stderr/Stdout are nil they default to os.Stderr/os.Stdout.
func New(opts Options) *Runner {
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	return &Runner{opts: opts}
}

// Run logs the command and, unless DryRun is set, executes it synchronously.
func (r *Runner) Run(c Cmd) error {
	r.logCmd(c)
	if r.opts.DryRun {
		return nil
	}
	// #nosec G204 -- runner deliberately spawns user-supplied commands; sanitisation is the caller's responsibility
	cmd := exec.Command(c.Name, c.Args...)
	if c.Cwd != "" {
		cmd.Dir = c.Cwd
	}
	cmd.Stdout = r.opts.Stdout
	cmd.Stderr = r.opts.Stderr
	if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed (%s): %w", c.Name, err)
	}
	return nil
}

func (r *Runner) logCmd(c Cmd) {
	if r.opts.Quiet {
		return
	}
	cwd := c.Cwd
	if cwd == "" {
		cwd = "."
	}
	parts := append([]string{c.Name}, c.Args...)
	fmt.Fprintf(r.opts.Stderr, "→ [%s] $ %s\n", cwd, strings.Join(parts, " "))
	if r.opts.Verbose && len(c.Env) > 0 {
		fmt.Fprintf(r.opts.Stderr, "    env: %s\n", strings.Join(c.Env, " "))
	}
}
