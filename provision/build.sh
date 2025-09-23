#!/bin/bash

# TasMap プロビジョン用ビルドスクリプト
# このスクリプトは以下の処理を行います：
# 1. フロントエンドのビルド
# 2. Dockerイメージのビルド  
# 3. Dockerイメージのtar保存
# 4. ボリューム作成用のコマンド表示
#
# 使用方法: ./build.sh <出力ディレクトリ名>
# 例: ./build.sh production
# 出力先: provision/<出力ディレクトリ名>/images/

set -e  # エラーが発生した場合は即座に終了

# 引数チェック
if [ $# -eq 0 ]; then
    echo "エラー: 出力ディレクトリ名を指定してください"
    echo ""
    echo "使用方法: $0 <出力ディレクトリ名>"
    echo "例: $0 production"
    echo "出力先: \$(pwd からの相対パス)/provision/<出力ディレクトリ名>/images/"
    echo ""
    exit 1
fi

OUTPUT_DIR_NAME="$1"

echo "=== TasMap プロビジョン用ビルド開始 ==="
echo "出力ディレクトリ名: $OUTPUT_DIR_NAME"

# プロジェクトルートディレクトリを取得
SCRIPT_DIR=$(cd $(dirname $0); pwd)
PROJECT_ROOT=$(dirname "$SCRIPT_DIR")
OUTPUT_DIR="$PROJECT_ROOT/provision/$OUTPUT_DIR_NAME"

# 実行場所からの相対パスを計算
EXEC_DIR=$(pwd)
# WindowsパスをLinux形式に変換してから相対パス計算
PROJECT_ROOT_LINUX=$(echo "$PROJECT_ROOT" | sed 's|\\|/|g' | sed 's|C:|/mnt/c|')
RELATIVE_PATH=$(realpath --relative-to="$EXEC_DIR" "$PROJECT_ROOT_LINUX")

echo "プロジェクトルート: $PROJECT_ROOT"
echo "出力先ディレクトリ: $OUTPUT_DIR"

# 1. フロントエンドのビルド
echo ""
echo "=== 1. フロントエンドのビルド ==="
cd "$PROJECT_ROOT/gant-proto"

# VITE_API_BASEが設定されている場合は削除
if [ -f .env ]; then
    echo "既存の.envファイルからVITE_API_BASEを削除中..."
    sed -i '/^VITE_API_BASE/d' .env
fi

echo "npm run build を実行中..."
npm run build

if [ ! -d "dist" ]; then
    echo "エラー: distディレクトリが生成されませんでした"
    exit 1
fi

# distディレクトリをweb/distにコピー
echo "distディレクトリを backend/docker/web/dist にコピー中..."
rm -rf "$PROJECT_ROOT/backend/docker/web/dist"
cp -r dist "$PROJECT_ROOT/backend/docker/web/dist"

echo "フロントエンドビルド完了"

# 2. Dockerイメージのビルド
echo ""
echo "=== 2. Dockerイメージのビルド ==="
cd "$PROJECT_ROOT/backend"

echo "docker-compose-for-provision-build.yml を使用してビルド中..."
docker compose -f docker-compose-for-provision-build.yml build

echo "Dockerイメージビルド完了"

# 3. Dockerイメージのtar保存
echo ""
echo "=== 3. Dockerイメージのtar保存 ==="
cd "$PROJECT_ROOT"

# 保存先ディレクトリを作成
mkdir -p "$OUTPUT_DIR/images"

echo "PostgreSQLイメージを保存中..."
docker save -o "$OUTPUT_DIR/images/tasmap_postgres.tar" laurensia/tasmap-postgres

echo "APIイメージを保存中..."
docker save -o "$OUTPUT_DIR/images/tasmap_api.tar" laurensia/tasmap-api

echo "マイグレーションイメージを保存中..."
docker save -o "$OUTPUT_DIR/images/tasmap_migration.tar" laurensia/tasmap-migration

echo "セッションイメージを保存中..."
docker save -o "$OUTPUT_DIR/images/tasmap_session.tar" laurensia/tasmap-session

echo "Webイメージを保存中..."
docker save -o "$OUTPUT_DIR/images/tasmap_web.tar" laurensia/tasmap-web

echo "Dockerイメージ保存完了"

# 4. 完了メッセージと次のステップの案内
echo ""
echo "=== ビルド完了 ==="
if [ "$RELATIVE_PATH" = "." ]; then
    RELATIVE_OUTPUT_PATH="provision/$OUTPUT_DIR_NAME/images"
else
    RELATIVE_OUTPUT_PATH="$RELATIVE_PATH/provision/$OUTPUT_DIR_NAME/images"
fi
echo "全てのイメージが $RELATIVE_OUTPUT_PATH ディレクトリに保存されました："
echo "  - tasmap_postgres.tar"
echo "  - tasmap_api.tar"
echo "  - tasmap_migration.tar"
echo "  - tasmap_session.tar"
echo "  - tasmap_web.tar"
echo ""
echo "=== 次のステップ（デプロイ先で実行） ==="
echo "1. イメージの読み込み："
echo "   docker load -i $RELATIVE_OUTPUT_PATH/tasmap_postgres.tar"
echo "   docker load -i $RELATIVE_OUTPUT_PATH/tasmap_api.tar"
echo "   docker load -i $RELATIVE_OUTPUT_PATH/tasmap_migration.tar"
echo "   docker load -i $RELATIVE_OUTPUT_PATH/tasmap_session.tar"
echo "   docker load -i $RELATIVE_OUTPUT_PATH/tasmap_web.tar"
echo ""
echo "2. ボリュームの作成："
echo "   docker volume create dbdata_tasmap"
echo ""
echo "3. アプリケーションの起動："
if [ "$RELATIVE_PATH" = "." ]; then
    RELATIVE_MEE_PATH="provision/mee"
else
    RELATIVE_MEE_PATH="$RELATIVE_PATH/provision/mee"
fi
echo "   cd $RELATIVE_MEE_PATH"
echo "   docker compose up -d"
echo ""
echo "=== プロビジョン用ビルド完了 ==="