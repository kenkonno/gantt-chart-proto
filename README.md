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


## メモ程度なので漏れがありそう


