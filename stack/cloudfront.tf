resource "aws_cloudfront_function" "beer_rewrite" {
  name    = "snooker-beer-rewrite-${var.env}"
  runtime = "cloudfront-js-2.0"
  publish = true
  code    = <<-EOT
    async function handler(event) {
      const request = event.request;
      request.uri = request.uri.replace(/^\/beer\//, '/');
      return request;
    }
  EOT
}

resource "aws_cloudfront_origin_access_control" "s3" {
  name                              = "snooker-${var.env}"
  signing_behavior                  = "always"
  origin_access_control_origin_type = "s3"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_response_headers_policy" "robots" {
  name = "snooker-robots-${var.env}"

  custom_headers_config {
    items {
      header   = "X-Robots-Tag"
      value    = "noindex, nofollow, noarchive, nosnippet"
      override = true
    }
  }
}

resource "aws_cloudfront_response_headers_policy" "no_cache" {
  name = "snooker-no-cache-${var.env}"

  custom_headers_config {
    items {
      header   = "Cache-Control"
      value    = "no-cache, max-age=0, must-revalidate"
      override = true
    }
  }
}

resource "aws_cloudfront_response_headers_policy" "beer_json" {
  name = "snooker-beer-json-${var.env}"

  custom_headers_config {
    items {
      header   = "Cache-Control"
      value    = "max-age=31536000, immutable"
      override = true
    }
  }
}

resource "aws_cloudfront_response_headers_policy" "index_html" {
  name = "snooker-index-html-${var.env}"

  custom_headers_config {
    items {
      header   = "X-Robots-Tag"
      value    = "noindex, nofollow, noarchive, nosnippet"
      override = true
    }
    items {
      header   = "Cache-Control"
      value    = "no-cache, max-age=0, must-revalidate"
      override = true
    }
  }
}

resource "aws_cloudfront_distribution" "snooker" {
  comment             = "snooker (${var.env})"
  enabled             = true
  default_root_object = "index.html"
  price_class         = "PriceClass_100"
  http_version        = "http2and3"
  aliases             = [var.cloudfront.domain]

  origin {
    domain_name              = aws_s3_bucket.static_web.bucket_regional_domain_name
    origin_id                = "s3-static"
    origin_access_control_id = aws_cloudfront_origin_access_control.s3.id
  }

  origin {
    domain_name              = "beerdb-data-pro-20260609135102107600000001.s3.eu-west-1.amazonaws.com"
    origin_id                = "beerdb"
    origin_access_control_id = aws_cloudfront_origin_access_control.s3.id
  }

  origin {
    domain_name              = "pub-db-data-pro-20260619095306867700000001.s3.eu-west-1.amazonaws.com"
    origin_id                = "pub-db"
    origin_access_control_id = aws_cloudfront_origin_access_control.s3.id
  }

  origin {
    domain_name = "${aws_api_gateway_rest_api.main.id}.execute-api.${var.aws_region}.amazonaws.com"
    origin_path = "/stg"
    origin_id   = "apig"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    target_origin_id       = "s3-static"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true

    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_optimized.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.robots.id

    allowed_methods = ["GET", "HEAD"]
    cached_methods  = ["GET", "HEAD"]
  }

  ordered_cache_behavior {
    path_pattern               = "index.html"
    target_origin_id           = "s3-static"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.index_html.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]
  }

  ordered_cache_behavior {
    path_pattern               = "pubs.json"
    target_origin_id           = "pub-db"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.no_cache.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]
  }

  ordered_cache_behavior {
    path_pattern               = "ngsw.json"
    target_origin_id           = "s3-static"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.no_cache.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]
  }

  ordered_cache_behavior {
    path_pattern               = "ngsw-worker.js"
    target_origin_id           = "s3-static"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.no_cache.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]
  }

  ordered_cache_behavior {
    path_pattern               = "beer/index.json"
    target_origin_id           = "beerdb"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.no_cache.id
    origin_request_policy_id   = data.aws_cloudfront_origin_request_policy.all_viewer_except_host.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.beer_rewrite.arn
    }
  }

  ordered_cache_behavior {
    path_pattern               = "beer/*.json"
    target_origin_id           = "beerdb"
    viewer_protocol_policy     = "redirect-to-https"
    compress                   = true
    cache_policy_id            = data.aws_cloudfront_cache_policy.caching_optimized.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.beer_json.id
    origin_request_policy_id   = data.aws_cloudfront_origin_request_policy.all_viewer_except_host.id
    allowed_methods            = ["GET", "HEAD"]
    cached_methods             = ["GET", "HEAD"]

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.beer_rewrite.arn
    }
  }

  # ordered_cache_behavior {
  #   path_pattern             = "api/*"
  #   target_origin_id         = "apig"
  #   viewer_protocol_policy   = "redirect-to-https"
  #   compress                 = true
  #   cache_policy_id          = data.aws_cloudfront_cache_policy.caching_disabled.id
  #   origin_request_policy_id = data.aws_cloudfront_origin_request_policy.all_viewer_except_host.id
  #   allowed_methods          = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
  #   cached_methods           = ["GET", "HEAD"]
  # }

  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 404
    response_code      = 200
    response_page_path = "/index.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = var.cloudfront.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}

module "s3_policy_cloudfront" {
  source = "github.com/ockendenjo/tfmods//s3-policy-cloudfront"
  bucket = aws_s3_bucket.static_web
  cloudfront_arns = [
    aws_cloudfront_distribution.snooker.arn,
  ]
}

resource "aws_route53_record" "snooker" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = var.cloudfront.domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.snooker.domain_name
    zone_id                = aws_cloudfront_distribution.snooker.hosted_zone_id
    evaluate_target_health = false
  }
}

data "aws_cloudfront_cache_policy" "caching_optimized" {
  name = "Managed-CachingOptimized"
}

data "aws_cloudfront_cache_policy" "caching_disabled" {
  name = "Managed-CachingDisabled"
}

data "aws_cloudfront_origin_request_policy" "all_viewer_except_host" {
  name = "Managed-AllViewerExceptHostHeader"
}

data "aws_route53_zone" "primary" {
  name         = var.zone_domain
  private_zone = false
}
