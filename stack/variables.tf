variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
}

variable "aws_region" {
  description = "AWS Region"
  type        = string
}

variable "binary_bucket" {
  description = "S3 bucket with compiled lambda binaries"
  type        = string
}

variable "cloudfront" {
  type = object({
    domain          = string
    certificate_arn = string
  })
}

variable "cognito" {
  type = object({
    domain_prefix    = string
    google_client_id = string
    callback_urls    = list(string)
    logout_urls      = list(string)
  })
}

variable "duration" {
  type = object({
    start_date = string
    end_date   = string
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

variable "permissions_boundary_arn" {
  type = string
}

variable "manifest_file" {
  type    = string
  default = "default.json"
}

variable "zone_domain" {
  description = "Route 53 Hosted Zone ID"
  type        = string
}
