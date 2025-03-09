# Simple-Ecs-on-Ec2

## 説明

ECS on EC2 を構築します。
RDSは対象外です。

## リソース

- Api Gateway
- ECR
- EC2(EC2インスタンス自体は作成しない)
  - Auto Scaling グループ
  - 定義テンプレート
- S3（ECR用の環境変数）
- ECS
  - ネットワーク設定：awcvpc
  - ServiceConnect(CloudMap)

## required

- SSH キーペア
- VPC（注意！特に前回パブリックで構築してしまったので、次回は必ずプライベートVPCで構築すること）
- AWSの自動生成リソースがあるから何回かapplyする必要あり。

## apply後
 - DBサーバーの構築
   - 追加設定でDB名を必ず指定すること。
   - パラメータグループで force_ssl を 0 にする
 - Redisサーバーの構築
   - 何かわからんけどredis oss とかになってつながらなくなった
   - 独自クラスター、で全部シンプルにして暗号化を全部いいえにすればつながった（たぶんデフォルトで暗号化ＯＮだとクライアント側で対応が必要そう。てかこの辺ＤＢもあるからデフォに改善すべきだな）
   - 独自のキャッシュを設計 - クラスターキャッシュ - 無効 の組み合わせでトライ cache.t3.micro
   - セキュリティ設定で全部無効化
   - サービスは作り直さないとなぜかキャパシティプロバイダから発見されなかったので作り直す
 - ECRにコンテナーをpushする
   - ローカルの aws profileに追加する
   - docker-compose.yamlに定義を追加する（不要なものをコメントアウト）
   - ECRメモの項を参照しECRへpushする
   - ☆ECSに紐づいている名前空間は使えないっぽい。たぶん一個したのサービスコネクトと同じ話。CloudMapのサービス名がぶつかっている模様。
   - サービスコネクトは作り直さないといけなかった。作り直したらAPIGateway側の統合設定も直す。
   - キャパシティプロバイダが反応しないので手動で作成しなおす。
     - autoscaling group から作成しなおし。ssh-key pair がなぜかmdsになっていた。
 - CodePipeLineを構築する
   - S3にフロントのファイルをdeployする
 - s3にenvファイルを設置する
 - Apigateway で api ステージにデプロイする（スロットの設定に注意）
 - Cloudfront
    - APIGatewayとの連携
    - /api/*のビヘイビアを指定する
- ECS

# ECR メモ
ログインする
docker-compose.yamlにECRの定義を追加すれば ログイン・build・pushでいけるみたい

$ aws ecr get-login-password --region ap-northeast-1 --profile=dev-laurensia | docker login --username AWS --password-stdin 866026585491.dkr.ecr.ap-northeast-1.amazonaws.com
$ aws ecr get-login-password --region ap-northeast-1 --profile=epson-prod | docker login --username AWS --password-stdin 339712996936.dkr.ecr.ap-northeast-1.amazonaws.com
$ aws ecr get-login-password --region ap-northeast-1 --profile=mds-prod | docker login --username AWS --password-stdin 084828592402.dkr.ecr.ap-northeast-1.amazonaws.com
Login Succeeded
docker-compose build
docker-compose push

# MDS構築時のメモ
なんか色々上手くいかなくなっている。
結局手動でやってりうので上手くいかなかったところはterraformから一旦外した方が便利になりそう。
- オートスケーリング
- サービスコネクト
- redisのSSL
- DBのSSL

# 運用について
TOBE: まとめていく。現状だとフロントはmasterへpush。バックエンドはbuild後各ECRへpushしてサービスから今日再デプロイをする必要がある。

一旦は、各お客様の環境は 手動のcodepipelineにして手動デプロイかな。

一旦はローカルでスクリプト組めば回るけど将来的にはもうちょい考えないといけないね。
というかスキーマ分けしたサービスを構築できればアプリケーションは１個で済むはずだ。


## 手動構築した時の注意事項

- ECSのサービスは作成時のみネットワークをオプションで設定できるのでVPCを間違わないように注意（編集だと出てこない）
- ApiGateway -> Cloud Map -> ECS は ECS側がService Connect じゃないとNG（ポート指定付きじゃないといけないみたい） ECS側はポートマッピングされるのでなんでもOK
- DBは色々悩んで無料枠使ったほうが結局安い
- AutoscalingはMAX 2, 必要 1 としてmigrationタスクを立ち上げれるようにする
- ECSはメモリをインスタンスから-250Mクライを最大予約にしないと立ち上がらない
- ローリングアップデート（ダウンタイム無し）設定は 0%, 100% とする（本番設定はしなくてもいいかも）
- CloudFront to Apigateway は CloudaFront側で /api* としたとき ApiGatewayのステージも api とする
- RDSは最初にDB作成しないとなぜかパブリック接続できなかった（DNSの伝播問題なだけかも 
- Apiagtewayのログは左メニューの設定を事前にする、しかもなぜか必要Roleを設定したのに上手くいかんかった、時間が解決した
- VPCは適当にデフォで作ったやつで全然よさそう

## メモ

- クラスターのコンテナインスタンスはAuto Scalingグループが表示されているだけ
- AutoScalingグループを作成するとec2インスタンスが起動される
    - 肝心のEBS設定は起動テンプレートのみで指定可能っぽい
    - VPNは作ってあげたほうがいいのかも？

まずはコンソールポチポチで実際に動作するところまで頑張る。
※cloudfrontからECSは指定できないのでELBをかませる必要がありそう。

手でぽちぽちするならECRが先かも。ともかく今日やり切るぞ。

リポジトリはdocker-composeのサービスごとに作成する

よくわからないけどタスク定義を追加

- Postgresだけの定義にしたら動作した
- Apiを上げようとするとたぶんタスクのリソース上限に引っかかって終わってる？
- 試しにタスク別にするか
    - まあスケールとか考えたら別々にするのがいいのか

BuildはgitからpullしてbuildしてECRに上げるようにする。
環境変数は困る。なんかs3に上げればOKな気配はある

名前空間でHOSTを解決させるみたい

---

# 強化合宿二日目

結論

1. ECS関連は慣れたから多分大丈夫
    1. 起動テンプレート
    2. AutoScaling
    3. クラスター
    4. サービス
    5. タスク
    6. 特にハマりはリソース不足だった
2. 名前解決とかをするためにVPCネットワークにECSを配置する
    1. 起動テンプレートで定義する
    2. ホストゾーンは結局使うので利用する

## とにかく何でもメモ

1. VPC(
   env-manual-ecs-vpc) https://ap-northeast-1.console.aws.amazon.com/vpcconsole/home?region=ap-northeast-1#VpcDetails:VpcId=vpc-01aae518b895c636a
    1. publicを使ってとりあえず作る。
    2. SubネットにIPV4自動割り当て設定をする
2. EC2起動テンプレート
    1. ECSから自動生成されたものを使う必要がある。
    2. 気を付ける必要はないがネットワークは設定しない 
    2. https://ap-northeast-1.console.aws.amazon.com/ec2/home?region=ap-northeast-1#LaunchTemplates:
3. AutoScalingGroup
    1. ECSから自動生成されたものを使う必要がある。と思ったけど作成時にVPC選べるからそうじゃないかも
    2. https://ap-northeast-1.console.aws.amazon.com/ec2/home?region=ap-northeast-1#AutoScalingGroupDetails:id=Infra-ECS-Cluster-env-manual-ecs-cluster-2-e641b6e0-ECSAutoScalingGroup-MryMg96oA9LO;view=details
4. ECS(env-manual-ecs-cluster-2) https://ap-northeast-1.console.aws.amazon.com/ecs/v2/clusters/env-manual-ecs-cluster-2/services?region=ap-northeast-1
    1. セキュリティグループデフォルトにしてしまったのは失敗
    1. タスク定義(AwsVpcを使う)
        1. dev-manual-ecs-service-task-definition-api https://ap-northeast-1.console.aws.amazon.com/ecs/v2/task-definitions/dev-manual-ecs-service-task-definition-api?status=ACTIVE&region=ap-northeast-1
        2. dev-manual-ecs-service-task-definition-postgres https://ap-northeast-1.console.aws.amazon.com/ecs/v2/task-definitions/dev-manual-ecs-service-task-definition-postgres?status=ACTIVE&region=ap-northeast-1
        3. プロビジョニング戦略はEC2の方を選ばないと１インスタンスに複数入れられない
       4. メモリーはちょっと小さめにするのがコツっぽい
    2. サービス定義
        1. 名前空間
            1. 新しく作る必要があるのはたぶんServiceConnect専用だから（デフォルトのやつ）
            2. 作るときは真ん中のやつを選んでDNSでも解決できるようにする
            Route53のホストゾーンが作られる https://us-east-1.console.aws.amazon.com/route53/v2/hostedzones?region=us-east-1#ListRecordSets/Z051555121PM6LOF4CWKN
            3. ここに勝手にレコードが更新されるっぽい
        2. env-manual-ecs-service-postgres
        3. env-manual-ecs-service-api-1
4. ロール
   5. ecsTaskExecutionRole
   6. https://us-east-1.console.aws.amazon.com/iamv2/home?region=us-east-1#/roles/details/ecsTaskExecutionRole?section=permissions
7. セキュリティグループ
   8. https://ap-northeast-1.console.aws.amazon.com/ec2/v2/home?region=ap-northeast-1#SecurityGroup:groupId=sg-0f9b35024f6ef7192
9. DB接続
   10. intellijからはまだできない
   11. 対象のEC2に接続して「 docker container exec -it b7214eb15672 bash 」で「psql -U gantt-chart-proto」でみれる
   12. と思たtけどAPIサーバーを通すときに同じようなことはやるはず。
13. 永続かチェック
    1. タスクの終了 
       15. NG どういうことだよｗ
    2. サービスの終了
    3. インスタンスの終了
    4. おそらく起動テンプレートを分けて起動テンプレート自体のDBSを永続化させる必要がありそう。まあそりゃそうか。DBとAPPが同じスケール設定て意味わからんモナ。

気になること

1. 起動テンプレートのサブネットが１つしか選べない、２個選びたい。一旦そのまま進める。
2. autoscalingがなぜか容量の更新中になる
3. キャパシティープロバイダーってやつになった
    1. たぶんAutoScalingGroupの上位に属する者っぽい
4. サービスディスカバーでヘルスチェックが必要っぽい
5. ネットワークがらみ？
    6. どうやらサービスのネットワーキングが間違っていて、新規作成の時しか表示と変更ができないよう。💩
    7. ECS on EC2 の場合は ECSからAutoSacleGroupingを作成しないとクラスターのVPCが新しく設定したものに適応されない模様
8. service-discoverをpostgresにしちゃった → いいのかも
9. VPCバンだと ENI というネットワークインターフェース制限があるから1EC2インスタンス = 1タスクになるっぽい（でかいのだと増やせるみたいだけど）
   10. ということでmicroインスタンス２つでやる。
11. migrationの実行はコンソールから都度になるとかなりきつい。
    12. cliからタスク追加する要にすれば楽にはなる。リリース系はいったん手動で考えてるから仕方ないかな。
    13. migration事態はとりあえずできた
14. 外からつながらなくなった。
    1. つながらない分にはまあいのだが。
    2. 見てるセキュリティグループがちがったっぽい


## 二日目まとめ

結論：無料枠使い倒した方がいい

ただインフラ回りやECSについてはかなり使い倒せたし、よかった。

ECSの目標としてはAPI通信できるところまでかな。

DBのパブリックアクセスはなんか伝播の問題だったっぽい。

寝る前にまとめ。
利用リソース群

 - Cloudfrontを全ての受けにする
 - FrontはS3に設置
 - DBは素直にRDS
 - ApiはECS
 - MigrationもECS（都度タスク実行すればＯＫ、スケーリングで常時1,MAX2にすればよい）
 - Cloudfront と ECS は HTTP統合の cloud map（Service discover） を使えば上手くいきそう 
 - VPCはよくわからんけどデフォのプリセットでなんとなく上手くいく

たぶんこれで月額１万行かないくらいじゃないかなあ

 - CI/CD


## 三日目

今日はApiGateway to ECS まで頑張ろう

1. docker -> localhost:80 (成功)
```agsl
root@ip-10-0-130-227:/go/src/github.com/kenkonno/gantt-chart-proto/backend/api# curl localhost
Hello World!
```

2. EC2 ->  Service discover（成功 ECSのログにも出ていた）
```agsl
curl api.dev-manual-ecs-namespace
   Hello World!
```

3. 色々あって成功した。 

メモ：Cloudfront to apigateway は HTTPSオンリー

ApiGatewayからはServiceConnectを使うべし
https://docs.aws.amazon.com/ja_jp/apigateway/latest/developerguide/http-api-develop-integrations-private.html#http-api-develop-integrations-private-Cloud-Map

a52ed93802be


## 四日目

手動でサービスが動くところまで作って terraform import でテンプレート化まで頑張る

simple cloudfrontはOK
simple ECS を作る
Attach API を作る（simple cloudfront と simple ECS を結合させる）
simple RDS を作る

一旦 migrate タスクを実行して cloudfrontにビルド済み資産を上げる

- メモ
  - ECS上でgoのバイナリをbuildしてるけどNG。build済みimageにする必要がある（public サブネットじゃなくすため）
  - と思ったけどdownloadはdockerfileでやってるから問題ないのかも。
  - ORIGIN環境変数を追加した
    - 既存のEC2インスタンス上に追加済み

とりあえず動いたー
これで手動構築は完了といえるだろう

terraform化を目指す
https://zoo200.net/terraform-import/

- terraform取り込みコマンド週
  - import
    - apigatewayv2
      - terraform import -var-file=dev.terraform.tfvars aws_apigatewayv2_api.api_gateway u90s3jbw2h
      - terraform import -var-file=dev.terraform.tfvars aws_apigatewayv2_route.route u90s3jbw2h/g0az5th
      - terraform import -var-file=dev.terraform.tfvars aws_apigatewayv2_integration.api_integration u90s3jbw2h/xz2ib34
      - terraform import -var-file=dev.terraform.tfvars aws_apigatewayv2_vpc_link.api_vpc_link rhdah7
      - terraform import -var-file=dev.terraform.tfvars aws_apigatewayv2_stage.api_stage u90s3jbw2h
      - terraform import -var-file=dev.terraform.tfvars aws_cloudwatch_log_group.api_gateway_log_group /aws/apigateway/welcome2
    - servicediscovery
      - terraform import -var-file=dev.terraform.tfvars aws_service_discovery_http_namespace.service_discovery_http_namespace ns-jvcxb2ely7ewkskj
      - terraform import -var-file=dev.terraform.tfvars aws_service_discovery_service.service_discovery_service srv-b2qqtese6yf7ulnt
      - terraform import -var-file=dev.terraform.tfvars aws_service_discovery_private_dns_namespace.service_discovery_private_dns_namespace ns-jvcxb2ely7ewkskj:vpc-01aae518b895c636a
      - 
      - terraform import -var-file=dev.terraform.tfvars aws_service_discovery_instance.service_discovery_instance 8bd80673b23d4e25b490566049dc9d33
    - ECR
      - terraform import -var-file=dev.terraform.tfvars aws_ecr_repository.api_ecr_repository dev-manual-test/gantt_api
    - ECS
      - terraform import -var-file=dev.terraform.tfvars aws_iam_role.ecsInstanceRole ecsInstanceRole
      - terraform import -var-file=dev.terraform.tfvars aws_launch_template.api_launch_template lt-04a2c640944efe5e0
      - terraform import -var-file=dev.terraform.tfvars aws_autoscaling_group.api_autoscaling_group Infra-ECS-Cluster-env-manual-ecs-cluster-2-e641b6e0-ECSAutoScalingGroup-MryMg96oA9LO
      - terraform import -var-file=dev.terraform.tfvars aws_ecs_cluster.api_ecs_cluster env-manual-ecs-cluster-2
      - terraform import -var-file=dev.terraform.tfvars aws_ecs_service.api_ecs_service env-manual-ecs-cluster-2/env-manual-ecs-service-api
      - terraform import -var-file=dev.terraform.tfvars aws_ecs_task_definition.api_ecs_task_definition arn:aws:ecs:ap-northeast-1:866026585491:task-definition/dev-manual-ecs-service-task-definition-api:12
      - terraform import -var-file=dev.terraform.tfvars aws_ecs_task_definition.migration_ecs_task_definition arn:aws:ecs:ap-northeast-1:866026585491:task-definition/dev-manual-ecs-service-task-definition-migration:2
      - terraform import -var-file=dev.terraform.tfvars aws_ecs_capacity_provider.api_capacity_provider_strategy Infra-ECS-Cluster-env-manual-ecs-cluster-2-e641b6e0-EC2CapacityProvider-lBoxExbVXogj
      - terraform import -var-file=dev.terraform.tfvars aws_iam_role.ecsTaskExecutionRole ecsTaskExecutionRole


  - show
    - terraform state show aws_iam_role.ecsTaskExecutionRole

  - plan
    - terraform plan -var-file=dev.terraform.tfvars


----

AWS EC2 Exec なるものがあるらしい

# CLIでサービスに対して有効か

aws ecs update-service \
--cluster env-manual-ecs-cluster-2 \
--task-definition dev-manual-ecs-service-task-definition-postgres \
--service env-manual-ecs-service-postgres \
--desired-count 1 \
--enable-execute-command \
--profile=dev-laurensia

aws ecs update-service \
--cluster env-manual-ecs-cluster-2 \
--task-definition dev-manual-ecs-service-task-definition-api \
--service env-manual-ecs-service-api \
--desired-count 1 \
--enable-execute-command \
--profile=dev-laurensia

# CHECK

aws ecs describe-tasks \
--cluster env-manual-ecs-cluster-2 \
--tasks 3be2c73bcec04041bf241a5f06c4ebe1  \
--profile=dev-laurensia | grep enableExecuteCommand

aws ecs describe-tasks \
--cluster env-manual-ecs-cluster-2 \
--tasks a65ce359c2644b5895f21b8274230fd8  \
--profile=dev-laurensia | grep enableExecuteCommand

# タスクロールの更新

# 接続

aws ecs execute-command  `
--cluster env-manual-ecs-cluster-2 `
--task a65ce359c2644b5895f21b8274230fd8 `
--container api `
--interactive `
--command "bash" `
--profile=dev-laurensia

aws ecs execute-command `
--cluster env-manual-ecs-cluster-2 `
--task 3be2c73bcec04041bf241a5f06c4ebe1 `
--container postgres `
--interactive `
--command "bash" `
--profile=dev-laurensia

ちなみになぜかpowerシェルじゃないと動かなった。

root@d4d8f26ffc6d:/go/src/github.com/kenkonno/gantt-chart-proto/backend/api# cat /etc/hosts
127.0.0.1 localhost
::1 localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
127.255.0.1 posgres-5432-tcp.dev-manual-ecs-namespace
2600:f0f0:0:0:0:0:0:1 posgres-5432-tcp.dev-manual-ecs-namespace
172.17.0.4 d4d8f26ffc6d
root@d4d8f26ffc6d:/go/src/github.com/kenkonno/gantt-chart-proto/backend/api#