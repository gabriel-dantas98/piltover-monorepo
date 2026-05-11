variable "domain_name" {
  description = "Fully qualified domain name the site will be served at (e.g. docs.example.com)."
  type        = string
}

variable "zone_id" {
  description = "Route53 hosted zone ID where the DNS record + ACM validation records live."
  type        = string
}

variable "bucket_name" {
  description = "Override the S3 bucket name. Defaults to the domain name."
  type        = string
  default     = ""
}

variable "create_route53_record" {
  description = "Whether to create the A-record alias pointing at CloudFront. Set to false if you manage DNS elsewhere."
  type        = bool
  default     = true
}

variable "tags" {
  description = "Tags applied to every resource."
  type        = map(string)
  default     = {}
}
