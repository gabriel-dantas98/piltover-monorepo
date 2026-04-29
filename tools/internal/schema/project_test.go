package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProject_Minimal(t *testing.T) {
	yaml := `
name: brag-cli
kind: cli
language: go
`
	p, err := ParseProject([]byte(yaml))
	require.NoError(t, err)
	assert.Equal(t, "brag-cli", p.Name)
	assert.Equal(t, KindCLI, p.Kind)
	assert.Equal(t, LangGo, p.Language)
}

func TestParseProject_FullWithCommands(t *testing.T) {
	yaml := `
name: piltover-docs
kind: app
language: ts
tags: [docs, public]
commands:
  lint: bun run lint
  test: bun run test
  build: bun run build
release:
  strategy: none
`
	p, err := ParseProject([]byte(yaml))
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"docs", "public"}, p.Tags)
	assert.Equal(t, "bun run lint", p.Commands.Lint)
	assert.Equal(t, ReleaseNone, p.Release.Strategy)
}

func TestParseProject_RejectsUnknownKind(t *testing.T) {
	_, err := ParseProject([]byte(`
name: x
kind: zoo
language: go
`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "kind")
}

func TestParseProject_RejectsUnknownLanguage(t *testing.T) {
	_, err := ParseProject([]byte(`
name: x
kind: cli
language: cobol
`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "language")
}

func TestLoadFromDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "project.yaml")
	require.NoError(t, os.WriteFile(path, []byte(`name: x
kind: package
language: ts
`), 0o644))
	p, err := LoadFromDir(dir)
	require.NoError(t, err)
	assert.Equal(t, "x", p.Name)
}
