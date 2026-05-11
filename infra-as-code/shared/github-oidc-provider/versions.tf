terraform {
  required_version = ">= 1.8"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
  # Backend is configured at init time via -backend-config flags pointing at the
  # outputs of account-bootstrap. See README.
  backend "s3" {}
}
