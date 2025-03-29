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
  region     = var.region
  default_tags {
    tags = {
      env     = var.env
      service = var.serviceName
    }
  }
}

####################################################
#                      Route53?
####################################################


####################################################
#                  Service discovery
####################################################
# Service Connect 版に一応書き換えたこれで上手くいけばいいが・・・

####################################################
#                  API Gateway
####################################################
resource "aws_service_discovery_http_namespace" "service_connect_namespace" {
  name        = "${var.env}-${var.serviceName}-namespace"
  description = "HTTP Namespace for Service Connect"
}

# API GatewayがCloud Mapにアクセスするためのロール
resource "aws_apigatewayv2_api" "api_gateway" {
  api_key_selection_expression = "$request.header.x-api-key"
  disable_execute_api_endpoint = false
  name                         = "${var.env}-${var.serviceName}-api"
  protocol_type                = "HTTP"
  route_selection_expression   = "$request.method $request.path"
}

resource "aws_apigatewayv2_route" "api_route" {
  api_id               = aws_apigatewayv2_api.api_gateway.id
  api_key_required     = false
  authorization_scopes = []
  authorization_type   = "NONE"
  request_models       = {}
  route_key            = "ANY /{proxy+}"
  target               = "integrations/${aws_apigatewayv2_integration.api_integration.id}"
}
resource "aws_apigatewayv2_vpc_link" "api_vpc_link" {
  name               = "${var.env}-${var.serviceName}-apigateway-for-ecs-vpc-link"
  security_group_ids = var.securityGroupIds
  subnet_ids         = var.privateSubnetIds
}
resource "aws_apigatewayv2_integration" "api_integration" {
  api_id                 = aws_apigatewayv2_api.api_gateway.id
  connection_id          = aws_apigatewayv2_vpc_link.api_vpc_link.id
  connection_type        = "VPC_LINK"
  integration_method     = "ANY"
  integration_type       = "HTTP_PROXY"

  # CloudMapサービスのARNを使用する apiはECSのdiscovery name に合わせている
#  integration_uri        = "http://api.${aws_service_discovery_http_namespace.service_connect_namespace.name}:80"
  integration_uri        = aws_service_discovery_http_namespace.service_connect_namespace.arn

  payload_format_version = "1.0"
  timeout_milliseconds   = 30000
}



// TODO: importしていない
resource "aws_cloudwatch_log_group" "api_gateway_log_group" {
  name = "/aws/apigateway/${var.env}-${var.serviceName}-api_gateway_log_group"
}

resource "aws_apigatewayv2_stage" "api_stage" {
  api_id      = aws_apigatewayv2_api.api_gateway.id
  auto_deploy = false
  name        = "api"
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_log_group.arn
    format          = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId"
  }
  default_route_settings {
    data_trace_enabled       = false
    detailed_metrics_enabled = true
    throttling_burst_limit   = 5000
    throttling_rate_limit    = 10000
  }
}
# APIGateway用のロギングロール
resource "aws_iam_role" "api_gateway_cloudwatch_role" {
  name = "${var.env}-${var.serviceName}-api-gateway-cloudwatch-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      }
    ]
  })
}

