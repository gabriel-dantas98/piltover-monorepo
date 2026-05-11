package scaffold

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_BasicTemplate(t *testing.T) {
	tmplFS := fstest.MapFS{
		"cli-go/project.yaml.tmpl": &fstest.MapFile{
			Data: []byte("name: {{.KebabName}}\nkind: cli\nlanguage: go\n"),
		},
		"cli-go/cmd/{{.Name}}/main.go.tmpl": &fstest.MapFile{
			Data: []byte("package main\n\nfunc main() { println(\"hi from {{.KebabName}}\") }\n"),
		},
	}
	dest := t.TempDir()
	vars := Vars{Name: "foo", KebabName: "foo"}
	require.NoError(t, Render(tmplFS, "cli-go", dest, vars))

	// project.yaml rendered
	data, err := os.ReadFile(filepath.Join(dest, "project.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "name: foo")

	// nested path with {{.Name}} resolved, .tmpl suffix removed
	data, err = os.ReadFile(filepath.Join(dest, "cmd", "foo", "main.go"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "hi from foo")
}

func TestRender_RefusesNonexistentKind(t *testing.T) {
	tmplFS := fstest.MapFS{
		"cli-go/project.yaml.tmpl": &fstest.MapFile{Data: []byte("hi")},
	}
	err := Render(tmplFS, "nonexistent", t.TempDir(), Vars{Name: "x", KebabName: "x"})
	require.Error(t, err)
}

func TestRender_KeepsNonTmplFilesVerbatim(t *testing.T) {
	tmplFS := fstest.MapFS{
		"stack/README.md": &fstest.MapFile{Data: []byte("# Static\n")},
	}
	dest := t.TempDir()
	require.NoError(t, Render(tmplFS, "stack", dest, Vars{Name: "x", KebabName: "x"}))
	data, err := os.ReadFile(filepath.Join(dest, "README.md"))
	require.NoError(t, err)
	assert.Equal(t, "# Static\n", string(data))
}

func ensureExists(t *testing.T, root, rel string) {
	t.Helper()
	_, err := fs.Stat(os.DirFS(root), rel)
	require.NoError(t, err)
}
