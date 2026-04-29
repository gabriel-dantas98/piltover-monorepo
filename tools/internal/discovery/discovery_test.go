package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover_FindsProjects(t *testing.T) {
	root := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(root, "apps/web"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "clis/foo"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "node_modules/junk"), 0o755))

	require.NoError(t, os.WriteFile(filepath.Join(root, "apps/web/project.yaml"),
		[]byte("name: web\nkind: app\nlanguage: ts\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "clis/foo/project.yaml"),
		[]byte("name: foo\nkind: cli\nlanguage: go\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "node_modules/junk/project.yaml"),
		[]byte("name: junk\nkind: cli\nlanguage: go\n"), 0o644))

	projects, err := Discover(root)
	require.NoError(t, err)
	require.Len(t, projects, 2, "node_modules must be skipped")

	names := []string{projects[0].Name, projects[1].Name}
	assert.ElementsMatch(t, []string{"web", "foo"}, names)
}

func TestDiscover_EmptyRepo(t *testing.T) {
	projects, err := Discover(t.TempDir())
	require.NoError(t, err)
	assert.Empty(t, projects)
}

func TestDiscover_SurfaceParseErrors(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "broken"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "broken/project.yaml"),
		[]byte("name: x\nkind: zoo\nlanguage: go\n"), 0o644))

	_, err := Discover(root)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "broken")
}