# ロールにCloudWatchへの書き込み権限を付与するポリシー
resource "aws_iam_role_policy" "api_gateway_cloudwatch_policy" {
  name = "${var.env}-${var.serviceName}-api-gateway-cloudwatch-policy"
  role = aws_iam_role.api_gateway_cloudwatch_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams",
          "logs:PutLogEvents",
          "logs:GetLogEvents",
          "logs:FilterLogEvents"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

# APIGatewayがこのロールを使用するよう設定するアカウント設定
resource "aws_api_gateway_account" "api_gateway_account" {
  cloudwatch_role_arn = aws_iam_role.api_gateway_cloudwatch_role.arn
}

// TODO: ECRは先に外だしで作ったほうがよくない？
#####################################################
##                       ECR
#####################################################
resource "aws_ecr_repository" "api_ecr_repository" {
  image_tag_mutability = "MUTABLE" // NOTE: docker image の上書き
  name                 = "${var.env}-${var.serviceName}/api"

  encryption_configuration {
    encryption_type = "AES256"
  }

  image_scanning_configuration {
    scan_on_push = false
  }
}
resource "aws_ecr_repository" "migration_ecr_repository" {
  image_tag_mutability = "IMMUTABLE" // NOTE: docker image の上書き
  name                 = "${var.env}-${var.serviceName}/migration"

  encryption_configuration {
    encryption_type = "AES256"
  }

  image_scanning_configuration {
    scan_on_push = false
  }
}
#####################################################
##                       ECS
#####################################################
resource "aws_iam_role" "ecsInstanceRole" {
  assume_role_policy = jsonencode(
    {
      Statement = [
        {
          Action    = "sts:AssumeRole"
          Effect    = "Allow"
          Principal = {
            Service = "ec2.amazonaws.com"
          }
          Sid = ""
        },{
          Action    = "sts:AssumeRole"
          Effect    = "Allow"
          Principal = {
            Service = "ecs-tasks.amazonaws.com"
          }
          Sid = ""
        },
      ]
      Version = "2008-10-17"
    }
  )
  force_detach_policies = false
  managed_policy_arns   = [
    "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role",
  ]
  max_session_duration = 3600
  name                 = "ecsInstanceRole"
  path                 = "/"
}

# instance role なるものを付ける必要がある
resource "aws_iam_instance_profile" "ecsInstanceRole" {
  name = aws_iam_role.ecsInstanceRole.name
  role = aws_iam_role.ecsInstanceRole.name
}

resource "aws_iam_role" "ecsTaskExecutionRole" {
  assume_role_policy = jsonencode(
    {
      Statement = [
        {
          Action    = "sts:AssumeRole"
          Effect    = "Allow"
          Principal = {
            Service = "ecs-tasks.amazonaws.com"
          }
          Sid = ""
        },
      ]
      Version = "2008-10-17"
    }
  )
  force_detach_policies = false
  managed_policy_arns   = [
    "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
    "arn:aws:iam::aws:policy/CloudWatchFullAccess",
    "arn:aws:iam::aws:policy/CloudWatchFullAccessV2",
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
  ]
  max_session_duration = 3600
  name                 = "ecsTaskExecutionRole"
  path                 = "/"
  tags                 = {}
  tags_all             = {}

  inline_policy {
    name   = "inline-ECS-Exec-Policy"
    policy = jsonencode(
      {
        Statement = [
          {
            Action = [
              "ssmmessages:CreateControlChannel",
              "ssmmessages:CreateDataChannel",
              "ssmmessages:OpenControlChannel",
              "ssmmessages:OpenDataChannel",
            ]
            Effect   = "Allow"
            Resource = "*"
          },
        ]
        Version = "2012-10-17"
      }
    )
  }
}
// TODO: apiってついてるけど起動テンプレートは１つになった。
resource "aws_launch_template" "api_launch_template" {
  default_version         = 1
  disable_api_stop        = false
  disable_api_termination = false
  image_id                = "ami-0fd4bb49e3134c823"
  instance_type           = "t2.micro"
  key_name                = var.sshKeyPairName
  vpc_security_group_ids  = var.securityGroupIds
  user_data = base64encode(<<DATA
#!/bin/bash
echo ECS_CLUSTER=${aws_ecs_cluster.api_ecs_cluster.name} >> /etc/ecs/ecs.config;
DATA
)
  iam_instance_profile {
    arn = aws_iam_instance_profile.ecsInstanceRole.arn
  }
}
# 初回実行で作成られる者らしく、おまじないみたいなもの
resource "aws_iam_service_linked_role" "autoscaling_group" {
  aws_service_name = "autoscaling.amazonaws.com" # こっちかも
}
resource "aws_autoscaling_group" "api_autoscaling_group" {
  // TODO: aws_autoscaling_group availability_zones 何がいいんだろう。デフォルト設定で２個だからそれでいいのかな。
  vpc_zone_identifier       = var.privateSubnetIds
  capacity_rebalance        = false
  default_cooldown          = 300
  default_instance_warmup   = 0
  desired_capacity          = 1
  enabled_metrics           = []
  load_balancers            = []
  max_instance_lifetime     = 0
  max_size                  = 2
  metrics_granularity       = "1Minute"
  min_size                  = 0
  name                      = "Infra-ECS-Cluster-${var.env}-${var.serviceName}-ecs-cluster"
  protect_from_scale_in     = true
  service_linked_role_arn   = aws_iam_service_linked_role.autoscaling_group.arn
  suspended_processes       = []
  target_group_arns         = []
  termination_policies      = []
  # ヘルスチェック設定
  health_check_grace_period = 300      # ヘルスチェックの猶予期間: 300秒
  health_check_type         = "EC2"    # ヘルスチェックタイプ（EC2またはELB）

  launch_template {
    id      = aws_launch_template.api_launch_template.id
    version = "$Latest"
  }
  # 置き換え動作: 終了する前に起動
  instance_refresh {
    strategy = "Rolling"
    preferences {
      min_healthy_percentage = 90
      instance_warmup        = 300
    }
  }

}

resource "aws_ecs_task_definition" "api_ecs_task_definition" {
  container_definitions = jsonencode(
    [
      {
        cpu              = 0
        environment      = []
        environmentFiles = [
          {
            type  = "s3"
            value = "${aws_s3_bucket.ecs_env_bucket.arn}/.env"
          },
        ]
        essential        = true
        image            = "${aws_ecr_repository.api_ecr_repository.repository_url}:latest"
        logConfiguration = {
          logDriver = "awslogs"
          options   = {
            awslogs-create-group  = "true"
            awslogs-group         = "/ecs/${var.env}-${var.serviceName}-service-api"
            awslogs-region        = var.region
            awslogs-stream-prefix = "ecs"
          }
          secretOptions = []
        }
        mountPoints  = []
        name         = "api"
        portMappings = [
          {
            appProtocol   = "http"
            containerPort = 80
            hostPort      = 80
            name          = "api-http"
            protocol      = "tcp"
          },
        ]
        ulimits     = []
        volumesFrom = []
      },
    ]
  )
  cpu                      = "768"
  family                   = "${var.env}-${var.serviceName}-ecs-service-task-definition-api"
  memory                   = "717"
  network_mode             = "awsvpc"
  requires_compatibilities = ["EC2",]
  task_role_arn            = aws_iam_role.ecsTaskExecutionRole.arn // TODO: task_execution_role 本当は分けたほうがいい
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn // TODO: task_execution_role

  runtime_platform {
    cpu_architecture        = "X86_64"
    operating_system_family = "LINUX"
  }
}

resource "aws_ecs_task_definition" "migration_ecs_task_definition" {
  container_definitions = jsonencode(
    [
      {
        cpu              = 0
        environment      = []
        environmentFiles = [
          {
            type  = "s3"
            value = "${aws_s3_bucket.ecs_env_bucket.arn}/.env"
          },
        ]
        essential        = true
        image            = "${aws_ecr_repository.migration_ecr_repository.repository_url}:latest"
        logConfiguration = {
          logDriver = "awslogs"
          options   = {
            awslogs-create-group  = "true"
            awslogs-group         = "/ecs/${var.env}-${var.serviceName}-service-migration"
            awslogs-region        = var.region
            awslogs-stream-prefix = "ecs"
          }
          secretOptions = []
        }
        mountPoints  = []
        name         = "migration"
        portMappings = []
        ulimits      = []
        volumesFrom  = []
      },
    ]
  )
  cpu                      = "768"
  family                   = "${var.env}-${var.serviceName}-ecs-service-task-definition-migration"
  memory                   = "717"
  network_mode             = "awsvpc"
  requires_compatibilities = ["EC2",]
  task_role_arn            = aws_iam_role.ecsTaskExecutionRole.arn // TODO: task_execution_role 本当は分けたほうがいい
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn // TODO: task_execution_role

  runtime_platform {
    cpu_architecture        = "X86_64"
    operating_system_family = "LINUX"
  }
}

// NOTE:importしていない TODO: これいらないかも
resource "aws_cloudwatch_log_group" "api_capacity_provider_gateway_log_group" {
  name = "/aws/capacity-provider/${var.env}-${var.serviceName}-api"
}

# TODO: CapacityProviderは手動でECSクラスターに設定する。
resource "aws_ecs_capacity_provider" "api_capacity_provider_strategy" {
  name = "Infra-ECS-Cluster-${aws_ecs_cluster.api_ecs_cluster.name}-EC2CapacityProvider"

  auto_scaling_group_provider {
    auto_scaling_group_arn         = aws_autoscaling_group.api_autoscaling_group.arn
    managed_termination_protection = "DISABLED"

    managed_scaling {
      instance_warmup_period    = 300
      maximum_scaling_step_size = 10000
      minimum_scaling_step_size = 1
      status                    = "ENABLED"
      target_capacity           = 100
    }
  }
}
resource "aws_ecs_cluster" "api_ecs_cluster" {
  name = "${var.env}-${var.serviceName}-ecs-cluster-api"
  configuration {
    execute_command_configuration {
      logging = "DEFAULT"
    }
  }
  # Service Connect名前空間を有効化
  service_connect_defaults {
    namespace = aws_service_discovery_http_namespace.service_connect_namespace.arn
  }

  setting {
    name  = "containerInsights"
    value = "disabled"
  }
}

resource "aws_ecs_service" "api_ecs_service" {
  cluster                            = aws_ecs_cluster.api_ecs_cluster.arn
  deployment_maximum_percent         = 100
  deployment_minimum_healthy_percent = 0
  desired_count                      = 1
  enable_ecs_managed_tags            = true
  enable_execute_command             = true
  health_check_grace_period_seconds  = 0
  #  iam_role                           = "/aws-service-role/ecs.amazonaws.com/AWSServiceRoleForECS" # たぶんAWS default なんかコメントアウトするといい感じになるらしい
  name                               = "${var.env}-${var.serviceName}-ecs-service-api"
  propagate_tags                     = "NONE"
  scheduling_strategy                = "REPLICA"
  task_definition                    = aws_ecs_task_definition.api_ecs_task_definition.arn
  triggers                           = {}

  capacity_provider_strategy {
    base              = 0
    capacity_provider = aws_ecs_capacity_provider.api_capacity_provider_strategy.name // TODO: 次回確認 なんかキャパシティプロバイダの名前とかID関係の指定にもともと失敗していた気がする。
    weight            = 1
  }

  // TODO: 次回デプロイ時に確認。ServiceConnectが設定されていればOK なぜかうまくいかない。。。一旦手順に加える

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  deployment_controller {
    type = "ECS"
  }

  network_configuration {
    assign_public_ip = false
    security_groups  = var.securityGroupIds
    subnets          = var.privateSubnetIds
  }


  ordered_placement_strategy {
    field = "instanceId"
    type  = "spread"
  }
  placement_constraints {
    type = "distinctInstance"
  }
  service_connect_configuration {
    enabled   = true
    namespace = aws_service_discovery_http_namespace.service_connect_namespace.arn

    service {
      port_name      = "api-http"  # タスク定義のポートマッピング名と一致させる
      discovery_name = "api"
      client_alias {
        port     = 80
        dns_name = "api"
      }
    }

    log_configuration {
      log_driver = "awslogs"
      options = {
        "awslogs-group"         = "/ecs/${var.env}-${var.serviceName}-service-connect",
        "awslogs-region"        = var.region,
        "awslogs-stream-prefix" = "service-connect"
      }
    }
  }

}
# キャパシティプロバイダーの定義
resource "aws_ecs_capacity_provider" "api_capacity_provider" {
  name = "${var.env}-${var.serviceName}-capacity-provider"

  auto_scaling_group_provider {
    auto_scaling_group_arn = aws_autoscaling_group.api_autoscaling_group.arn

    managed_scaling {
      maximum_scaling_step_size = 1000
      minimum_scaling_step_size = 1
      status                    = "ENABLED"
      target_capacity           = 100
    }
    managed_termination_protection = "ENABLED"
  }
}
# キャパシティプロバイダーとクラスターの関連付け
resource "aws_ecs_cluster_capacity_providers" "api_cluster_capacity_providers" {
  cluster_name = aws_ecs_cluster.api_ecs_cluster.name

  capacity_providers = [aws_ecs_capacity_provider.api_capacity_provider_strategy.name]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = aws_ecs_capacity_provider.api_capacity_provider_strategy.name
  }
}

