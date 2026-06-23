terraform {
  backend "s3" {
    # Backend configuration should be provided via backend-config or partial config
    # terraform init -backend-config="tfvars/backend-dev.hcl"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.81.0, != 6.1.0"
    }
  }
}

provider "aws" {
  allowed_account_ids = [var.aws_account_id]
  region              = var.aws_region

  default_tags {
    tags = {
      Environment = var.env
      Project     = "snooker"
    }
  }
}

output "github_env" {
  description = "Environment variables for GitHub"
  value = {
    WEB_BUCKET      = aws_s3_bucket.static_web.id
    DISTRIBUTION_ID = aws_cloudfront_distribution.snooker.id
  }
}
