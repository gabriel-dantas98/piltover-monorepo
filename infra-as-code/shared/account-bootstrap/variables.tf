variable "aws_region" {
  description = "AWS region used by both the state bucket and the DynamoDB lock table."
  type        = string
  default     = "us-east-1"
}

variable "bucket_name" {
  description = "Name of the S3 bucket that stores OpenTofu state for the account."
  type        = string
  default     = "piltover-tofu-state"
}

variable "lock_table_name" {
  description = "Name of the DynamoDB table used for OpenTofu state locking."
  type        = string
  default     = "piltover-tofu-locks"
}

variable "tags" {
  description = "Common tags applied to every resource."
  type        = map(string)
  default = {
    Project   = "piltover"
    ManagedBy = "opentofu"
    Stack     = "account-bootstrap"
  }
}
