package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRepoWithStack(t *testing.T, name string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "docker-stacks", name)
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte("services: {}\n"), 0o600))
	return root
}

func TestStacks_Ls(t *testing.T) {
	root := newTestRepoWithStack(t, "postgres")
	var stdout bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--root", root, "stacks", "ls"})
	require.NoError(t, cmd.Execute())
	out := stdout.String()
	assert.Contains(t, out, "postgres")
	assert.Contains(t, out, "docker-stacks/postgres")
}

func TestStacks_UpDryRun(t *testing.T) {
	root := newTestRepoWithStack(t, "postgres")
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "--dry-run", "stacks", "up", "postgres"})
	require.NoError(t, cmd.Execute())
	out := stderr.String()
	assert.Contains(t, out, "→ [.] $ docker compose -f docker-stacks/postgres/compose.yaml up -d")
}

func TestStacks_DownDryRun(t *testing.T) {
	root := newTestRepoWithStack(t, "postgres")
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "--dry-run", "stacks", "down", "postgres"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, stderr.String(), "down")
	assert.NotContains(t, stderr.String(), "-v")
}

func TestStacks_NukeDryRun(t *testing.T) {
	root := newTestRepoWithStack(t, "postgres")
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "--dry-run", "stacks", "nuke", "postgres"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, stderr.String(), "down -v")
}

func TestStacks_NotFound(t *testing.T) {
	root := newTestRepoWithStack(t, "postgres")
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "stacks", "up", "missing"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing")
}
