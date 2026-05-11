package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runNew(t *testing.T, root string, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(append([]string{"--root", root, "new"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func TestNew_Cli(t *testing.T) {
	root := t.TempDir()
	_, err := runNew(t, root, "cli", "demo")
	require.NoError(t, err)

	projPath := filepath.Join(root, "clis", "demo", "project.yaml")
	data, err := os.ReadFile(projPath) //nolint:gosec // test only reads temp dir
	require.NoError(t, err)
	assert.Contains(t, string(data), "name: demo")
	assert.Contains(t, string(data), "kind: cli")

	mainPath := filepath.Join(root, "clis", "demo", "cmd", "demo", "main.go")
	data, err = os.ReadFile(mainPath) //nolint:gosec // test only reads temp dir
	require.NoError(t, err)
	assert.Contains(t, string(data), "hello from demo")
}

func TestNew_Package(t *testing.T) {
	root := t.TempDir()
	_, err := runNew(t, root, "package", "lib")
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(root, "packages", "lib", "package.json")) //nolint:gosec // test only reads temp dir
	require.NoError(t, err)
	assert.Contains(t, string(data), "@piltover/lib")
}

func TestNew_App(t *testing.T) {
	root := t.TempDir()
	_, err := runNew(t, root, "app", "web")
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(root, "apps", "web", "app", "page.tsx"))
	require.NoError(t, err)
}

func TestNew_Stack(t *testing.T) {
	root := t.TempDir()
	_, err := runNew(t, root, "stack", "redis")
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(root, "docker-stacks", "redis", "compose.yaml")) //nolint:gosec // test only reads temp dir
	require.NoError(t, err)
	assert.Contains(t, string(data), "piltover-redis")
}

func TestNew_RefusesExisting(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "clis", "exists"), 0o750))
	_, err := runNew(t, root, "cli", "exists")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exists")
}

func TestNew_RefusesUnknownKind(t *testing.T) {
	_, err := runNew(t, t.TempDir(), "plugin", "x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "plugin")
}
