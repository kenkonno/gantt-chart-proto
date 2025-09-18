rm -rf ./gen/*

# gen
java -jar \
  "../backend/api/openapi-generator-cli-6.6.0.jar" \
  generate -g typescript-axios -i "../backend/api/api_spec/bundled.yaml" -o "gen"

mv ./gen/api.ts ./src/api
mv ./gen/base.ts ./src/api
mv ./gen/configuration.ts ./src/api
mv ./gen/index.ts ./src/api
mv ./gen/common.ts ./src/apgi