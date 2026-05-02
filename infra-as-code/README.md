# infra-as-code/

OpenTofu modules (`modules/`, reusable child modules) and apply-once bootstrap
stacks (`shared/`, e.g. github-oidc-provider, account-bootstrap).

Apps consume modules from `apps/<name>/infra/` (root module per app).

State backend: S3 + DynamoDB lock, configured automatically by `piltover tf`.
