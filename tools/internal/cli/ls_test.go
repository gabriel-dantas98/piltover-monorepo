package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLs_PrintsDiscoveredProjects(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "clis/foo"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "clis/foo/project.yaml"),
		[]byte("name: foo\nkind: cli\nlanguage: go\n"), 0o644))

	var stdout bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--root", root, "ls"})
	require.NoError(t, cmd.Execute())

	out := stdout.String()
	assert.Contains(t, out, "clis/foo")
	assert.Contains(t, out, "cli")
	assert.Contains(t, out, "go")
	assert.Contains(t, out, "foo")
}
