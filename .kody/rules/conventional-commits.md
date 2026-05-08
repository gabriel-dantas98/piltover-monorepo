---
title: "Use Conventional Commits"
scope: "commit"
path: ["**/*"]
severity_min: "medium"
languages: ["go", "ts", "python", "hcl", "shell"]
buckets: ["style-conventions"]
enabled: true
---

Every commit follows Conventional Commits: `feat:`, `fix:`, `docs:`, `chore:`,
`refactor:`, `test:`, `ci:`, `build:`, `perf:`, with an optional scope in
parentheses (e.g. `feat(tools/cli): add --json flag`). Commitlint enforces the
header rules; the `lefthook commit-msg` hook runs it locally on every commit.
