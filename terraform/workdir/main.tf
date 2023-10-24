terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
  required_version = ">= 1.1.8"
}

provider "aws" {
  access_key = var.accessKey
  secret_key = var.secretKey
  region     = "ap-northeast-1"
}
