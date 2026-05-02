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

The cache key is derived from every `.go` file under `tools/cmd` and
`tools/internal` plus `tools/go.sum`. Any source change invalidates the cache
and forces a rebuild.
