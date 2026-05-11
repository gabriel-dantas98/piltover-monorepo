# Plan 5 — Docker stacks + `piltover stacks` + `piltover new` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to execute task-by-task.

**Goal:** Replace the v0 stubs for `piltover stacks` and `piltover new` with real implementations, and ship the two canonical local-dev docker-compose stacks (`postgres` and `localstack`).

**Architecture:** `piltover stacks` is a thin wrapper around `docker compose -f docker-stacks/<name>/compose.yaml` honouring the engine's logging contract. `piltover new` scaffolds subprojects (`app`, `cli`, `package`, `stack`) from embedded templates into the right top-level folder, dropping a `project.yaml` so the engine discovers it immediately.

**Tech Stack:** Go 1.23, Cobra, `text/template` (stdlib), `docker compose` v2, postgres 17, [localstack](https://docs.localstack.cloud/).

**Source spec:** `docs/superpowers/specs/2026-04-29-monorepo-foundation-design.md` §4.3, §8.

---

## What Plan 5 delivers

| Artefact | Purpose |
|---|---|
| `tools/internal/stacks/` | Helpers to list, validate, and dispatch docker-compose stacks. |
| `tools/internal/cli/stacks.go` | Real `piltover stacks ls\|up\|down\|nuke <name>` command tree (replaces v0 stub). |
| `tools/internal/cli/stacks_test.go` | Unit tests with a stubbed runner. |
| `tools/internal/scaffold/` | Embedded templates + writer helpers. |
| `tools/internal/cli/new.go` | Real `piltover new <kind> <name>` (replaces v0 stub). |
| `tools/internal/cli/new_test.go` | Unit tests over a temp dir. |
| `tools/templates/cli-go/` | Embedded template tree for `new cli`. |
| `tools/templates/package-ts/` | Embedded template tree for `new package`. |
| `tools/templates/app-ts/` | Embedded template tree for `new app` (minimal Next stub). |
| `tools/templates/stack/` | Embedded template tree for `new stack`. |
| `docker-stacks/postgres/` | Postgres 17 + pgAdmin local stack. |
| `docker-stacks/localstack/` | localstack community for SQS/S3/SNS/DDB/Lambda emulation. |

## What Plan 5 explicitly defers

- `piltover tf` real impl + OpenTofu modules → **Plan 4**.
- More docker stacks (redis, observability, kafka, ai-local) — added on demand.
- `new plugin`, `new action`, `new infra-module` — same shape, can extend later.

---

## File Structure (this plan only)

```
tools/
├── configs/defaults.yaml                   (unchanged)
├── internal/
│   ├── stacks/
│   │   ├── stacks.go                       List+validate stacks under docker-stacks/
│   │   └── stacks_test.go
│   ├── scaffold/
│   │   ├── scaffold.go                     Render embedded template → tree of files
│   │   └── scaffold_test.go
│   ├── cli/
│   │   ├── stacks.go                       newStacksCmd: ls/up/down/nuke
│   │   ├── stacks_test.go
│   │   ├── new.go                          newScaffoldCmd: new <kind> <name>
│   │   └── new_test.go
│   └── …
└── templates/
    ├── cli-go/
    │   ├── project.yaml.tmpl
    │   ├── go.mod.tmpl
    │   ├── cmd/{{.Name}}/main.go.tmpl
    │   └── README.md.tmpl
    ├── package-ts/
    │   ├── project.yaml.tmpl
    │   ├── package.json.tmpl
    │   ├── tsconfig.json.tmpl
    │   ├── src/index.ts.tmpl
    │   └── README.md.tmpl
    ├── app-ts/
    │   ├── project.yaml.tmpl
    │   ├── package.json.tmpl
    │   ├── tsconfig.json.tmpl
    │   ├── app/page.tsx.tmpl
    │   ├── app/layout.tsx.tmpl
    │   ├── next.config.mjs.tmpl
    │   └── README.md.tmpl
    └── stack/
        ├── compose.yaml.tmpl
        ├── .env.example.tmpl
        └── README.md.tmpl

docker-stacks/
├── postgres/
│   ├── compose.yaml
│   ├── .env.example
│   └── README.md
└── localstack/
    ├── compose.yaml
    ├── .env.example
    └── README.md
```

---

## Task 1: Stacks discovery package (TDD)

**Files:**
- Create: `tools/internal/stacks/stacks.go`
- Test: `tools/internal/stacks/stacks_test.go`

`stacks.List(root) ([]Stack, error)` walks `docker-stacks/*/compose.yaml`. `Stack` is
`{Name, Path, ComposeFile}`. `stacks.Resolve(root, name)` returns one or an error.

- [ ] **Step 1: Failing tests**

Tests cover: empty dir, valid stack discovery, name resolution, "stack not found"
error. Use `t.TempDir()` + fixture compose files.

- [ ] **Step 2: Implement**

```go
package stacks

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
)

type Stack struct {
	Name        string
	Path        string // repo-relative path to stack dir
	ComposeFile string // repo-relative path to compose.yaml
}

func List(root string) ([]Stack, error) {
	dir := filepath.Join(root, "docker-stacks")
	entries, err := fs.ReadDir(rootFS(root), "docker-stacks")
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			return nil, nil
		}
		return nil, err
	}
	var out []Stack
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		compose := filepath.Join(dir, e.Name(), "compose.yaml")
		if _, err := fs.Stat(rootFS(root), filepath.Join("docker-stacks", e.Name(), "compose.yaml")); err != nil {
			continue
		}
		out = append(out, Stack{
			Name:        e.Name(),
			Path:        filepath.Join("docker-stacks", e.Name()),
			ComposeFile: filepath.Join("docker-stacks", e.Name(), "compose.yaml"),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func Resolve(root, name string) (Stack, error) {
	all, err := List(root)
	if err != nil {
		return Stack{}, err
	}
	for _, s := range all {
		if s.Name == name {
			return s, nil
		}
	}
	return Stack{}, fmt.Errorf("stack %q not found under docker-stacks/", name)
}
```

`rootFS` is `os.DirFS(root)`. Adjust the imports accordingly.

- [ ] **Step 3: Run tests, commit.**

```bash
cd tools && go test ./internal/stacks/...
git add tools/internal/stacks/
git commit -m "feat(tools/stacks): add docker-stacks discovery package"
```

---

## Task 2: `piltover stacks` real command

**Files:**
- Modify: `tools/internal/cli/stacks.go` (replace stub)
- Modify: `tools/internal/cli/root.go` (use `newStacksCmd(g)` instead of `newStubCmd`)
- Test: `tools/internal/cli/stacks_test.go`

Subcommands: `ls`, `up <name>`, `down <name>`, `nuke <name>` (down with `-v`).

`up`/`down`/`nuke` invoke `docker compose -f <compose> up -d` / `down` / `down -v`
through the existing `runner.Runner`. `ls` prints a tabwriter of `NAME / PATH`.

Tests use a fake runner that records calls instead of executing them. Reuse the
runner's `--dry-run` mode where possible.

- [ ] **Step 1: Failing tests** for: `ls` output shape; `up postgres` issues exactly
  `docker compose -f docker-stacks/postgres/compose.yaml up -d` via the runner;
  `nuke` adds `-v`; "stack not found" exits non-zero with a clear message.

- [ ] **Step 2: Implement** the Cobra subcommands.

- [ ] **Step 3: Wire into root.** Replace
  `root.AddCommand(newStubCmd("stacks", ...))` with `root.AddCommand(newStacksCmd(g))`.

- [ ] **Step 4: `make verify` clean, commit.**

```bash
git commit -m "feat(tools/cli): implement piltover stacks ls|up|down|nuke"
```

---

## Task 3: Scaffold package (TDD)

**Files:**
- Create: `tools/internal/scaffold/scaffold.go`
- Test: `tools/internal/scaffold/scaffold_test.go`

Exposes `scaffold.Render(kind, name, destDir, fs.FS) error`. Walks the template tree
(embedded `fs.FS`), executes each file with `text/template` against the input
`Vars{Name string, KebabName string, GoModule string, …}`, and writes the result to
`destDir/<rendered path>`. Removes the `.tmpl` suffix from filenames. Template paths
themselves may contain `{{.Name}}` placeholders (e.g. `cmd/{{.Name}}/main.go.tmpl`).

- [ ] **Step 1: Failing tests** with a fixture `fs.FS` (use `fstest.MapFS`).

- [ ] **Step 2: Implement.**

- [ ] **Step 3: Commit.**

---

## Task 4: Embedded templates

**Files:**
- Create: `tools/templates/cli-go/**`
- Create: `tools/templates/package-ts/**`
- Create: `tools/templates/app-ts/**`
- Create: `tools/templates/stack/**`

Each template tree is a minimal but **working** project:

- `cli-go/`: `project.yaml` (kind=cli, language=go), `go.mod`, `cmd/<name>/main.go`
  printing "hello from <name>", `README.md`.
- `package-ts/`: `project.yaml` (kind=package, language=ts), `package.json` with
  biome+vitest, `tsconfig.json`, `src/index.ts` exporting a hello function,
  `README.md`.
- `app-ts/`: `project.yaml`, `package.json` (Next.js + Fumadocs-compatible), `app/`
  with a single page, `next.config.mjs`, `tsconfig.json`, `README.md`. Aim for
  minimal; the user will customize.
- `stack/`: `compose.yaml` placeholder, `.env.example`, `README.md`.

Each template uses `{{.Name}}` and `{{.KebabName}}` consistently. The package name
field in `project.yaml` is the kebab-case form.

- [ ] **Step 1: Write all templates.**
- [ ] **Step 2: Add `//go:embed all:templates` directive** to a Go file in
  `tools/internal/scaffold/`. Confirm `go build ./...` still works.
- [ ] **Step 3: Commit.**

---

## Task 5: `piltover new` real command

**Files:**
- Modify: `tools/internal/cli/new.go` (replace stub)
- Modify: `tools/internal/cli/root.go` (use `newScaffoldCmd(g)`)
- Test: `tools/internal/cli/new_test.go`

Routes:

| Invocation | Destination |
|---|---|
| `piltover new cli <name>` | `clis/<name>/` |
| `piltover new package <name>` | `packages/<name>/` |
| `piltover new app <name>` | `apps/<name>/` |
| `piltover new stack <name>` | `docker-stacks/<name>/` |

Other kinds (plugin / action / infra-module) print "not implemented for v1; create
manually until extended."

The command refuses to overwrite an existing destination (returns an error). It
prints a "next steps" hint at the end (`piltover ls`, `make tools`).

- [ ] **Step 1: Failing tests** for each route + the destination-exists error.

- [ ] **Step 2: Implement.**

- [ ] **Step 3: Wire into root, lint+test+build, commit.**

---

## Task 6: `docker-stacks/postgres/`

**Files:**
- Create: `docker-stacks/postgres/compose.yaml`
- Create: `docker-stacks/postgres/.env.example`
- Create: `docker-stacks/postgres/README.md`

### Step 6.1: `compose.yaml`

```yaml
name: piltover-postgres

services:
  postgres:
    image: postgres:17-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-piltover}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-piltover}
      POSTGRES_DB: ${POSTGRES_DB:-piltover}
    ports:
      - "${POSTGRES_PORT:-5432}:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-piltover}"]
      interval: 5s
      timeout: 5s
      retries: 5

  pgadmin:
    image: dpage/pgadmin4:latest
    restart: unless-stopped
    profiles: ["ui"]
    environment:
      PGADMIN_DEFAULT_EMAIL: ${PGADMIN_EMAIL:-admin@piltover.local}
      PGADMIN_DEFAULT_PASSWORD: ${PGADMIN_PASSWORD:-admin}
    ports:
      - "${PGADMIN_PORT:-5050}:80"
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  pgdata:
```

### Step 6.2: `.env.example`

```env
POSTGRES_USER=piltover
POSTGRES_PASSWORD=piltover
POSTGRES_DB=piltover
POSTGRES_PORT=5432

# Optional UI (enable via: docker compose --profile ui up)
PGADMIN_EMAIL=admin@piltover.local
PGADMIN_PASSWORD=admin
PGADMIN_PORT=5050
```

### Step 6.3: `README.md`

Document:
- What's inside (postgres 17 + optional pgAdmin profile).
- Default ports + credentials.
- Connection string examples (`postgres://piltover:piltover@localhost:5432/piltover`).
- `piltover stacks up postgres` to start.
- Reset everything: `piltover stacks nuke postgres`.

- [ ] Commit: `feat(docker-stacks): add postgres stack (postgres 17 + optional pgAdmin)`

---

## Task 7: `docker-stacks/localstack/`

**Files:**
- Create: `docker-stacks/localstack/compose.yaml`
- Create: `docker-stacks/localstack/.env.example`
- Create: `docker-stacks/localstack/README.md`

### Step 7.1: `compose.yaml`

```yaml
name: piltover-localstack

services:
  localstack:
    image: localstack/localstack:latest
    restart: unless-stopped
    ports:
      - "${LOCALSTACK_PORT:-4566}:4566"
    environment:
      SERVICES: ${LOCALSTACK_SERVICES:-s3,sqs,sns,dynamodb,lambda,iam,sts,logs}
      DEBUG: ${LOCALSTACK_DEBUG:-0}
      PERSISTENCE: ${LOCALSTACK_PERSISTENCE:-1}
      DOCKER_HOST: unix:///var/run/docker.sock
    volumes:
      - localstack-data:/var/lib/localstack
      - /var/run/docker.sock:/var/run/docker.sock

volumes:
  localstack-data:
```

### Step 7.2: `.env.example`

```env
LOCALSTACK_PORT=4566
LOCALSTACK_SERVICES=s3,sqs,sns,dynamodb,lambda,iam,sts,logs
LOCALSTACK_DEBUG=0
LOCALSTACK_PERSISTENCE=1
```

### Step 7.3: `README.md`

Document:
- Endpoint: `http://localhost:4566`
- AWS CLI shim: `aws --endpoint-url=http://localhost:4566 s3 ls`
- Snippet to create an S3 bucket against it.
- Reset: `piltover stacks nuke localstack`.

- [ ] Commit: `feat(docker-stacks): add localstack stack for AWS local emulation`

---

## Task 8: Documentation updates

**Files:**
- Modify: `apps/docs/content/repo/engine.mdx` (move `stacks`/`new` out of "stubs")
- Modify: `apps/docs/content/agents/engine-api.mdx` (same)
- Modify: `AGENTS.md` (status column: `piltover stacks` → productized, `piltover new` → productized)
- Create: `apps/docs/content/guides/local-dev-stacks.mdx` (how to use the stacks)
- Create: `apps/docs/content/guides/scaffold-a-subproject.mdx` (how `piltover new` works)
- Modify: `apps/docs/content/meta.json` (add the 2 new guides under "Guides")

- [ ] **Step 1: Author the 2 guide pages.** Use the conventions of the existing
  guides (frontmatter, code fences, mermaid where useful).

- [ ] **Step 2: Update engine.mdx / engine-api.mdx / AGENTS.md status tables.**

- [ ] **Step 3: Update `meta.json` sidebar order.**

- [ ] **Step 4: `bun run build` in `apps/docs`, commit.**

---

## Task 9: End-to-end smoke

**Files:** none (verification only)

- [ ] **Step 1:** `make tools && piltover --version`
- [ ] **Step 2:** `piltover stacks ls` shows postgres + localstack.
- [ ] **Step 3:** `piltover stacks up postgres` (start). Verify `pg_isready` works
  through `docker compose ps`.
- [ ] **Step 4:** `piltover stacks down postgres` (clean stop).
- [ ] **Step 5:** `piltover new cli example-cli`. Confirm `clis/example-cli/` exists
  with a working `project.yaml`. Run `piltover lint clis/example-cli` to validate.
  Remove the scaffolded dir.
- [ ] **Step 6:** `piltover --dry-run stacks up localstack` — prints the docker
  command without running it.

If everything is green, open the PR.

---

## Task 10: Open PR + merge

- [ ] **Step 1:** `git push -u origin feat/plan-5-stacks-and-scaffolds`
- [ ] **Step 2:** Open PR via `gh pr create`. Body lists the new commands + stacks +
  docs.
- [ ] **Step 3:** Wait for CI green. Merge.

---

## Self-Review (controller)

1. **Spec coverage:** Plan 5 closes spec §4.3 stubs `stacks`/`new` and spec §8 (local
   docker-compose stacks). ✓
2. **Logging contract:** `stacks up/down/nuke` go through the runner. `new` performs
   filesystem writes only — log them too (one line per file written) to keep the
   discipline.
3. **Type consistency:** Stack / Vars structs are referenced consistently across the
   stacks/, scaffold/, and cli packages. ✓
4. **YAGNI:** No premature features (no stack composition, no template registry, no
   `new` for plugins/actions). Defer to demand.