# タスク以降は手動（マニュアル操作推奨）
resource "aws_s3_bucket" "ecs_env_bucket" {
  bucket = "${var.env}-${var.serviceName}-service-env"
}
####################################################
#                  S3 Bucket
####################################################
resource "aws_s3_bucket" "frontend_bucket" {
  bucket = "${var.env}-${var.serviceName}-frontend"
}

resource "aws_s3_bucket_ownership_controls" "frontend_bucket_ownership" {
  bucket = aws_s3_bucket.frontend_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "frontend_bucket_access" {
  bucket                  = aws_s3_bucket.frontend_bucket.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "frontend_bucket_versioning" {
  bucket = aws_s3_bucket.frontend_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# CloudFrontからのアクセスのみを許可するバケットポリシー
resource "aws_s3_bucket_policy" "frontend_bucket_policy" {
  bucket = aws_s3_bucket.frontend_bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowCloudFrontServicePrincipal"
        Effect    = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.frontend_bucket.arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.frontend_distribution.arn
          }
        }
      }
    ]
  })
}

####################################################
#            CloudFront Origin Access Control
####################################################
resource "aws_cloudfront_origin_access_control" "frontend_oac" {
  name                              = "${var.env}-${var.serviceName}-oac"
  description                       = "OAC for ${var.env}-${var.serviceName} frontend"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

####################################################
#                  CloudFront
####################################################
resource "aws_cloudfront_distribution" "frontend_distribution" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "${var.env}-${var.serviceName}-frontend"
  default_root_object = "index.html"
  price_class         = "PriceClass_200" # 北米、欧州、アジア

  # S3オリジン
  origin {
    domain_name              = aws_s3_bucket.frontend_bucket.bucket_regional_domain_name
    origin_id                = "S3-${aws_s3_bucket.frontend_bucket.bucket}"
    origin_access_control_id = aws_cloudfront_origin_access_control.frontend_oac.id
  }

  # APIGatewayオリジン
  origin {
    domain_name = "${aws_apigatewayv2_api.api_gateway.id}.execute-api.${var.region}.amazonaws.com"
    origin_id   = "ApiGateway-${aws_apigatewayv2_api.api_gateway.name}"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  # S3のデフォルトキャッシュ動作
  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "S3-${aws_s3_bucket.frontend_bucket.bucket}"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    compress               = true
  }

  # APIGatewayへのパスパターン（第一優先）
  ordered_cache_behavior {
    path_pattern     = "/api/*"
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "ApiGateway-${aws_apigatewayv2_api.api_gateway.name}"

    forwarded_values {
      query_string = true
      headers      = ["Authorization", "Origin", "Referer"]
      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 0
    max_ttl                = 0
  }

  # SPAルーティング対応（404エラーをindex.htmlにリダイレクト）
  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/index.html"
    error_caching_min_ttl = 10
  }

  custom_error_response {
    error_code         = 404
    response_code      = 200
    response_page_path = "/index.html"
    error_caching_min_ttl = 10
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
    # カスタムドメインを使用する場合は以下を有効化
    # acm_certificate_arn = aws_acm_certificate.cert.arn
    # ssl_support_method  = "sni-only"
    # minimum_protocol_version = "TLSv1.2_2021"
  }

  tags = {
    Name = "${var.env}-${var.serviceName}-cloudfront"
  }
}

