package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minimalFixture = `---
title: "No secrets in code"
scope: "pull_request"
path:
  - "**/*.go"
  - "**/*.ts"
severity_min: "medium"
languages:
  - go
  - typescript
buckets:
  - security
enabled: true
---

Do not commit secrets or credentials directly in source code.
`

func TestParse_Minimal(t *testing.T) {
	r, err := Parse([]byte(minimalFixture))
	require.NoError(t, err)
	assert.Equal(t, "No secrets in code", r.Title)
	assert.Equal(t, "pull_request", r.Scope)
	assert.Equal(t, []string{"**/*.go", "**/*.ts"}, r.Path)
	assert.Equal(t, "medium", r.SeverityMin)
	assert.Equal(t, []string{"go", "typescript"}, r.Languages)
	assert.Equal(t, []string{"security"}, r.Buckets)
	assert.True(t, r.Enabled)
	assert.Contains(t, r.Body, "Do not commit secrets")
}

func TestParse_MissingFrontmatter(t *testing.T) {
	_, err := Parse([]byte("This is just markdown without frontmatter\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing frontmatter fence")
}

func TestParse_InvalidYAML(t *testing.T) {
	bad := "---\ntitle: [\nbad yaml here\n---\n\nbody\n"
	_, err := Parse([]byte(bad))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "yaml")
}

func TestValidate_RequiresAllFields(t *testing.T) {
	r := &Rule{
		Slug:  "missing-title",
		Scope: "pull_request",
		Path:  []string{"**/*.go"},
	}
	err := r.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title")
}
