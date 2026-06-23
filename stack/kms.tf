resource "aws_kms_key" "main" {
  description             = "snooker-${var.env}"
  deletion_window_in_days = 7
  enable_key_rotation     = false
  multi_region            = false

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "EnableKeyAdministration"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${var.aws_account_id}:root"
        }
        Action = [
          "kms:Create*",
          "kms:Describe*",
          "kms:Enable*",
          "kms:List*",
          "kms:Put*",
          "kms:Update*",
          "kms:Revoke*",
          "kms:Disable*",
          "kms:Decrypt",
          "kms:Get*",
          "kms:Delete*",
          "kms:TagResource",
          "kms:UntagResource",
          "kms:ScheduleKeyDeletion",
          "kms:CancelKeyDeletion",
          "kms:ReplicateKey",
        ]
        Resource = "*"
      },
      # CloudFront only uses the primary key - no equivalent statement in the replica key policy
      {
        Sid    = "AllowUseFromCloudFront"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action = [
          "kms:Decrypt",
          "kms:GenerateDataKey*",
        ]
        Resource = "*"
        Condition = {
          StringEquals = {
            "aws:SourceArn" = aws_cloudfront_distribution.snooker.arn
          }
        }
      },
      {
        Sid    = "AllowUseFromSnookerRoles"
        Effect = "Allow"
        Principal = {
          AWS = "*"
        }
        Action = [
          "kms:Decrypt",
          "kms:DescribeKey",
          "kms:Encrypt",
          "kms:GenerateDataKey*",
          "kms:ReEncrypt*",
        ]
        Resource = "*"
        Condition = {
          StringLike = {
            "aws:PrincipalArn" = "arn:aws:iam::${var.aws_account_id}:role/snooker-*"
          }
        }
      },
      {
        Sid    = "AllowCloudWatchLogs"
        Effect = "Allow"
        Principal = {
          Service = "logs.eu-west-1.amazonaws.com"
        }
        Action = [
          "kms:Encrypt*",
          "kms:Decrypt*",
          "kms:ReEncrypt*",
          "kms:GenerateDataKey*",
          "kms:Describe*",
        ]
        Resource = "*"
        Condition = {
          ArnLike = {
            "kms:EncryptionContext:aws:logs:arn" = "arn:aws:logs:eu-west-1:${var.aws_account_id}:*"
          }
        }
      }
    ]
  })
}

resource "aws_kms_alias" "main" {
  name          = "alias/snooker-${var.env}"
  target_key_id = aws_kms_key.main.key_id
}
