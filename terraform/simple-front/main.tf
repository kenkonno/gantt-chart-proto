terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.21.0"
    }
  }
  required_version = ">= 1.1.8"
}

provider "aws" {
  access_key = var.accessKey
  secret_key = var.secretKey
  region     = "ap-northeast-1"
  default_tags {
    tags = {
      env     = var.env
      service = var.serviceName
    }
  }
}

# バケット作成
resource "aws_s3_bucket" "front_bucket" {
  bucket = "${var.env}-${var.serviceName}-bucket"
}

# バケットポリシー（cloudfront - s3の許可）
resource "aws_s3_bucket_policy" "front_s3_bucket_policy" {
  bucket = aws_s3_bucket.front_bucket.bucket
  policy = jsonencode({
    "Version" : "2008-10-17",
    "Id" : "PolicyForCloudFrontPrivateContent",
    "Statement" : [
      {
        "Sid" : "AllowCloudFrontServicePrincipal",
        "Effect" : "Allow",
        "Principal" : {
          "Service" : "cloudfront.amazonaws.com"
        },
        "Action" : "s3:GetObject",
        "Resource" : "${aws_s3_bucket.front_bucket.arn}/*",
        "Condition" : {
          "StringEquals" : {
            "AWS:SourceArn" : aws_cloudfront_distribution.front_cloudfront.arn
          }
        }
      }
    ]
  })
}

# cloudfrontの推奨設定の設定。
resource "aws_cloudfront_origin_access_control" "front_origin_access_control" {
  name                              = aws_s3_bucket.front_bucket.bucket_regional_domain_name
  description                       = ""
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# cloudfront作成
resource "aws_cloudfront_distribution" "front_cloudfront" {
  origin {
    domain_name         = aws_s3_bucket.front_bucket.bucket_regional_domain_name
    origin_id           = aws_s3_bucket.front_bucket.bucket_regional_domain_name
    connection_attempts = 3
    connection_timeout  = 10
    origin_path         = ""
    origin_access_control_id = aws_cloudfront_origin_access_control.front_origin_access_control.id
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations        = []
    }
  }
  enabled             = true
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = aws_s3_bucket.front_bucket.bucket_regional_domain_name
    viewer_protocol_policy = "allow-all"
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
    minimum_protocol_version       = "TLSv1"
  }

  // TODO: 次回構築時に確認 エラーページのリライト
  custom_error_response {
    error_code = 403
    error_caching_min_ttl = 10
    response_code = 200
    response_page_path = "/index.html"
  }

  #  EC2との連携は別のterraformファイルにして後から更新する形にしたほうがよさそう
  #  ordered_cache_behavior {
  #    allowed_methods        = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
  #    cached_methods         = ["GET", "HEAD"]
  #    compress               = false,
  #    default_ttl            = 0,
  #    max_ttl                = 0,
  #    min_ttl                = 0,
  #    path_pattern           = "/api/*",
  #    smooth_streaming       = false,
  #    target_origin_id       = "ec2-3-112-66-219.ap-northeast-1.compute.amazonaws.com", // TODO: Ec2 origin id
  #    viewer_protocol_policy = "redirect-to-https"
  #  }

}

output "url" {
  value = aws_cloudfront_distribution.front_cloudfront.domain_name
}

output "s3" {
  value = "https://s3.console.aws.amazon.com/s3/buckets/${aws_s3_bucket.front_bucket.bucket}"
}