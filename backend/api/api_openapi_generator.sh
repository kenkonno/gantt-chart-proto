rm -rf ./gen/*
rm -rf ./openapi_models/*
rm -rf ./openapi/*

# bundle
npx @redocly/cli bundle ./api_spec/api.yaml -o ./api_spec/bundled.yaml

# gen
java -jar \
  "openapi-generator-cli-6.6.0.jar" \
  generate -g go-gin-server -i "./api_spec/bundled.yaml" -o "gen"

mv ./gen/go/model* ./openapi_models
mv ./gen/go/* ./openapi/

# exec convert
go run ./openapi_converter/main.go

cd ./openapi_models/

find ./ -type f -name "*.go" | xargs -t -I{} sed -i -e 's/json:"\([^"]*\)"/json:"\1" form:"\1"/g' {}

find ./ -type f | xargs -t -I{} sed -i -e "s/openapi/openapi_models/" {}

sed ../openapi/routers.go -i -e "s/NewRouter()/NewRouter(router *gin.Engine)/"
sed ../openapi/routers.go -i -e "s/router := gin.Default()//"

# 既存のスクリプトの後に追加
