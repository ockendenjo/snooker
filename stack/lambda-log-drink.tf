module "lambda_log_drink" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  project_name             = "snooker"
  aws_env                  = var.env
  name                     = "log-drink"
  permissions_boundary_arn = var.permissions_boundary_arn
  s3_bucket                = var.binary_bucket
  s3_object_key            = local.manifest["log-drink"]
  kms_key_arn              = aws_kms_key.main.arn

  environment = {
    DRINKS_TABLE_NAME = aws_dynamodb_table.drinks_ew1.id
    USER_TABLE_NAME   = aws_dynamodb_table.users.id
    START_DATE        = var.duration.start_date
    END_DATE          = var.duration.end_date
  }
}

module "apig_log_drink" {
  source = "github.com/ockendenjo/tfmods//apig-endpoint"

  rest_api        = aws_api_gateway_rest_api.main
  http_method     = "POST"
  path            = "logDrink"
  lambda          = module.lambda_log_drink
  parent_id       = aws_api_gateway_resource.api.id
  authorizer_id   = aws_api_gateway_authorizer.main.id
  authorizer_type = "COGNITO_USER_POOLS"
}

module "iam_kms_log_drink" {
  source        = "github.com/ockendenjo/tfmods//iam-kms"
  role_id       = module.lambda_log_drink.role_id
  kms_arn       = aws_kms_key.main.arn
  allow_encrypt = true
}

module "iam_dynamodb_log_drink" {
  source  = "github.com/ockendenjo/tfmods.git//iam-dynamodb"
  role_id = module.lambda_log_drink.role_id
  dynamo_table_arns = [
    aws_dynamodb_table.drinks_ew1.arn,
    aws_dynamodb_table.users.arn,
  ]
}
