# s3-static-site

A reusable OpenTofu module that provisions a CloudFront-fronted static site backed
by S3, with an ACM certificate, optional Route53 alias, and SigV4-signed access via
Origin Access Control.

This is a **child module** — it does not configure a backend or providers. Consumer
root modules provide both (notably the `aws.us_east_1` alias for ACM).

## Usage

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
      configuration_aliases = [aws.us_east_1]
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

module "site" {
  source = "git::https://github.com/gabriel-dantas98/piltover-monorepo.git//infra-as-code/modules/s3-static-site?ref=infra-modules/s3-static-site/v0.1.0"

  providers = {
    aws.us_east_1 = aws.us_east_1
  }

  domain_name = "demo.example.com"
  zone_id     = "Z0123456789ABCDEFGHIJ"
  tags = {
    Project = "demo"
  }
}
```

## Inputs

| Name | Type | Default | Purpose |
|---|---|---|---|
| `domain_name` | string | — | FQDN to serve from (e.g. `docs.example.com`). |
| `zone_id` | string | — | Route53 hosted zone ID. |
| `bucket_name` | string | `""` (= `domain_name`) | Override the S3 bucket name. |
| `create_route53_record` | bool | `true` | Toggle the A-record alias. |
| `tags` | map(string) | `{}` | Tags on every resource. |

## Outputs

| Name | Purpose |
|---|---|
| `bucket_name` | S3 bucket name. |
| `distribution_id` | CloudFront distribution ID. |
| `distribution_domain` | Default CloudFront domain. |
| `domain_name` | Echo of the input. |

## Notes

- ACM cert is created in `us-east-1` regardless of the bucket region — that's a
  CloudFront requirement.
- The S3 bucket is private; CloudFront's OAC is the only allowed reader.
- For local docs hosting (this repo's `apps/docs`), GitHub Pages is used instead.
  This module is the template for any future static site that needs richer edge
  behaviour (custom error pages, signed URLs, geo restrictions, etc.).
