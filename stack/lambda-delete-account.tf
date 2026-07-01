module "lambda_delete_account" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "snooker"
  aws_env                  = var.env
  name                     = "delete-account"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["delete-account"]
  kms_key_arn              = aws_kms_key.main.arn

  environment = {
    USER_TABLE_NAME = aws_dynamodb_table.users.id
    USER_POOL_ID    = aws_cognito_user_pool.main.id
  }
}

module "apig_delete_account" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api        = aws_api_gateway_rest_api.main
  http_method     = "POST"
  path            = "deleteAccount"
  lambda          = module.lambda_delete_account
  parent_id       = aws_api_gateway_resource.api.id
  authorizer_id   = aws_api_gateway_authorizer.main.id
  authorizer_type = "COGNITO_USER_POOLS"
}

module "iam_kms_delete_account" {
  source  = "github.com/ockendenjo/tfmods//iam-kms"
  role_id = module.lambda_delete_account.role_id
  kms_arn = aws_kms_key.main.arn
}

module "iam_dynamodb_delete_account" {
  source  = "github.com/ockendenjo/tfmods.git//iam-dynamodb"
  role_id = module.lambda_delete_account.role_id
  dynamo_table_arns = [
    aws_dynamodb_table.users.arn,
  ]
  allow_index_use = true
}

resource "aws_iam_role_policy" "iam_cognito_delete_account" {
  role = module.lambda_delete_account.role_id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["cognito-idp:AdminDeleteUser"]
      Resource = aws_cognito_user_pool.main.arn
    }]
  })
}
