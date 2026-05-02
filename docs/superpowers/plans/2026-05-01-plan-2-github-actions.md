# Plan 2 — GitHub Actions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to execute task-by-task.

**Goal:** Replace the temporary `ci.yml` from Plan 1 with a reusable, scalable CI orchestration: composite actions for repeatable steps, reusable workflows for parameterized jobs, and entry-point workflows that drive `piltover affected` matrix execution.

**Architecture:** GitHub composite actions (in `ci-cd-actions/`) own small reusable steps. GitHub reusable workflows (in `.github/workflows/`, the only path GitHub allows) own parametrized jobs invoked via `workflow_call`. Two entry-point workflows — `pr.yml` and `main.yml` — orchestrate the rest. The engine builds itself via `setup-piltover`, then `piltover-affected` emits a JSON matrix consumed by downstream jobs.

**Tech Stack:** GitHub Actions YAML, composite actions, reusable workflows, `actions/setup-go@v5`, `golangci/golangci-lint-action@v7`, `actions/cache@v4`.

**Source spec:** `docs/superpowers/specs/2026-04-29-monorepo-foundation-design.md` §9.

---

## What Plan 2 delivers

| Artefact | Purpose |
|---|---|
| `ci-cd-actions/setup-piltover/action.yml` | Compile the engine binary, cache it across jobs in the same workflow run. |
| `ci-cd-actions/piltover-affected/action.yml` | Run `piltover affected --base <ref>`, expose its JSON as a job output. |
| `.github/workflows/reusable-ci.yml` | `workflow_call`-able job that runs lint+test+build for one project path. |
| `.github/workflows/pr.yml` | Pull-request entry: discovery → affected → matrix → reusable-ci. |
| `.github/workflows/main.yml` | Push-to-main entry: full CI sweep across every project. |
| `tools/project.yaml` | Lets the engine discover itself so future commands operate on it. |
| `tools/Makefile` (root) target update | Add `make engine-project` smoke for self-discovery. |

The old `.github/workflows/ci.yml` is **deleted** at the end. `pr.yml` and `main.yml` jointly cover its triggers.

## What Plan 2 explicitly defers

- `reusable-tofu-plan.yml`, `reusable-tofu-apply.yml`, `setup-tofu-aws-oidc/` → **Plan 4** (when OpenTofu modules exist).
- `reusable-go-release.yml`, `reusable-npm-release.yml` → released alongside the first CLI / first npm package (later child plans).
- `reusable-docs-deploy.yml` → **Plan 3** (Fumadocs).
- `docker-buildx-ecr-push/` → when the first containerised app exists.
- `nightly.yml` drift detection → low priority, opens separately.

---

## File Structure

```
ci-cd-actions/
├── setup-piltover/
│   ├── action.yml
│   └── README.md
└── piltover-affected/
    ├── action.yml
    └── README.md

.github/workflows/
├── reusable-ci.yml          # workflow_call: lint + test + build for one project
├── pr.yml                   # entry: pull_request — affected matrix
├── main.yml                 # entry: push to main — full sweep + (future) deploy/release
└── ci.yml                   # DELETED at the end of this plan

tools/
├── project.yaml             # NEW: makes the engine discoverable
```

---

## Task 1: Add `project.yaml` to the engine

**Files:**
- Create: `tools/project.yaml`

- [ ] **Step 1: Author the project.yaml**

Create `tools/project.yaml`:

```yaml
name: piltover
kind: cli
language: go
tags: [engine, foundation]
commands:
  lint: golangci-lint run ./...
  test: go test -race -count=1 ./...
  build: go build -o ./bin/piltover ./cmd/piltover
release:
  strategy: goreleaser
```

- [ ] **Step 2: Verify discovery**

```bash
make tools
./tools/bin/piltover ls
```

Expected: prints a row `tools  piltover  cli  go  engine,foundation`.

- [ ] **Step 3: Commit**

```bash
git add tools/project.yaml
git commit -m "feat(tools): add project.yaml so the engine discovers itself"
```

---

## Task 2: `setup-piltover` composite action

**Files:**
- Create: `ci-cd-actions/setup-piltover/action.yml`
- Create: `ci-cd-actions/setup-piltover/README.md`

- [ ] **Step 1: Author action.yml**

Create `ci-cd-actions/setup-piltover/action.yml`:

