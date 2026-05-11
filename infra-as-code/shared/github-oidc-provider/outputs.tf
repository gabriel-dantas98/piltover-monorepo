output "role_arn" {
  description = "ARN of the IAM role GitHub Actions can assume."
  value       = aws_iam_role.github_actions.arn
}

output "provider_arn" {
  description = "ARN of the IAM OIDC provider."
  value       = aws_iam_openid_connect_provider.github.arn
}
