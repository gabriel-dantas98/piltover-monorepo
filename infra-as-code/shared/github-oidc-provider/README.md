# github-oidc-provider

Creates the AWS IAM resources GitHub Actions needs to assume an AWS role via OIDC
(no static `AWS_ACCESS_KEY_ID` ever stored in GitHub Secrets).

## What it ships

| Resource | Why |
|---|---|
| `aws_iam_openid_connect_provider.github` | Tells AWS to trust `token.actions.githubusercontent.com`. |
| `aws_iam_role.github_actions` | The role GitHub Actions assumes. Has **no** attached policies by default — attach least-privilege policies per workflow. |

## Apply procedure

Apply this **after** [`account-bootstrap`](../account-bootstrap/README.md). You'll need
the `state_bucket`, `state_bucket_region`, and `lock_table` outputs from there.

```bash
# Pull the outputs:
cd infra-as-code/shared/account-bootstrap
STATE_BUCKET=$(tofu output -raw state_bucket)
STATE_REGION=$(tofu output -raw state_bucket_region)
LOCK_TABLE=$(tofu output -raw lock_table)
cd ../../..

# Initialise this stack against the S3 backend:
piltover tf init infra-as-code/shared/github-oidc-provider -- \
  -backend-config="bucket=${STATE_BUCKET}" \
  -backend-config="key=shared/github-oidc-provider/terraform.tfstate" \
  -backend-config="region=${STATE_REGION}" \
  -backend-config="dynamodb_table=${LOCK_TABLE}" \
  -backend-config="encrypt=true"

# Apply.
piltover tf apply infra-as-code/shared/github-oidc-provider
```

## Use in a workflow

After apply, grab the `role_arn` output. In any GitHub Actions workflow that needs
AWS access, add:

```yaml
permissions:
  id-token: write
  contents: read

steps:
  - uses: ./ci-cd-actions/setup-tofu-aws-oidc
    with:
      role_arn: <the role_arn output>
      aws_region: us-east-1
```

## Variables

| Name | Default | Purpose |
|---|---|---|
| `aws_region` | `us-east-1` | Region for IAM (note: IAM is global but providers need a region). |
| `github_org` | `gabriel-dantas98` | Owner of the repo. |
| `github_repo` | `piltover-monorepo` | Repository name. |
| `role_name` | `github-actions-piltover` | IAM role name. |
| `subject_patterns` | `main` branch + pull requests | OIDC sub-claim patterns that may assume the role. |

Tighten `subject_patterns` (e.g. drop `pull_request`) for stricter prod environments.
