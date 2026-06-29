data "aws_s3_object" "manifest_file" {
  bucket = var.binary_bucket
  key    = "lambda_manifests/${var.manifest_file}"
}

locals {
  manifest = jsondecode(data.aws_s3_object.manifest_file.body)
}
