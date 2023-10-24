# 完全別環境を作成してチェック

簡易チェック（できたかできてないかだけ）
- S3 OK
- cloud map OKっぽいDNSも勝手に出来上がった
- ECR OK
- ECS
  - cluster OK 
  - task definition api OK
  - task definition migration OK
  - service api NG
  - service migration NG
- apigateway 
  - VPC LINK OK
  - stage OK
  - Route OK(destroy後の２回目)
- 起動テンプレート OK(destroy後の２回目)
- AutoScalingGroup NG

なんかよくわからんけど、なんかいかapplyしてれば勝手にできるっぽいよ

## 初回実行のログ

```agsl
╷
│ Error: creating Service Discovery HTTP Namespace (other-env-development-test-ecs-namespace): DuplicateRequest: Another operation of type CreatePrivateDnsNamespace and id 67ocdeyrreszd2ephjosvwxtormzjxau-jh7tteuq is in progress
│ {
│   RespMetadata: {
│     StatusCode: 400,
│     RequestID: "48e20595-6873-42dc-9c9b-ddb81ff5157b"
│   },
│   DuplicateOperationId: "67ocdeyrreszd2ephjosvwxtormzjxau-jh7tteuq",
│   Message_: "Another operation of type CreatePrivateDnsNamespace and id 67ocdeyrreszd2ephjosvwxtormzjxau-jh7tteuq is in progress"
│ }
│
│   with aws_service_discovery_http_namespace.service_discovery_http_namespace,
│   on main.tf line 31, in resource "aws_service_discovery_http_namespace" "service_discovery_http_namespace":
│   31: resource "aws_service_discovery_http_namespace" "service_discovery_http_namespace" {
│
╵
╷
│ Error: creating EC2 Launch Template (terraform-20231024163627180000000003): InvalidIamInstanceProfileArn.Malformed: The ARN ‘arn:aws:iam::796476764001:role/ecsInstanceRole’ is not valid. The expected format is arn:aws:iam::<account-id>:instance-profile/<instance-profile-name>.
│       status code: 400, request id: 586f3935-7af8-4b42-89b7-a99f5a5255c7
│
│   with aws_launch_template.api_launch_template,
│   on main.tf line 227, in resource "aws_launch_template" "api_launch_template":
│  227: resource "aws_launch_template" "api_launch_template" {
│
╵
╷
│ Error: creating ECS Capacity Provider (Infra-ECS-Cluster-other-env-development-test-ecs-cluster-api-EC2CapacityProvider): ClientException: The ECS Service Linked Role does not exist. Please ensure that it exists and try again.
│
│   with aws_ecs_capacity_provider.api_capacity_provider_strategy,
│   on main.tf line 379, in resource "aws_ecs_capacity_provider" "api_capacity_provider_strategy":
│  379: resource "aws_ecs_capacity_provider" "api_capacity_provider_strategy" {
│
```
エラー一覧
- ServiceDiscoveryのHTTP Namespace
  - たぶん重複してたから削除
- LunchTemplate
  - │ Error: creating EC2 Launch Template (terraform-20231024163627180000000003): InvalidIamInstanceProfileArn.Malformed: The ARN ‘arn:aws:iam::796476764001:role/ecsInstanceRole’ is not valid. The expected format is arn:aws:iam::<account-id>:instance-profile/<instance-profile-name>.
  - 自動で作られるものをアタッチするべき？ :role と :instance-profile なのでなんか違うのかも
  - instance-roleなるものを作る
- ECS Capacity Provider
  - creating ECS Capacity Provider (Infra-ECS-Cluster-other-env-development-test-ecs-cluster-api-EC2CapacityProvider): ClientException: The ECS Service Linked Role does not exist. Please ensure that it exists and try again.
  - どうやらLInkedRoleがないらしい

## 動作確認してないけどいったん動いた！

次回
1. rdsの構築
2. アプリケーションの構築
   1. docker-compose.ymlのrepositoryを変えてpushする。
   2. .envファイルをS3に設置する
   3. cloudfrontは無いけどApigatewayのURLから確認できるはず。
      4. と思ったけど無理なので attach用の terraform を追加する
      5. simple-front を実施後、apigateway と cloudfront を指定して origin と ビヘイビアを作成するだけ（一旦手動でもいいかも）
   6. フロントS3にファイルをbuildしてあげる。リクエスト先はother-envに書き換えること。
