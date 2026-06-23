variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
}

variable "aws_region" {
  description = "AWS Region"
  type        = string
}

variable "cloudfront" {
  type = object({
    domain          = string
    certificate_arn = string
  })
}

variable "env" {
  description = "Environment name (dev or pro)"
  type        = string

  validation {
    condition     = contains(["dev", "pro"], var.env)
    error_message = "Environment must be either 'dev' or 'pro'."
  }
}

variable "zone_domain" {
  description = "Route 53 Hosted Zone ID"
  type        = string
}
