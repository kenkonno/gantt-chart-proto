質問や計画は必ず日本語で出力してください。

# モデルからAPIを生成する方法

1. backend/models/db配下のファイルはgormで利用するORマッパーの定義群です。
2. 指示で指定された定義をもとに以下の作業を行います。
   1. backend/api/GanttChartApi.yaml の更新を行う
      1. 一般的な GET/POST/DELETE/PUT を追加すること
   2. backend/api/api_openapi_generator.sh の実行を行う
   3. backend/api/interactorに追加を行う
      1. フォルダ名は model名のsnake_caseの複数形とすること（例：feature_option → feature_options）
      2. APIの名前は [METHOD]_[MODELのsnake_case（複数形）]
      3. APIの内容は一般的なシンプルなCRUDの制御で良い。
   4. backend/repositoryの追加
      1. repository配下にはfactory関数を追加する
      2. interfacesにはインターフェースを追加する
      3. commonの配下には実装を追加する
      4. simulation_modeでも利用するモデルの場合はsimulation配下にも実装を追加する
         1. 全てのテーブル名に「simulation_」を付与すること
         2. backend/models/db/配下にも 「simulation_[モデル名]」のファイルを追加すること
            1. これの実装方法は元のmodelを利用し個別の定義はしないこと。
      5. backend/migration/main.go に 追加したモデルのmigrationコードを実装すること
         1. 例：
   ```go
   // 通常のモデル
   migrate(db.FeatureOption{})
   
   // simulationモデル
   migrate(db.SimulationFeatureOption{})
   ```

**一連の作業は既存のAPIのファイルを参考にして実装を進めること。例えば holidaysに関するファイルを参照してください**

**特に lo のライブラリを利用すること**