variable "accessKey" {
  type = string
}
variable "secretKey" {
  type = string
}
variable "env" {
  type = string
}
variable "serviceName" {
  type = string
}
variable "sshKeyPairName" {
  type = string # dev-test
}
variable "region" {
  type = string # ap-northeast-1
}
variable "vpc" {
  type = string
}
variable "securityGroupIds" { # 主にECS実行インスタンス関連のsecurityGroup, HTTP, HTTPS を受け入れさせる
  type = list(string)
}
variable "privateSubnetIds" {
  type = list(string)
}