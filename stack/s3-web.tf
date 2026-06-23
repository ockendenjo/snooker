resource "aws_s3_bucket" "static_web" {
  bucket_prefix = "snooker-${var.env}-web-"
  force_destroy = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "static_web" {
  bucket = aws_s3_bucket.static_web.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.main.arn
    }
    bucket_key_enabled = true
  }
}
