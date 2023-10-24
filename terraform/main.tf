terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
  required_version = ">= 1.1.8"
}

locals {
  env = "dev"
  masterAccountId = "420302062688" // TODO: アカウントが固定値 シークレット情報
}

provider "aws" {
  profile = "default"
  region  = "us-east-1"
  default_tags {
    tags = {
      env = local.env
    }
  }
}

// dev要のロール
resource "aws_iam_role" "main_role" {
  name = "${local.env}-laurensia-role"

  assume_role_policy = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        "Effect": "Allow",
        "Principal": {
          "AWS": "arn:aws:iam::${local.masterAccountId}:root"
        },
        "Action": "sts:AssumeRole",
        "Condition": {}
      },
    ]
  })
}


// お試しS3
resource "aws_s3_bucket" "front_bucket" {
  ass
  bucket = "${local.env}-laurensia-front-bucket"
}