####################################################
#                  CodePipeline
####################################################
# フロントエンド用のS3アーティファクトバケット
# フロントエンド用のS3アーティファクトバケット
resource "aws_s3_bucket" "codepipeline_bucket" {
  bucket = "${var.env}-${var.serviceName}-codepipeline-artifacts"
}

resource "aws_s3_bucket_ownership_controls" "codepipeline_bucket_ownership" {
  bucket = aws_s3_bucket.codepipeline_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

# アーティファクトのライフサイクル設定を追加
resource "aws_s3_bucket_lifecycle_configuration" "codepipeline_bucket_lifecycle" {
  bucket = aws_s3_bucket.codepipeline_bucket.id

  rule {
    id     = "delete-old-artifacts"
    status = "Enabled"

    # 古いバージョンを7日後に削除
    expiration {
      days = 1
    }

    # 不完全なマルチパートアップロードを1日後に削除
    abort_incomplete_multipart_upload {
      days_after_initiation = 1
    }
  }
}

resource "aws_s3_bucket_versioning" "codepipeline_bucket_versioning" {
  bucket = aws_s3_bucket.codepipeline_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# CodeBuild用ロール
resource "aws_iam_role" "codebuild_role" {
  name = "${var.env}-${var.serviceName}-codebuild-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "codebuild.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "codebuild_policy" {
  role = aws_iam_role.codebuild_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:GetObjectVersion",
          "s3:GetBucketAcl",
          "s3:GetBucketLocation"
        ]
        Resource = [
          aws_s3_bucket.codepipeline_bucket.arn,
          "${aws_s3_bucket.codepipeline_bucket.arn}/*",
          aws_s3_bucket.frontend_bucket.arn,
          "${aws_s3_bucket.frontend_bucket.arn}/*"
        ]
      }
    ]
  })
}

