# tools/

Home of the `piltover` engine — a thin Go wrapper that discovers subprojects,
runs their lint/test/build, orchestrates CI, and logs every underlying command
it executes.

## Build

From the repo root:

```bash
make tools
```

Or directly:

```bash
cd tools && go build -o ./bin/piltover ./cmd/piltover
```

The compiled binary is gitignored. The canonical install location is your `$GOBIN`
(set via `make tools`).

## Layout

- `cmd/piltover/` — entry point.
- `internal/runner/` — logged-exec wrapper. Every subprocess spawn must go through this package.
- `internal/schema/` — `project.yaml` types + parser.
- `internal/discovery/` — repo walker that finds `project.yaml` files.
- `internal/cli/` — Cobra command implementations.
- `configs/` — shared lint/format default configs.

## Tests

```bash
cd tools && go test ./...
```
