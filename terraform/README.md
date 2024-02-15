# やるぞTerraform

# アカウントを作ったら最初にやること
 - スイッチロールの作成
 - terraform操作できるフルアクセスのユーザーの作成

## Motivation

Samで何とかなると思ったけど寒い結果になって危機感を感じた。

ガントチャートプロジェクトでいよいよきちんと環境管理する必要が出てきた。

## 目標

- 一発激安環境を構築できるようにする
    - ECS on EC2 and EBS にする
    - cloudfront - s3
    - CodePipeline
- 実際に１か月動かさなかったときにどのくらいお金がかかるのか？（これは自分で払ってもいい）
- 開発環境・本番環境をデプロイできるようにする
- ChangeRoleをできるようにする
- terraformの学習

## 実行方法

Planの確認

``terraform.exe plan -var-file [epson-prod].terraform.tfvars``

適応

``terraform.exe apply -var-file [epson-prod].terraform.tfvars``

## TODO
- TODO: 実行フォルダを環境ごとに用意して管理する必要がある。
- ECRへのpushがdocker-compose.yamlに依存しているので環境ごとに手動となる。
- 


## メモ

- [簡単なチュートリアル](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/aws-build)

- 


## やること

- 今日の日記
  - どうやらswitch-roleの奴はそもそもアカウントが別れている模様。
  - アカウントを分けるとかなり便利になって、assume_roleでアカウント間を行き来できるようにするみたい。
  - 構成としては以下の通り。今のlaurensiaをマスターとして、dev,prod アカウントを作成してterraformの環境変数でデプロイ先を変更させるのがよさそう。
    - マスター組織
      - dev
      - stage
      - prod
  - terraformについてはなんとなく理解はした。
  - 次回はメールドレス(dev, prod)とアカウントの解説。クレジットカードはlaurensiaの奴を使う。

- 試しにS3をロール別に見れるようにする
  - リソースはアカウントに紐づく
  - 一番おおもとのアカウントを対象に assume_role を設定すると環境ごとが出来上がりそう
    - なんかちがった。安全の確保としてもそもそもアカウントを分けるのがベストプラクティスっぽい。
  - 一旦AWSでぽちぽちしてterraform上でサンプルを変更する試行錯誤で秘伝のたれを作るのがよさそう
  - なんかパッと見た感じフルオートカスタムって感じではなさそう。
  - Step.1でマスターに紐づくassume_roleを作って、その情報を手動で設定して各環境でdeployするのがわかりやすそう
    - 自分でもできるだけでプロなら別の方法があるかも。