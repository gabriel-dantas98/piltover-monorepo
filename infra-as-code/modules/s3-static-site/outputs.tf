output "bucket_name" {
  description = "Name of the S3 bucket."
  value       = aws_s3_bucket.site.bucket
}

output "distribution_id" {
  description = "ID of the CloudFront distribution."
  value       = aws_cloudfront_distribution.site.id
}

output "distribution_domain" {
  description = "Default CloudFront domain (e.g. d111111abcdef8.cloudfront.net)."
  value       = aws_cloudfront_distribution.site.domain_name
}

output "domain_name" {
  description = "Public domain name (echo of the input)."
  value       = var.domain_name
}