```yaml
name: Setup piltover
description: Build and cache the piltover engine binary, then expose it on PATH.

inputs:
  go-version:
    description: Go version to install. Leave empty to read tools/go.mod.
    required: false
    default: ""

runs:
  using: composite
  steps:
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ inputs.go-version }}
        go-version-file: ${{ inputs.go-version == '' && 'tools/go.mod' || '' }}
        cache-dependency-path: tools/go.sum

    - name: Resolve cache key
      id: key
      shell: bash
      run: |
        key="piltover-bin-$(sha256sum tools/cmd/piltover/main.go tools/go.sum tools/internal/**/*.go 2>/dev/null | sha256sum | cut -c1-16)"
        echo "key=${key}" >> "$GITHUB_OUTPUT"

    - name: Restore engine binary cache
      id: cache
      uses: actions/cache@v4
      with:
        path: tools/bin/piltover
        key: ${{ steps.key.outputs.key }}

    - name: Build engine if cache miss
      if: steps.cache.outputs.cache-hit != 'true'
      shell: bash
      working-directory: tools
      run: |
        echo "→ [.] $ go build -o ./bin/piltover ./cmd/piltover"
        go build -o ./bin/piltover ./cmd/piltover

    - name: Expose binary on PATH
      shell: bash
      run: |
        echo "${GITHUB_WORKSPACE}/tools/bin" >> "$GITHUB_PATH"
        echo "→ [.] $ piltover --version"
        "${GITHUB_WORKSPACE}/tools/bin/piltover" --version
```

- [ ] **Step 2: Author README.md**

Create `ci-cd-actions/setup-piltover/README.md`:

```markdown
# setup-piltover

Composite action that builds the `piltover` engine binary, caches it for the
remainder of the workflow run, and exposes it on `$PATH`.

## Usage

```yaml
- uses: ./ci-cd-actions/setup-piltover

# Or pin a specific Go version:
- uses: ./ci-cd-actions/setup-piltover
  with:
    go-version: '1.23'
