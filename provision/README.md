# オンデマンド化

## 概要

オンデマンド化の提供方法の手順を示す。

## 事前準備

1. DBの初期化が必要か確認する。必要であれば、backend/docker/postgres-init/init-data/init.sql へ設置する。

## 手順

1. build.sh を 対象の顧客がわかるように引数付きで実行する。（ex. sh build.sh mee
2. provision/XXXX ディレクトリ配下にimageが出力される
3. .env, docker-compose.yaml の設置を行う。.envのパスワードは変更すること。
4. provision/XXXX ディレクトリからサービスが起動することを確認する

## 顧客へ提供する構築手順

1. 作業ディレクトリ移動
2. docker image と .env, docker-compose.yaml を作業ディレクトリへ格納
3. docker imageの読み込み
    ```aiignore
    # PostgreSQLイメージをtarから読み込み
    docker load -i tasmap_postgres.tar
    # APIイメージをtarから読み込み
    docker load -i tasmap_api.tar
    # マイグレーションイメージをtarから読み込み
    docker load -i tasmap_migration.tar
    # セッションイメージをtarから読み込み
    docker load -i tasmap_session.tar
    # Webイメージをtarから読み込み
    docker load -i tasmap_web.tar
   
    # （オプション）データ初期化イメージをtarから読み込み
    docker load -i tasmap_init_data.tar
    ```
4. データベース用のvolume作成
    ```aiignore
    docker volume create dbdata_tasmap
    ```
5. アプリケーションの起動
    ```aiignore
    docker compose up -d
    ```
6. (オプション)データの初期化
    ```aiignore
    docker compose --profile init run --rm tasmap_db_init
    ```

以上。

## バージョンアップ方法

dockerイメージを読み込みなおしてコンテナーを再起動してください。

コンテナの読み込み直し。

```aiignore
# PostgreSQLイメージをtarから読み込み
docker load -i tasmap_postgres.tar
# APIイメージをtarから読み込み
docker load -i tasmap_api.tar
# マイグレーションイメージをtarから読み込み
docker load -i tasmap_migration.tar
# セッションイメージをtarから読み込み
docker load -i tasmap_session.tar
# Webイメージをtarから読み込み
docker load -i tasmap_web.tar

# （オプション）データ初期化イメージをtarから読み込み
docker load -i tasmap_init_data.tar
```

サービスの再起動

```aiignore
docker compose down

docker compose up -d
```

## その他

1. 各ポートはdocker-compose.yamlを参照してください。HOSTからコンテナーはdefaultではWebサーバー用の80番ポートのみ公開しています。