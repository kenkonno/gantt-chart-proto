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
3. 積み上げ機能。
   1. 特にfooterメニューはデザインも意識して領域を作成する必要がある。
   2. footerメニューがONの場合は領域を別途確保して上げる必要がある
   3. gantt-chartライブラリ側の機能として実装する。
4. 削除と上下移動




いままで人数を調整する機能だったが、担当者ベースになったので難しい。


所感
残り２周と考えると結構急ピッチでデモ版まで上げる必要がある。
facilityIdの連動は必須そう。
チケットのコメントも簡単で有効

マイルストーン、積み上げは見た目だけとりあえず作ったほうが良い。