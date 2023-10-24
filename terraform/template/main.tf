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

resource "aws_s3_bucket" "example_bucket" {
  bucket = "${var.env}-${var.serviceName}-example-bucket"
}