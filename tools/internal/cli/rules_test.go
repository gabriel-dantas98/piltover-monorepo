package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fixtureRule1 = `---
title: "No secrets in code"
scope: "pull_request"
path:
  - "**/*.go"
severity_min: "medium"
languages:
  - go
buckets:
  - security
enabled: true
---

Do not commit secrets or credentials directly in source code.
`

const fixtureRule2 = `---
title: "Require error handling"
scope: "pull_request"
path:
  - "**/*.go"
severity_min: "high"
languages:
  - go
buckets:
  - quality
enabled: false
---

All returned errors must be handled.
`

const fixtureRuleInvalid = `---
scope: "pull_request"
path:
  - "**/*.go"
severity_min: "medium"
languages:
  - go
buckets:
  - security
enabled: true
---

Missing title field.
`

func writeRuleFixture(t *testing.T, root, name, content string) {
	t.Helper()
	dir := filepath.Join(root, ".kody", "rules")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}

func TestRules_Ls(t *testing.T) {
	root := t.TempDir()
	writeRuleFixture(t, root, "no-secrets.md", fixtureRule1)
	writeRuleFixture(t, root, "require-error-handling.md", fixtureRule2)

	var stdout bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--root", root, "rules", "ls"})
	require.NoError(t, cmd.Execute())

	out := stdout.String()
	assert.Contains(t, out, "no-secrets")
	assert.Contains(t, out, "require-error-handling")
	assert.Contains(t, out, "pull_request")
	assert.Contains(t, out, "medium")
	assert.Contains(t, out, "high")
}

func TestRules_Lint_Pass(t *testing.T) {
	root := t.TempDir()
	writeRuleFixture(t, root, "no-secrets.md", fixtureRule1)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "rules", "lint"})
	require.NoError(t, cmd.Execute())

	assert.Contains(t, stdout.String(), "ok")
}

func TestRules_Lint_Fails(t *testing.T) {
	root := t.TempDir()
	writeRuleFixture(t, root, "bad-rule.md", fixtureRuleInvalid)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--root", root, "rules", "lint"})
	err := cmd.Execute()
	require.Error(t, err)

	errOut := stderr.String()
	assert.Contains(t, errOut, "title")
}

func TestRules_SyncDocs(t *testing.T) {
	root := t.TempDir()
	writeRuleFixture(t, root, "no-secrets.md", fixtureRule1)
	writeRuleFixture(t, root, "require-error-handling.md", fixtureRule2)

	// Write a stale .mdx that should be deleted.
	outDir := filepath.Join(root, "apps", "docs", "content", "rules")
	require.NoError(t, os.MkdirAll(outDir, 0o750))
	stalePath := filepath.Join(outDir, "old-stale-rule.mdx")
	require.NoError(t, os.WriteFile(stalePath, []byte("stale content"), 0o600))

	var stdout bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--root", root, "rules", "sync-docs"})
	require.NoError(t, cmd.Execute())

	out := stdout.String()
	assert.Contains(t, out, "synced")
	assert.Contains(t, out, "2")

	// Assert generated files exist with expected content.
	mdx1, err := os.ReadFile(filepath.Join(outDir, "no-secrets.mdx")) // #nosec G304 -- test reads from t.TempDir()
	require.NoError(t, err)
	content1 := string(mdx1)
	assert.Contains(t, content1, "No secrets in code")
	assert.Contains(t, content1, "security")
	assert.Contains(t, content1, "pull_request")

	mdx2, err := os.ReadFile(filepath.Join(outDir, "require-error-handling.mdx")) // #nosec G304 -- test reads from t.TempDir()
	require.NoError(t, err)
	content2 := string(mdx2)
	assert.Contains(t, content2, "Require error handling")

	// Stale file should be deleted.
	_, err = os.Stat(stalePath)
	assert.True(t, os.IsNotExist(err), "stale .mdx should have been deleted")

	// Verify the no-secrets.mdx contains MDX structure.
	assert.Contains(t, content1, "import { Card }")
	assert.True(t, strings.HasPrefix(content1, "---\n"), "should start with frontmatter")
}