# CodeBuildプロジェクト
resource "aws_codebuild_project" "frontend_build" {
  name          = "${var.env}-${var.serviceName}-ui"
  description   = "Vue.jsフロントエンドのビルドプロジェクト"
  build_timeout = "15"
  service_role  = aws_iam_role.codebuild_role.arn

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    type                        = "LINUX_CONTAINER"
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-x86_64-standard:4.0"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "NODE_ENV"
      value = "development"
    }

    environment_variable {
      name  = "VUE_APP_API_BASE"
      value = "https://${aws_cloudfront_distribution.frontend_distribution.domain_name}"
    }
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = <<-EOF
      version: 0.2
      phases:
        install:
          commands:
            - cd gant-proto
            - ls -la
            - npm install
        pre_build:
          commands:
            - ls -la
            - npm install
        build:
          commands:
            - ls -la
            - npm run build
      artifacts:
        files:
          - '**/*'
        base-directory: 'gant-proto/dist'
        discard-paths: no
    EOF
  }

  logs_config {
    cloudwatch_logs {
      group_name  = "/aws/codebuild/${var.env}-${var.serviceName}-frontend-build"
      stream_name = "%CODEBUILD_ID%"
    }
  }
}

# CodePipeline用ロール
resource "aws_iam_role" "codepipeline_role" {
  name = "${var.env}-${var.serviceName}-codepipeline-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "codepipeline.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "codepipeline_policy" {
  name = "${var.env}-${var.serviceName}-codepipeline-policy"
  role = aws_iam_role.codepipeline_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:GetObjectVersion",
          "s3:GetBucketVersioning",
          "s3:PutObject"
        ]
        Resource = [
          aws_s3_bucket.codepipeline_bucket.arn,
          "${aws_s3_bucket.codepipeline_bucket.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "codebuild:BatchGetBuilds",
          "codebuild:StartBuild"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "codestar-connections:UseConnection",
          "codestar-connections:GetConnection",
          "codestar-connections:ListConnections"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:*"
        ]
        Resource = [
          aws_s3_bucket.frontend_bucket.arn,
          "${aws_s3_bucket.frontend_bucket.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "cloudfront:CreateInvalidation"
        ]
        Resource = aws_cloudfront_distribution.frontend_distribution.arn
      }
    ]
  })
}

