package tf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTfDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	target := filepath.Join(root, "infra-as-code", "shared", "demo")
	require.NoError(t, os.MkdirAll(target, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(target, "main.tf"), []byte("# empty\n"), 0o600))
	return root
}

func TestResolve_Found(t *testing.T) {
	root := newTfDir(t)
	tgt, err := Resolve(root, "infra-as-code/shared/demo")
	require.NoError(t, err)
	assert.Equal(t, "infra-as-code/shared/demo", tgt.Path)
}

func TestResolve_MissingDir(t *testing.T) {
	_, err := Resolve(t.TempDir(), "infra-as-code/shared/missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing")
}

func TestResolve_NoTfFiles(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "x", "y"), 0o750))
	_, err := Resolve(root, "x/y")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no *.tf files")
}

func TestAction_Args_Init(t *testing.T) {
	tgt := Target{Path: "infra/x"}
	args := ActionInit.Args(tgt, nil)
	assert.Equal(t, []string{"-chdir=infra/x", "init"}, args)
}

func TestAction_Args_Plan_WithExtra(t *testing.T) {
	tgt := Target{Path: "infra/x"}
	args := ActionPlan.Args(tgt, []string{"-out=tfplan", "-var", "foo=bar"})
	assert.Equal(t, []string{"-chdir=infra/x", "plan", "-out=tfplan", "-var", "foo=bar"}, args)
}

func TestAction_Args_Apply(t *testing.T) {
	tgt := Target{Path: "infra/x"}
	args := ActionApply.Args(tgt, nil)
	assert.Equal(t, []string{"-chdir=infra/x", "apply"}, args)
}

func TestAction_Validate_Rejects_Unknown(t *testing.T) {
	require.True(t, ValidAction("init"))
	require.True(t, ValidAction("plan"))
	require.True(t, ValidAction("apply"))
	require.True(t, ValidAction("destroy"))
	require.True(t, ValidAction("validate"))
	require.True(t, ValidAction("fmt"))
	require.False(t, ValidAction("nuke"))
}
