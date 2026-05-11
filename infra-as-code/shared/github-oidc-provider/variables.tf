variable "aws_region" {
  description = "AWS region for the IAM resources."
  type        = string
  default     = "us-east-1"
}

variable "github_org" {
  description = "GitHub organisation or user that owns the repo."
  type        = string
  default     = "gabriel-dantas98"
}

variable "github_repo" {
  description = "GitHub repository name."
  type        = string
  default     = "piltover-monorepo"
}

variable "role_name" {
  description = "Name of the IAM role assumable from GitHub Actions."
  type        = string
  default     = "github-actions-piltover"
}

variable "subject_patterns" {
  description = "List of OIDC subject patterns (sub claim) that may assume the role. Matches the standard 'repo:org/repo:ref:refs/heads/<branch>' format."
  type        = list(string)
  default = [
    "repo:gabriel-dantas98/piltover-monorepo:ref:refs/heads/main",
    "repo:gabriel-dantas98/piltover-monorepo:pull_request",
  ]
}

variable "tags" {
  description = "Common tags applied to every resource."
  type        = map(string)
  default = {
    Project   = "piltover"
    ManagedBy = "opentofu"
    Stack     = "github-oidc-provider"
  }
}