# GitHub接続
resource "aws_codestarconnections_connection" "github_connection" {
  name          = "github-connection"
  provider_type = "GitHub"
}

# Lambda関数のコードをインラインで定義
resource "local_file" "lambda_function" {
  content  = <<-EOF
const AWS = require('aws-sdk');

exports.handler = async (event, context) => {
  console.log('Received event:', JSON.stringify(event, null, 2));

  const cloudfront = new AWS.CloudFront();
  console.log(event);
  const params = event;

  const invalidationParams = {
    DistributionId: params.distribution_id,
    InvalidationBatch: {
      CallerReference: new Date().getTime().toString(),
      Paths: {
        Quantity: params.paths.length,
        Items: params.paths
      }
    }
  };

  try {
    const response = await cloudfront.createInvalidation(invalidationParams).promise();
    console.log('Invalidation created successfully:', JSON.stringify(response, null, 2));
    return {
      statusCode: 200,
      body: JSON.stringify('Invalidation created successfully!')
    };
  } catch (error) {
    console.error('Error creating invalidation:', error);
    throw error;
  }
};
  EOF
  filename = "${path.module}/invalidate_cloudfront_cache.js"
}

# Lambda関数用のZIPファイルを動的に作成
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = local_file.lambda_function.filename
  output_path = "${path.module}/invalidate_cloudfront_cache.zip"
  depends_on  = [local_file.lambda_function]
}

