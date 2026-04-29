package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

func TestProjectsForChangedFiles_MatchesContainingProject(t *testing.T) {
	projects := []*schema.Project{
		{Path: "apps/web", Name: "web"},
		{Path: "clis/foo", Name: "foo"},
	}
	changed := []string{"apps/web/src/index.ts", "README.md"}
	got := projectsForChangedFiles(projects, changed)
	if assert.Len(t, got, 1) {
		assert.Equal(t, "web", got[0].Name)
	}
}

func TestProjectsForChangedFiles_NoMatch(t *testing.T) {
	projects := []*schema.Project{{Path: "apps/web"}}
	got := projectsForChangedFiles(projects, []string{"docs/index.md"})
	assert.Empty(t, got)
}
