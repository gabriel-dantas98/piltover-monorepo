---
title: "Wrappers must log the underlying command before executing"
scope: "file"
path: ["tools/**/*.go", "clis/**/*.go"]
severity_min: "high"
languages: ["go"]
buckets: ["style-conventions", "operability"]
enabled: true
---

Any helper, runner, or wrapper that spawns an external process MUST print to stderr,
before invoking the process, the exact command in the form:

    → [<cwd>] $ <name> <args...>

The contract honours `--quiet` (suppress the line), `--verbose` (also print env vars),
and `--dry-run` (print but skip execution). The runner package at
`tools/internal/runner` is the canonical implementation; new wrappers should reuse it
rather than calling `exec.Command` directly.