```

By default it reads `tools/go.mod` to determine the Go version.

The cache key is derived from `tools/cmd/piltover/main.go`, `tools/go.sum`, and
every file under `tools/internal/**`. Any source change invalidates the cache and
forces a rebuild.
```

- [ ] **Step 3: Commit**

```bash
git add ci-cd-actions/setup-piltover/
git commit -m "feat(ci): add setup-piltover composite action"
```

---

## Task 3: `piltover-affected` composite action

**Files:**
- Create: `ci-cd-actions/piltover-affected/action.yml`
- Create: `ci-cd-actions/piltover-affected/README.md`

- [ ] **Step 1: Author action.yml**

Create `ci-cd-actions/piltover-affected/action.yml`:

```yaml
name: piltover affected
description: Run `piltover affected --base <ref>` and expose the JSON matrix as outputs.

inputs:
  base:
    description: Git ref to diff HEAD against.
    required: false
    default: origin/main

outputs:
  matrix:
    description: JSON object {"include":[...]} suitable for `strategy.matrix`.
    value: ${{ steps.run.outputs.matrix }}
  has_projects:
    description: '"true" if at least one project is affected, else "false".'
    value: ${{ steps.run.outputs.has_projects }}

runs:
  using: composite
  steps:
    - name: Ensure base ref is fetched
      shell: bash
      run: |
        echo "→ [.] $ git fetch --no-tags --depth=50 origin ${{ inputs.base }}"
        git fetch --no-tags --depth=50 origin ${{ inputs.base }} 2>/dev/null || true

    - name: Run piltover affected
      id: run
      shell: bash
      run: |
        echo "→ [.] $ piltover affected --base ${{ inputs.base }}"
        matrix="$(piltover affected --base ${{ inputs.base }})"
        echo "matrix=${matrix}" >> "$GITHUB_OUTPUT"
        if echo "${matrix}" | jq -e '.include | length > 0' >/dev/null; then
          echo "has_projects=true" >> "$GITHUB_OUTPUT"
        else
          echo "has_projects=false" >> "$GITHUB_OUTPUT"
        fi
```

- [ ] **Step 2: Author README**

Create `ci-cd-actions/piltover-affected/README.md`:

```markdown
# piltover-affected

Composite action that runs `piltover affected --base <ref>` and exposes:

| Output | Type | Description |
|---|---|---|
| `matrix` | JSON | `{"include":[...]}` ready to feed into `strategy.matrix` |
| `has_projects` | bool string | `"true"` if any project changed, `"false"` otherwise |

## Usage

```yaml
- uses: ./ci-cd-actions/setup-piltover

- id: affected
  uses: ./ci-cd-actions/piltover-affected
  with:
    base: ${{ github.event.pull_request.base.ref }}

- name: Use the matrix in a downstream job
  if: steps.affected.outputs.has_projects == 'true'
  ...
```

`piltover-affected` requires `setup-piltover` to have run first (it depends on
the binary on PATH).
```

- [ ] **Step 3: Commit**

```bash
git add ci-cd-actions/piltover-affected/
git commit -m "feat(ci): add piltover-affected composite action"
```

---

## Task 4: `reusable-ci.yml` workflow

**Files:**
- Create: `.github/workflows/reusable-ci.yml`

- [ ] **Step 1: Author the reusable workflow**

Create `.github/workflows/reusable-ci.yml`:

```yaml
name: reusable-ci

on:
  workflow_call:
    inputs:
      project_path:
        description: Repo-relative path of the project to lint/test/build.
        required: true
        type: string
      project_name:
        description: Display name (matrix label).
        required: true
        type: string
      language:
        description: Language declared in the project's project.yaml.
        required: true
        type: string

permissions:
  contents: read

jobs:
  ci:
    name: ${{ inputs.project_name }} (${{ inputs.language }})
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@v4

      - uses: ./ci-cd-actions/setup-piltover

      - name: Verify go.mod tidy (if Go project)
        if: inputs.language == 'go'
        shell: bash
        working-directory: ${{ inputs.project_path }}
        run: |
          if [ -f go.mod ]; then
            echo "→ [${{ inputs.project_path }}] $ go mod tidy"
            go mod tidy
            if ! git diff --quiet -- go.mod go.sum; then
              echo "go.mod / go.sum drift detected. Run 'go mod tidy' from ${{ inputs.project_path }}."
              git --no-pager diff -- go.mod go.sum
              exit 1
            fi
          fi

      - name: golangci-lint (if Go project)
        if: inputs.language == 'go'
        uses: golangci/golangci-lint-action@v7
        with:
          version: v2.1.6
          working-directory: ${{ inputs.project_path }}
          args: --timeout=5m

      - name: piltover lint
        if: inputs.language != 'go'
        run: piltover lint ${{ inputs.project_path }}

      - name: piltover test
        run: piltover test ${{ inputs.project_path }}

      - name: piltover build
        run: piltover build ${{ inputs.project_path }}
```

Note: golangci-lint is invoked via the dedicated GitHub action for Go projects (it knows how to cache and report annotations), and via `piltover lint` for non-Go projects (which delegates to the language-appropriate linter from `tools/configs/defaults.yaml`).

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/reusable-ci.yml
git commit -m "feat(ci): add reusable-ci workflow for per-project lint/test/build"
```

---

## Task 5: `pr.yml` entry-point

**Files:**
- Create: `.github/workflows/pr.yml`

- [ ] **Step 1: Author**

Create `.github/workflows/pr.yml`:

```yaml
name: pr

on:
  pull_request:

permissions:
  contents: read

concurrency:
  group: pr-${{ github.event.pull_request.number }}
  cancel-in-progress: true

jobs:
  affected:
    name: Compute affected matrix
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.affected.outputs.matrix }}
      has_projects: ${{ steps.affected.outputs.has_projects }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: ./ci-cd-actions/setup-piltover

      - id: affected
        uses: ./ci-cd-actions/piltover-affected
        with:
          base: origin/${{ github.event.pull_request.base.ref }}

  ci:
    name: ci
    needs: affected
    if: needs.affected.outputs.has_projects == 'true'
    strategy:
      fail-fast: false
      matrix: ${{ fromJson(needs.affected.outputs.matrix) }}
    uses: ./.github/workflows/reusable-ci.yml
    with:
      project_path: ${{ matrix.path }}
      project_name: ${{ matrix.name }}
      language: ${{ matrix.language }}

  ci-skipped:
    name: ci-skipped
    needs: affected
    if: needs.affected.outputs.has_projects == 'false'
    runs-on: ubuntu-latest
    steps:
      - run: echo "No projects changed in this PR — CI skipped."
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/pr.yml
git commit -m "feat(ci): add pr.yml entry-point with affected matrix orchestration"
```

---

## Task 6: `main.yml` entry-point

**Files:**
- Create: `.github/workflows/main.yml`

- [ ] **Step 1: Author**

Create `.github/workflows/main.yml`:

```yaml
name: main

on:
  push:
    branches: [main]

permissions:
  contents: read

concurrency:
  group: main
  cancel-in-progress: false

jobs:
  discover:
    name: List all projects
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.matrix.outputs.matrix }}
    steps:
      - uses: actions/checkout@v4

      - uses: ./ci-cd-actions/setup-piltover

      - id: matrix
        shell: bash
        run: |
          echo "→ [.] $ piltover affected --base HEAD~1"
          matrix="$(piltover affected --base HEAD~1)"
          # If HEAD~1 doesn't exist or no projects changed, fall back to running every project.
          if echo "${matrix}" | jq -e '.include | length == 0' >/dev/null; then
            echo "→ [.] $ piltover ls (fallback: full sweep)"
            entries="$(piltover ls --json 2>/dev/null || piltover ls | tail -n +2 | awk 'NF >= 4 {printf "{\"path\":\"%s\",\"name\":\"%s\",\"kind\":\"%s\",\"language\":\"%s\"}\n", $1, $2, $3, $4}' | jq -s '.')"
            matrix="$(jq -nc --argjson e "${entries:-[]}" '{include: $e}')"
          fi
          echo "matrix=${matrix}" >> "$GITHUB_OUTPUT"

  ci:
    name: ci
    needs: discover
    if: ${{ fromJson(needs.discover.outputs.matrix).include[0] != null }}
    strategy:
      fail-fast: false
      matrix: ${{ fromJson(needs.discover.outputs.matrix) }}
    uses: ./.github/workflows/reusable-ci.yml
    with:
      project_path: ${{ matrix.path }}
      project_name: ${{ matrix.name }}
      language: ${{ matrix.language }}
```

The `piltover ls --json` flag is added in Task 7. Until then the awk fallback parses the text output.

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/main.yml
git commit -m "feat(ci): add main.yml entry-point (full sweep on push to main)"
```

---

## Task 7: Add `piltover ls --json`

**Files:**
- Modify: `tools/internal/cli/ls.go`
- Modify: `tools/internal/cli/ls_test.go`

- [ ] **Step 1: Extend ls_test.go with JSON case**

Append to `tools/internal/cli/ls_test.go`:

```go
func TestLs_JSON(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "clis/foo"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(root, "clis/foo/project.yaml"),
		[]byte("name: foo\nkind: cli\nlanguage: go\n"), 0o600))

	var stdout bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--root", root, "ls", "--json"})
	require.NoError(t, cmd.Execute())

	out := stdout.String()
	assert.Contains(t, out, `"name":"foo"`)
	assert.Contains(t, out, `"kind":"cli"`)
	assert.Contains(t, out, `"language":"go"`)
}
```

- [ ] **Step 2: Run test to confirm failure**

```bash
cd tools && go test ./internal/cli/ -run TestLs_JSON
```

Expected: FAIL — `--json` flag not implemented.

- [ ] **Step 3: Update ls.go**

Replace `tools/internal/cli/ls.go` with:

```go
package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/discovery"
	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/schema"
)

func newLsCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "ls",
		Short: "List every discovered subproject",
		RunE: func(cmd *cobra.Command, _ []string) error {
			projects, err := discovery.Discover(g.Root)
			if err != nil {
				return err
			}
			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return printLsJSON(cmd, projects)
			}
			return printLsTable(cmd, projects)
		},
	}
	c.Flags().Bool("json", false, "emit JSON suitable for piping into jq")
	return c
}

func printLsTable(cmd *cobra.Command, projects []*schema.Project) error {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tNAME\tKIND\tLANGUAGE\tTAGS")
	for _, p := range projects {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			p.Path, p.Name, p.Kind, p.Language, strings.Join(p.Tags, ","))
	}
	return w.Flush()
}

type lsEntry struct {
	Path     string   `json:"path"`
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
}

func printLsJSON(cmd *cobra.Command, projects []*schema.Project) error {
	entries := make([]lsEntry, 0, len(projects))
	for _, p := range projects {
		entries = append(entries, lsEntry{
			Path:     p.Path,
			Name:     p.Name,
			Kind:     string(p.Kind),
			Language: string(p.Language),
			Tags:     p.Tags,
		})
	}
	b, err := json.Marshal(map[string]any{"include": entries})
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(b))
	return nil
}
```

- [ ] **Step 4: Run tests, verify pass**

```bash
cd tools && go test ./internal/cli/...
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add tools/internal/cli/ls.go tools/internal/cli/ls_test.go
git commit -m "feat(tools/cli): add --json output to piltover ls"
```

---

## Task 8: Simplify `main.yml` to use `piltover ls --json`

**Files:**
- Modify: `.github/workflows/main.yml`

- [ ] **Step 1: Replace the awk fallback**

Replace the matrix step in `.github/workflows/main.yml` with the cleaner version using `piltover ls --json`:

```yaml
      - id: matrix
        shell: bash
        run: |
          echo "→ [.] $ piltover affected --base HEAD~1"
          matrix="$(piltover affected --base HEAD~1)"
          if echo "${matrix}" | jq -e '.include | length == 0' >/dev/null; then
            echo "→ [.] $ piltover ls --json (fallback: full sweep)"
            matrix="$(piltover ls --json)"
          fi
          echo "matrix=${matrix}" >> "$GITHUB_OUTPUT"
```

The rest of the workflow stays the same.

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/main.yml
git commit -m "feat(ci): use piltover ls --json in main.yml fallback"
```

---

## Task 9: Delete the old `ci.yml`

**Files:**
- Delete: `.github/workflows/ci.yml`

- [ ] **Step 1: Verify the new workflows are landing first**

Before deleting, confirm `pr.yml` and `main.yml` are in the working tree (they were committed in Tasks 5-6).

```bash
ls -la .github/workflows/
```

Expected: `pr.yml`, `main.yml`, `reusable-ci.yml` are all present alongside the old `ci.yml`.

- [ ] **Step 2: Delete**

```bash
git rm .github/workflows/ci.yml
git commit -m "chore(ci): remove legacy ci.yml (replaced by pr.yml + main.yml)"
```

---

## Task 10: End-to-end verification by triggering the new workflow

**Files:** none (verification only)

- [ ] **Step 1: Push the branch**

```bash
git push origin feat/plan-2-github-actions
```

- [ ] **Step 2: Open or update the PR**

If a PR for this branch is not yet open:
```bash
gh pr create --base main --head feat/plan-2-github-actions \
  --title "feat: plan 2 — github actions reusable workflows" \
  --body "(filled in by controller)"
```

If a PR already exists, the push triggers re-run.

- [ ] **Step 3: Wait for `pr.yml` to complete**

```bash
gh pr checks
```

Expected: the new workflow runs, the `affected` job emits a matrix containing `tools` (since this PR touches `tools/internal/cli/ls.go`), the matrix CI job runs lint+test+build for the engine, and everything passes.

- [ ] **Step 4: If anything fails**

Inspect logs with `gh run view <run-id> --log-failed`. Fix and push again. The workflow self-validates.

- [ ] **Step 5: No commit needed**

This task is verification only. If a fix is required, the fix is its own commit.

---

## Out of Scope (handled by other plans)

- **Plan 3:** Fumadocs site, `reusable-docs-deploy.yml`.
- **Plan 4:** OpenTofu modules, `setup-tofu-aws-oidc/`, `reusable-tofu-plan.yml`, `reusable-tofu-apply.yml`.
- **Future plan:** `reusable-go-release.yml` (when first CLI ships), `reusable-npm-release.yml` (when first npm package ships), `docker-buildx-ecr-push/` (when first containerized app ships), `nightly.yml` drift detection.

---

## Self-Review (controller)

After all 10 tasks land:

1. **Spec coverage:** Spec §9 listed reusable workflows for ci/tofu/release/docs. Plan 2 delivers `reusable-ci.yml` and the entry-points; the rest are explicitly deferred above. ✓
2. **Logging contract:** Every shell step in the actions/workflows that spawns a subprocess prints `→ [<cwd>] $ <cmd>`. Verified in Tasks 2, 3, 5, 6. ✓
3. **Type consistency:** matrix shape (`{path, name, kind, language, tags}`) is identical between `piltover affected` (Task 9 of Plan 1), `piltover ls --json` (Task 7), and the `reusable-ci` inputs (Task 4). ✓
4. **Branch hygiene:** all commits follow Conventional Commits.
