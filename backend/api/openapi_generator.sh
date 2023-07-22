rm -rf ./gen/*
rm -rf ./openapi_models/*
rm -rf ./openapi/*


# gen
"C:\Program Files\JetBrains\IntelliJ IDEA 2022.2.1\jbr\bin\java.exe" -jar "C:\Users\gunka\AppData\Local\JetBrains\IntelliJIdea2022.2\openapi\codegen\6ef65fe5cce7fbb8bf9ac9374d0fb142\openapi-generator-cli-5.0.0.jar" generate -g go-gin-server -i "C:\Users\gunka\git\gantt-chart-proto\backend\api\GanttChartApi.yaml" -o "C:\Users\gunka\git\gantt-chart-proto\backend\api\gen"

mv ./gen/go/model* ./openapi_models
mv ./gen/go/* ./openapi/

# exec convert
go run ./openapi_converter/main.go

cd ./openapi_models/

find ./ -type f | xargs -t -I{} sed -i -e "s/openapi/openapi_models/" {}

sed ../openapi/routers.go -i -e "s/NewRouter()/NewRouter(router *gin.Engine)/"
sed ../openapi/routers.go -i -e "s/router := gin.Default()//"