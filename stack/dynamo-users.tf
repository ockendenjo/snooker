resource "aws_dynamodb_table" "users" {
  name         = "snooker-users-v2-${var.env}"
  billing_mode = "PAY_PER_REQUEST"
  table_class  = "STANDARD"

  hash_key = "email"

  attribute {
    name = "email"
    type = "S"
  }

  attribute {
    name = "id"
    type = "S"
  }

  global_secondary_index {
    name            = "id-index"
    projection_type = "ALL"

    key_schema {
      attribute_name = "id"
      key_type       = "HASH"
    }
  }

  deletion_protection_enabled = false

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.main.arn
  }
}
