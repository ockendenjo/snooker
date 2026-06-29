module "lambda_get_session_data" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "snooker"
  aws_env                  = var.env
  name                     = "get-session-data"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["get-session-data"]
  kms_key_arn              = aws_kms_key.main.arn

  environment = {
    USER_TABLE_NAME = aws_dynamodb_table.users.id
  }
}

module "apig_get_session_data" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api        = aws_api_gateway_rest_api.main
  http_method     = "POST"
  path            = "getSessionData"
  lambda          = module.lambda_get_session_data
  parent_id       = aws_api_gateway_resource.api.id
  authorizer_id   = aws_api_gateway_authorizer.main.id
  authorizer_type = "COGNITO_USER_POOLS"
}

module "iam_kms_get_session_data" {
  source  = "github.com/ockendenjo/tfmods//iam-kms"
  role_id = module.lambda_get_session_data.role_id
  kms_arn = aws_kms_key.main.arn
}

module "iam_dynamodb_get_session_data" {
  source  = "github.com/ockendenjo/tfmods.git//iam-dynamodb"
  role_id = module.lambda_get_session_data.role_id
  dynamo_table_arns = [
    aws_dynamodb_table.users.arn,
  ]
}
