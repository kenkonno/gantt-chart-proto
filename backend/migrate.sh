# migrationは間違えないように別の docker-compose-migration.yamlにしている
# -p はプロジェクト名を変更してる。デフォルトだとディレクトリが使われるようです。
docker compose -f docker-compose-migration.yml -p migration run migration