# CloudFront Cache Invalidation用のLambda関数
resource "aws_lambda_function" "invalidate_cloudfront_cache" {
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  function_name    = "${var.env}-${var.serviceName}-invalidate-cloudfront-cache"
  role             = aws_iam_role.lambda_execution_role.arn
  handler          = "invalidate_cloudfront_cache.handler"
  runtime          = "nodejs16.x"
  timeout          = 30
}

# Lambda用IAMロール
resource "aws_iam_role" "lambda_execution_role" {
  name = "${var.env}-${var.serviceName}-lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "${var.env}-${var.serviceName}-lambda-policy"
  role = aws_iam_role.lambda_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "cloudfront:CreateInvalidation"
        ]
        Resource = "*"
      }
    ]
  })
}

# CodePipelineがLambda関数を呼び出すための権限を追加
resource "aws_iam_role_policy" "codepipeline_lambda_policy" {
  name = "codepipeline-lambda-policy"
  role = aws_iam_role.codepipeline_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "lambda:InvokeFunction"
        ],
        Resource = aws_lambda_function.invalidate_cloudfront_cache.arn,
        Effect   = "Allow"
      }
    ]
  })
}

# 手動トリガーのみのCodePipeline
resource "aws_codepipeline" "frontend_pipeline" {
  name     = "${var.env}-${var.serviceName}-frontend-pipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = aws_s3_bucket.codepipeline_bucket.bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        ConnectionArn        = aws_codestarconnections_connection.github_connection.arn
        FullRepositoryId     = var.github_repository # 例: "your-organization/your-repo"
        BranchName           = var.github_branch     # 例: "main"
        DetectChanges        = "false"               # ブランチの変更によるトリガーを無効化
      }
    }
  }

  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.frontend_build.name
      }
    }
  }

  stage {
    name = "Manual-Approval"

    action {
      name     = "Manual-Approval"
      category = "Approval"
      owner    = "AWS"
      provider = "Manual"
      version  = "1"

      configuration = {
        CustomData = "フロントエンドをS3バケットにデプロイする承認をお願いします"
      }
    }
  }

  stage {
    name = "Deploy"

    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "S3"
      input_artifacts = ["build_output"]
      version         = "1"

      configuration = {
        BucketName = aws_s3_bucket.frontend_bucket.bucket
        Extract    = "true"
      }
    }
  }

  stage {
    name = "Invalidate"

    action {
      name      = "InvalidateCache"
      category  = "Invoke"
      owner     = "AWS"
      provider  = "Lambda"
      version   = "1"
      region    = var.region
      run_order = 1

      configuration = {
        FunctionName   = aws_lambda_function.invalidate_cloudfront_cache.function_name
        UserParameters = jsonencode({
          distribution_id = aws_cloudfront_distribution.frontend_distribution.id
          paths           = ["/*"]
        })
      }
    }
  }
}

# 変数定義
variable "github_repository" {
  description = "GitHubリポジトリ名 (例: organization/repo-name)"
  type        = string
}

variable "github_branch" {
  description = "GitHubブランチ名"
  type        = string
  default     = "main"
}

####################################################
#                  Outputs
####################################################
output "cloudfront_domain_name" {
  description = "CloudFront distribution domain name"
  value       = aws_cloudfront_distribution.frontend_distribution.domain_name
}

output "s3_bucket_name" {
  description = "Frontend S3 bucket name"
  value       = aws_s3_bucket.frontend_bucket.bucket
}

output "s3_bucket_arn" {
  description = "Frontend S3 bucket ARN"
  value       = aws_s3_bucket.frontend_bucket.arn
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID"
  value       = aws_cloudfront_distribution.frontend_distribution.id
}

