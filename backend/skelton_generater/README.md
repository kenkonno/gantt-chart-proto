# konno_cli

## 主な目的

---

model定義からapi-repositoryを一通で作成する。

openapi.yamlについてはシステムで吸収しない。運用で吸収する。（request, responseがややこしくなるため）

modelとresponseの紐づけは手動で実行する。

interactorのスケルトンまでは作成する。

## Outputs

 - repository
 - interactor

## Rules

- openapi.yamlの定義
  - APi名、リクエストモデル名、レスポンスモデル名が関連
    - interactorで名称を参照される
    - Api：[Method][ApiName]
    - Request：[Method][ApiName]Request
    - Response：[Method][ApiName]Response
  - 備考
    - yaml generatorはいい感じに interactorを呼び出してくれるだけ。システムでは考慮していない。

## How To Use

## command example

実行ディレクトリはbackend直下

```
go run ./main.go "../models/db/user"
```

構文
```
go run ./main.go "../models/db/[モデルへのパス(.goは除く)]"
```

注意事項

なぜかキャッシュクリアしないと依存関係が解決されないことがある模様。
分析のし直しでも問題ないかも。
