rm -rf ./gen/*

# gen
"C:\Program Files\JetBrains\IntelliJ IDEA 2022.2.1\jbr\bin\java.exe" -jar \
  "../backend/api/openapi-generator-cli-6.6.0.jar" \
  generate -g typescript-axios -i "../backend/api/GanttChartApi.yaml" -o "gen"

mv ./gen/api.ts ./src/api
mv ./gen/base.ts ./src/api
mv ./gen/configuration.ts ./src/api
mv ./gen/index.ts ./src/api
mv ./gen/common.ts ./src/api