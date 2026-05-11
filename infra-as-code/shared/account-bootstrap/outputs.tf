output "state_bucket" {
  description = "Name of the S3 bucket that stores OpenTofu state."
  value       = aws_s3_bucket.state.bucket
}

output "state_bucket_region" {
  description = "Region of the state bucket."
  value       = var.aws_region
}

output "lock_table" {
  description = "Name of the DynamoDB lock table."
  value       = aws_dynamodb_table.locks.name
}
