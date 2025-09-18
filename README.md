# Ken's All in Project set.

- フロント・バックエンド・バッチ完備
- gormモデル定義から ソースコードの自動生成
- openapi.yamlからのソースコード自動生成

# HowToUse

1. dockerのvolumeの作成（この名前はユニークなものにすること）
    1. 例）docker volume create dbdata_gantt
2. docker-compose.yamlの編集
    1. 20行目のvolumeと63行目のvolumeを1.iで作成したものに変更する。
    2. recreate.shの修正
3. パスorパッケージの修正
    1. DockerFileの修正
    2. openapi_generator.shの修正
    3. openapi_converterの修正


## 残タスク
1. チケットの関連
   1. コメント機能
   2. 更新履歴機能
2. マイルストーン機能
3. 削除と上下移動

# やること整理
1. EC2に上げる
2. 週次ビューの追加
3. 全体ビューの追加
4. マイルストーの追加
5. チケットコメント関連の追加


## s3 ec2にしたときの話

1. twipとの共存を目指した
2. twipのapiサーバーだけは落としたまま。batchとpostgresは共存
3. gantt側のコンテナ名を変更した
4. postgresはhost側のポートを5433にした（コンテナ側は変わらず5432）


## 環境増設したときのメモ
環境増やした時のメモ
・ボリュームの作成
[ec2-user@ip-172-31-16-212 gantt-chart-proto]$ docker volume create dbdata_gantt_2
dbdata_gantt_2

・docker-compose.yamlの変更（portsの変更。volumesの変更）
ports:
- "5434:5432"
volumes:
- dbdata_gantt_2:/var/lib/postgresql/data

    # golang アプリケーション
gantt_api:
container_name: gantt_api
ports:
- "8082:8081"

volumes:
dbdata_gantt_2:
external: true

・.envの変更
ENVIROMENT=DEVELPMENT

## Postgresql
POSTGRES_PORT=5434

## Api
API_PORT=8081


色々やったけど、結局IAM周りとかでだるくなったのでteraform化したいですね。

## 社内・デモ環境

docker-compose -p env1 down
docker-compose -p env1 up -d gantt_postgres
docker-compose -p env1 up -d gantt_session
docker-compose -p env1 up -d gantt_api

docker-compose -p env1 logs gantt_api

docker-compose -p env1 up gantt_migration

docker-compose -p env2 down
docker-compose -p env2 up -d gantt_postgres_2
docker-compose -p env2 up -d gantt_session_2
docker-compose -p env2 up -d gantt_api_2

docker-compose -p env2 logs gantt_api_2

docker-compose -p env2 up gantt_migration_2

docker logs $(docker ps -ql) -f

# AWS環境へbuild
$ aws ecr get-login-password --region ap-northeast-1 --profile=dev-laurensia | docker login --username AWS --password-stdin 866026585491.dkr.ecr.ap-northeast-1.amazonaws.com
$ aws ecr get-login-password --region ap-northeast-1 --profile=epson-prod | docker login --username AWS --password-stdin 339712996936.dkr.ecr.ap-northeast-1.amazonaws.com
$ aws ecr get-login-password --region ap-northeast-1 --profile=ftech-prod | docker login --username AWS --password-stdin
724772070484.dkr.ecr.ap-northeast-1.amazonaws.com
Login Succeeded
docker-compose build
docker-compose push

## トライアル環境データ削除作業手順
1. 環境を01～06のどれかを確認する
2. ssh_tunnel.shで対象の環境にアクセス
3. datagripからデータを確認する
4. pg_dumpでデータを出力、ファイル名の先頭にdemo0Xをつける（pg_dumpのバージョンは17を指定すること）
5. s3にアップロード(https://ap-northeast-1.console.aws.amazon.com/s3/buckets/tasmap-trial-db-backup?region=ap-northeast-1&bucketType=general&tab=objects)
6. adhookからデータ削除SQLを実行
7. データが削除されたことを確認
8. Web画面からadmin/itumonoでログインできることを確認。


1. https://d2l5ymvdgzq1kr.cloudfront.net
2. https://d1jcfwj1966idb.cloudfront.net
3. https://d1y9r0167uolwh.cloudfront.net
4. https://dggi1qvva9qqy.cloudfront.net
5. https://dyjhvazbfc80z.cloudfront.net
6. https://d1lqvglun25qns.cloudfront.net

## トライアル環境
docker-compose -p env01 -f docker-compose-demo.yml down
docker-compose -p env01 -f docker-compose-demo.yml up -d gantt_postgres
docker-compose -p env01 -f docker-compose-demo.yml up -d gantt_session
docker-compose -p env01 -f docker-compose-demo.yml up -d gantt_api
docker-compose -p env01 -f docker-compose-demo.yml logs gantt_api
docker-compose -p env01 -f docker-compose-demo.yml run gantt_migration

docker-compose -p env02 -f docker-compose-demo.yml down
docker-compose -p env02 -f docker-compose-demo.yml up -d gantt_postgres
docker-compose -p env02 -f docker-compose-demo.yml up -d gantt_session
docker-compose -p env02 -f docker-compose-demo.yml up -d gantt_api
docker-compose -p env02 -f docker-compose-demo.yml logs gantt_api
docker-compose -p env02 -f docker-compose-demo.yml run gantt_migration

docker-compose -p env03 -f docker-compose-demo.yml down
docker-compose -p env03 -f docker-compose-demo.yml up -d gantt_postgres
docker-compose -p env03 -f docker-compose-demo.yml up -d gantt_session
docker-compose -p env03 -f docker-compose-demo.yml up -d gantt_api
docker-compose -p env03 -f docker-compose-demo.yml logs gantt_api
docker-compose -p env03 -f docker-compose-demo.yml run gantt_migration

docker-compose -p env04 -f docker-compose-demo.yml down
docker-compose -p env04 -f docker-compose-demo.yml up -d gantt_postgres
docker-compose -p env04 -f docker-compose-demo.yml up -d gantt_session
docker-compose -p env04 -f docker-compose-demo.yml up -d gantt_api
docker-compose -p env04 -f docker-compose-demo.yml logs gantt_api
docker-compose -p env04 -f docker-compose-demo.yml run gantt_migration

docker-compose -p env05 -f docker-compose-demo.yml down
docker-compose -p env05 -f docker-compose-demo.yml up -d gantt_postgres
docker-compose -p env05 -f docker-compose-demo.yml up -d gantt_session
docker-compose -p env05 -f docker-compose-demo.yml up -d gantt_api
docker-compose -p env05 -f docker-compose-demo.yml logs gantt_api
docker-compose -p env05 -f docker-compose-demo.yml run gantt_migration

docker-compose -p env06 -f docker-compose-demo.yml down
docker-compose -p env06 -f docker-compose-demo.yml up -d gantt_postgres
docker-compose -p env06 -f docker-compose-demo.yml up -d gantt_session
docker-compose -p env06 -f docker-compose-demo.yml up -d gantt_api
docker-compose -p env06 -f docker-compose-demo.yml logs gantt_api
docker-compose -p env06 -f docker-compose-demo.yml run gantt_migration

