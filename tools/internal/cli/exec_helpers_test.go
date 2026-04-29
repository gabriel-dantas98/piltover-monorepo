package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

func TestResolveCommand_PrefersOverride(t *testing.T) {
	defaults := map[schema.Language]LanguageDefaults{
		schema.LangGo: {Lint: "golangci-lint run ./..."},
	}
	p := &schema.Project{
		Language: schema.LangGo,
		Commands: schema.Commands{Lint: "golangci-lint run --fast ./..."},
	}
	got := resolveCommand(p, "lint", defaults)
	assert.Equal(t, "golangci-lint run --fast ./...", got)
}

func TestResolveCommand_FallsBackToDefault(t *testing.T) {
	defaults := map[schema.Language]LanguageDefaults{
		schema.LangGo: {Lint: "golangci-lint run ./..."},
	}
	p := &schema.Project{Language: schema.LangGo}
	got := resolveCommand(p, "lint", defaults)
	assert.Equal(t, "golangci-lint run ./...", got)
}

func TestResolveCommand_EmptyForLangNoneOrUnset(t *testing.T) {
	defaults := map[schema.Language]LanguageDefaults{
		schema.LangNone: {Lint: ""},
	}
	p := &schema.Project{Language: schema.LangNone}
	assert.Equal(t, "", resolveCommand(p, "lint", defaults))
}
