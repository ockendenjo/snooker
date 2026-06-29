# Cognito User Pool
resource "aws_cognito_user_pool" "main" {
  name = "snooker-${var.env}"

  auto_verified_attributes = ["email"]

  username_attributes = ["email"]
  username_configuration {
    case_sensitive = false
  }

  # Enable self-registration - users can sign up with email/password or via Google
  admin_create_user_config {
    allow_admin_create_user_only = false
  }

  schema {
    name                = "email"
    attribute_data_type = "String"
    mutable             = true
    required            = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }
}

resource "aws_cognito_user_pool_domain" "main" {
  domain                = var.cognito.domain_prefix
  user_pool_id          = aws_cognito_user_pool.main.id
  managed_login_version = 2
}

resource "aws_cognito_managed_login_branding" "client" {
  client_id    = aws_cognito_user_pool_client.main.id
  user_pool_id = aws_cognito_user_pool.main.id

  use_cognito_provided_values = true
}

data "aws_ssm_parameter" "google_client_secret" {
  name = "/snooker/${var.env}/google/client_secret"
}

resource "aws_cognito_identity_provider" "google" {
  user_pool_id  = aws_cognito_user_pool.main.id
  provider_name = "Google"
  provider_type = "Google"

  provider_details = {
    authorize_scopes              = "openid email"
    attributes_url                = "https://people.googleapis.com/v1/people/me?personFields="
    token_url                     = "https://www.googleapis.com/oauth2/v4/token"
    token_request_method          = "POST"
    authorize_url                 = "https://accounts.google.com/o/oauth2/v2/auth"
    attributes_url_add_attributes = "true"
    client_id                     = var.cognito.google_client_id
    client_secret                 = data.aws_ssm_parameter.google_client_secret.value
    oidc_issuer                   = "https://accounts.google.com"
  }

  attribute_mapping = {
    email = "email"
  }
}

resource "aws_cognito_user_pool_client" "main" {
  name         = "snooker-client-${var.env}"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret = false

  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["openid", "email", "profile"]

  callback_urls = var.cognito.callback_urls
  logout_urls   = var.cognito.logout_urls

  supported_identity_providers = ["COGNITO", "Google"]

  explicit_auth_flows = [
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_USER_PASSWORD_AUTH",
  ]

  refresh_token_validity = 30
  token_validity_units {
    refresh_token = "days"
  }

  prevent_user_existence_errors = "ENABLED"

  depends_on = [aws_cognito_identity_provider.google]
}

output "cognito" {
  description = "Cognito settings for UI"
  value = {
    domain   = "${var.cognito.domain_prefix}.auth.${var.aws_region}.amazoncognito.com"
    clientId = aws_cognito_user_pool_client.main.id
  }
}
