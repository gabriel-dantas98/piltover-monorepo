package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunChecks_PopulatesStatus(t *testing.T) {
	checks := []Check{
		{Name: "always-ok", Cmd: "true"},
		{Name: "always-fail", Cmd: "this-binary-does-not-exist-xyz"},
	}
	results := runChecks(checks)
	require.Len(t, results, 2)
	assert.True(t, results[0].OK, "expected first check to pass")
	assert.False(t, results[1].OK, "expected second check to fail")
}

func TestRenderJSON_Shape(t *testing.T) {
	results := []CheckResult{
		{Name: "go", OK: true, Detail: "go version go1.23"},
		{Name: "tofu", OK: false, Detail: "exec: not found"},
	}
	out := renderJSON(results)
	var parsed []CheckResult
	require.NoError(t, json.Unmarshal([]byte(out), &parsed))
	assert.Equal(t, results, parsed)
}
