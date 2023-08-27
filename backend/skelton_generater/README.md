# konno_cli

## 主な目的

---

model定義からapi-repositoryを一通で作成する。

このcliの目的は、新規作成のファイル作成の補助が目的。

常にこのCLIを通して自動生成したファイルで最新かすることではない。

openapi_models.yamlについてはシステムで吸収しない。運用で吸収する。（request, responseがややこしくなるため）

modelとresponseの紐づけは手動で実行する。

interactorのスケルトンまでは作成する。

## Outputs

 - repository
 - interactor

## Rules

- openapi_models.yamlの定義
  - APi名、リクエストモデル名、レスポンスモデル名が関連
    - interactorで名称を参照される
    - Api：[Method][ApiName]
    - Request：[Method][ApiName]Request
    - Response：[Method][ApiName]Response
  - 備考
    - yaml generatorはいい感じに interactorを呼び出してくれるだけ。システムでは考慮していない。
    - がいい感じにスケルトンは作る。


メモ
一旦作り切ったので、今後テンプレートをエンハンスしていけばすごくよくなると思います。


## How To Use

## command example
```
go run ./main.go "../models/db/user"
```

構文
```
go run ./main.go "../models/db/[モデルへのパス(.goは除く)]"
```


1. 実行ディレクトリはskelton_generater直下
2. 必要があれば Migrationをする
3. destディレクトリにあるファイル群をコピー
4. 必要があれば dest/yaml_info.txt を api.yamlにコピペ
   1. Apiのパス定義と モデル定義の２つに別れています。
5. openapi_generator.shの実行（バックエンド、フロント）


注意事項

なぜかキャッシュクリアしないと依存関係が解決されないことがある模様。
分析のし直しでも問題ないかも。
