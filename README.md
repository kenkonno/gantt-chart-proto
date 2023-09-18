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