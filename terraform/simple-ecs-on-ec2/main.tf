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
# たぶんいらない
#resource "aws_service_discovery_http_namespace" "service_discovery_http_namespace" {
#  name = "${var.env}-${var.serviceName}-ecs-namespace"
#}
resource "aws_service_discovery_service" "service_discovery_service" {
  name         = "api-80-tcp"
  namespace_id = aws_service_discovery_private_dns_namespace.service_discovery_private_dns_namespace.id
  type         = "HTTP"
}
resource "aws_service_discovery_private_dns_namespace" "service_discovery_private_dns_namespace" {
  name = "${var.env}-${var.serviceName}-ecs-namespace"
  vpc  = var.vpc
}

####################################################
#                  API Gateway
####################################################
// TODO: Logを許可するRoleの設定。ちょっとなんか難しそうだから最悪手動でやる
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
  integration_uri        = aws_service_discovery_service.service_discovery_service.arn
  payload_format_version = "1.0"
  request_parameters     = {}
  request_templates      = {}
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
    throttling_burst_limit   = 0
    throttling_rate_limit    = 0
  }
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
  health_check_grace_period = 0
  health_check_type         = "EC2"
  load_balancers            = []
  max_instance_lifetime     = 0
  max_size                  = 2
  metrics_granularity       = "1Minute"
  min_size                  = 0
  name                      = "Infra-ECS-Cluster-${var.env}-${var.serviceName}-ecs-cluster"
  protect_from_scale_in     = false
  service_linked_role_arn   = aws_iam_service_linked_role.autoscaling_group.arn
  suspended_processes       = []
  target_group_arns         = []
  termination_policies      = []

  launch_template {
    id      = aws_launch_template.api_launch_template.id
    version = "$Latest"
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
            name          = "api-80-tcp"
            protocol      = "tcp"
          },
        ]
        ulimits     = []
        volumesFrom = []
      },
    ]
  )
  cpu                      = "1024"
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
  cpu                      = "1024"
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
    field = "attribute:ecs.availability-zone"
    type  = "spread"
  }
  ordered_placement_strategy {
    field = "instanceId"
    type  = "spread"
  }

}
# タスク以降は手動（マニュアル操作推奨）



resource "aws_s3_bucket" "ecs_env_bucket" {
  bucket = "${var.env}-${var.serviceName}-service-env"
}