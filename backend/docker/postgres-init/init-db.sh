#!/bin/bash

echo "データベース初期化を開始します..."

# PostgreSQLの接続を待機
until pg_isready -h tasmap_postgres -U ${POSTGRES_USER:-postgres}; do
  echo "PostgreSQLの準備を待機中..."
  sleep 2
done

echo "PostgreSQLに接続可能です。初期化を実行します..."

# PGPASSWORD環境変数を設定してパスワード入力を回避
export PGPASSWORD=${POSTGRES_PASSWORD}

# 初期化SQLファイルを実行
if [ -f /init-data/init.sql ]; then
  psql -h tasmap_postgres -U ${POSTGRES_USER:-postgres} -d ${POSTGRES_DB:-postgres} -f /init-data/init.sql
  echo "データベース初期化が完了しました"
else
  echo "エラー: /init-data/init.sql が見つかりません"
  exit 1
fi
