rm -rf ./gen/*
rm -rf ./openapi_models/*
rm -rf ./openapi/*


# gen
"C:\Program Files\JetBrains\IntelliJ IDEA 2022.2.1\jbr\bin\java.exe" -jar \
  "openapi-generator-cli-6.6.0.jar" \
  generate -g go-gin-server -i "./GanttChartApi.yaml" -o "gen"

mv ./gen/go/model* ./openapi_models
mv ./gen/go/* ./openapi/

# exec convert
go run ./openapi_converter/main.go

cd ./openapi_models/

find ./ -type f | xargs -t -I{} sed -i -e "s/openapi/openapi_models/" {}

sed ../openapi/routers.go -i -e "s/NewRouter()/NewRouter(router *gin.Engine)/"
sed ../openapi/routers.go -i -e "s/router := gin.Default()//"