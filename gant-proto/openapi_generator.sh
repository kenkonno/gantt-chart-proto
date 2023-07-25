rm -rf ./gen/*

"C:\Program Files\JetBrains\IntelliJ IDEA 2022.2.1\jbr\bin\java.exe" -jar "C:\Users\gunka\AppData\Local\JetBrains\IntelliJIdea2022.2\openapi\codegen\6ef65fe5cce7fbb8bf9ac9374d0fb142\openapi-generator-cli-5.0.0.jar" generate -g typescript-axios -i "C:\Users\gunka\git\twitch_clip_project\backend\api\TwitchClipApi.yaml" -o "C:\Users\gunka\git\twitch_clip_project\front\twitch_clips\gen"

mv ./gen/api.ts ./src/api
mv ./gen/base.ts ./src/api
mv ./gen/configuration.ts ./src/api
mv ./gen/index.ts ./src/api