package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTfRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	target := filepath.Join(root, "infra-as-code", "shared", "bootstrap")
	require.NoError(t, os.MkdirAll(target, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(target, "main.tf"), []byte("# empty\n"), 0o600))
	return root
}

func TestTf_InitDryRun(t *testing.T) {
	root := newTfRepo(t)
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "--dry-run", "tf", "init", "infra-as-code/shared/bootstrap"})
	require.NoError(t, cmd.Execute())
	out := stderr.String()
	assert.Contains(t, out, "→ [.] $ tofu -chdir=infra-as-code/shared/bootstrap init")
}

func TestTf_PlanWithExtraArgs(t *testing.T) {
	root := newTfRepo(t)
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "--dry-run", "tf", "plan", "infra-as-code/shared/bootstrap", "--", "-out=tfplan"})
	require.NoError(t, cmd.Execute())
	out := stderr.String()
	assert.Contains(t, out, "tofu -chdir=infra-as-code/shared/bootstrap plan -out=tfplan")
}

func TestTf_MissingTarget(t *testing.T) {
	root := t.TempDir()
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "tf", "init", "infra-as-code/shared/missing"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing")
}

func TestTf_InvalidAction(t *testing.T) {
	root := newTfRepo(t)
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stderr)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "tf", "nuke", "infra-as-code/shared/bootstrap"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}
