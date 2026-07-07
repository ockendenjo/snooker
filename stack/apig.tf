resource "aws_api_gateway_rest_api" "main" {
  name        = "snooker-api-${var.env}"
  description = "Snooker (${var.env})"

  disable_execute_api_endpoint = false

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "main" {
  rest_api_id = aws_api_gateway_rest_api.main.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_authorizer.main,
      aws_api_gateway_resource.api,

      module.apig_delete_account,
      module.apig_get_session_data,
      module.apig_log_drink,
      module.apig_set_displayname,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_authorizer" "main" {
  name                             = "snooker-cognito-${var.env}"
  type                             = "COGNITO_USER_POOLS"
  provider_arns                    = [aws_cognito_user_pool.main.arn]
  rest_api_id                      = aws_api_gateway_rest_api.main.id
  identity_source                  = "method.request.header.Authorization"
  authorizer_result_ttl_in_seconds = 1
}

resource "aws_api_gateway_stage" "stg" {
  deployment_id = aws_api_gateway_deployment.main.id
  rest_api_id   = aws_api_gateway_rest_api.main.id
  stage_name    = "stg"
}

resource "aws_api_gateway_resource" "api" {
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "api"
  rest_api_id = aws_api_gateway_rest_api.main.id
}
