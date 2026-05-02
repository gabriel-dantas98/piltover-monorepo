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
