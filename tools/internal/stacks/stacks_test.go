package stacks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_EmptyRepo(t *testing.T) {
	stacks, err := List(t.TempDir())
	require.NoError(t, err)
	assert.Empty(t, stacks)
}

func TestList_FindsStacks(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"postgres", "localstack"} {
		dir := filepath.Join(root, "docker-stacks", name)
		require.NoError(t, os.MkdirAll(dir, 0o750))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte("services: {}\n"), 0o600))
	}
	// A noise directory without compose.yaml should be skipped.
	require.NoError(t, os.MkdirAll(filepath.Join(root, "docker-stacks", "noise"), 0o750))

	stacks, err := List(root)
	require.NoError(t, err)
	require.Len(t, stacks, 2)
	assert.Equal(t, "localstack", stacks[0].Name)
	assert.Equal(t, "postgres", stacks[1].Name)
	assert.Equal(t, filepath.Join("docker-stacks", "postgres"), stacks[1].Path)
	assert.Equal(t, filepath.Join("docker-stacks", "postgres", "compose.yaml"), stacks[1].ComposeFile)
}

func TestResolve_Found(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "docker-stacks", "postgres")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte("services: {}\n"), 0o600))

	s, err := Resolve(root, "postgres")
	require.NoError(t, err)
	assert.Equal(t, "postgres", s.Name)
}

func TestResolve_NotFound(t *testing.T) {
	_, err := Resolve(t.TempDir(), "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing")
}
