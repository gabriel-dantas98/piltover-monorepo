# account-bootstrap

Creates the AWS resources every other OpenTofu root module needs:

- An S3 bucket for remote state (versioned, encrypted, public-access blocked).
- A DynamoDB table for state locking.

This stack is **applied manually exactly once per AWS account** with the local state
backend declared in `versions.tf`. After it runs, every other root module references
the outputs via an S3 backend pointing at this bucket.

## Apply procedure

```bash
# 1) Configure local AWS credentials (e.g. via aws sso login) for the target account.
aws sts get-caller-identity   # sanity check

# 2) Initialise + apply.
piltover tf init  infra-as-code/shared/account-bootstrap
piltover tf apply infra-as-code/shared/account-bootstrap
```

That's it. Keep the resulting `terraform.tfstate` file local — do NOT migrate it to
S3 (the bucket it would migrate into IS one of the resources this stack creates).

## Outputs

After `apply`, capture the outputs:

```bash
tofu -chdir=infra-as-code/shared/account-bootstrap output
```

You'll need `state_bucket`, `state_bucket_region`, and `lock_table` to configure the
backend of every subsequent root module:

```hcl
terraform {
  backend "s3" {
    bucket         = "<state_bucket output>"
    key            = "<unique-key-per-stack>/terraform.tfstate"
    region         = "<state_bucket_region output>"
    dynamodb_table = "<lock_table output>"
    encrypt        = true
  }
}
```

## Variables

| Name | Default | Purpose |
|---|---|---|
| `aws_region` | `us-east-1` | Region for both resources. |
| `bucket_name` | `piltover-tofu-state` | S3 bucket name (must be globally unique — override if taken). |
| `lock_table_name` | `piltover-tofu-locks` | DynamoDB table name. |
| `tags` | (sensible defaults) | Common tags. |
