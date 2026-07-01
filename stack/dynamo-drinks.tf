resource "aws_dynamodb_table" "drinks_ew1" {
  name         = "snooker-drinks-${var.env}"
  billing_mode = "PAY_PER_REQUEST"
  table_class  = "STANDARD"

  hash_key  = "user_id"
  range_key = "tstamp"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "tstamp"
    type = "S"
  }

  attribute {
    name = "unknown_beer"
    type = "N"
  }

  deletion_protection_enabled = false

  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.main.arn
  }

  global_secondary_index {
    name            = "unknown-beer"
    projection_type = "KEYS_ONLY"

    key_schema {
      attribute_name = "unknown_beer"
      key_type       = "HASH"
    }
  }
}
