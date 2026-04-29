package runner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_LogsCommandPrefix(t *testing.T) {
	var stderr bytes.Buffer
	r := New(Options{Stderr: &stderr})
	err := r.Run(Cmd{Cwd: "apps/web", Name: "echo", Args: []string{"hello"}})
	require.NoError(t, err)
	assert.Contains(t, stderr.String(), "→ [apps/web] $ echo hello")
}

func TestRun_QuietSuppressesPrefix(t *testing.T) {
	var stderr bytes.Buffer
	r := New(Options{Stderr: &stderr, Quiet: true})
	err := r.Run(Cmd{Cwd: ".", Name: "echo", Args: []string{"x"}})
	require.NoError(t, err)
	assert.NotContains(t, stderr.String(), "→ [")
}

func TestRun_DryRunSkipsExecution(t *testing.T) {
	var stderr bytes.Buffer
	r := New(Options{Stderr: &stderr, DryRun: true})
	err := r.Run(Cmd{Cwd: "x", Name: "false"}) // would exit non-zero if executed
	require.NoError(t, err)
	assert.Contains(t, stderr.String(), "→ [x] $ false")
}

func TestRun_VerboseLogsEnv(t *testing.T) {
	var stderr bytes.Buffer
	r := New(Options{Stderr: &stderr, Verbose: true})
	err := r.Run(Cmd{Cwd: ".", Name: "echo", Args: []string{"y"}, Env: []string{"FOO=bar"}})
	require.NoError(t, err)
	out := stderr.String()
	assert.Contains(t, out, "→ [.] $ echo y")
	assert.Contains(t, out, "FOO=bar")
}

func TestRun_PropagatesNonZeroExit(t *testing.T) {
	var stderr bytes.Buffer
	r := New(Options{Stderr: &stderr})
	err := r.Run(Cmd{Cwd: ".", Name: "sh", Args: []string{"-c", "exit 7"}})
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "exit"))
}
