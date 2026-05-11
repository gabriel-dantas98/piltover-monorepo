# setup-tofu-aws-oidc

Composite action that assumes an AWS IAM role via GitHub OIDC, installs OpenTofu,
and exposes both on the runner.

## Usage

```yaml
permissions:
  id-token: write
  contents: read

steps:
  - uses: actions/checkout@v4
  - uses: ./ci-cd-actions/setup-tofu-aws-oidc
    with:
      role_arn: arn:aws:iam::123456789012:role/github-actions-piltover
      aws_region: us-east-1   # default
      tofu_version: "1.8.0"   # default
  - run: piltover tf plan infra-as-code/shared/github-oidc-provider
```

The `role_arn` is the `role_arn` output of
[`infra-as-code/shared/github-oidc-provider`](../../infra-as-code/shared/github-oidc-provider/README.md).

## Required permissions

The caller workflow must declare `id-token: write` at the workflow or job level —
that's how `aws-actions/configure-aws-credentials@v4` obtains the OIDC token.
