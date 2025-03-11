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


## 会社別環境構築手順

1. メールアドレスの新規作成
   1. https://secure.sakura.ad.jp/auth/login?url=https://secure.sakura.ad.jp/menu/top/index.php
   2. fouから始まるIDでログイン
   3. 契約中のサービスをクリック
   4. laurensiaのメールサーバーを選択
   5. メール一覧
   6. 新規追加
   7. アカウント情報管理シートに追加
   8. https://docs.google.com/spreadsheets/d/18YMkYk1Egan787MOglTR3o5UchJEg5U7K4GpCDNT_QU/edit?gid=0#gid=0
   9. 作成後webメールを開いておく
2. awsのアカウント作成する
   1. 新規作成ページ https://signin.aws.amazon.com/signup?request_type=register
   2. パスワードはメールアドレスと同じにする
   3. カード情報は以下を参照
   4. https://docs.google.com/document/d/1WYQ7CDRtyH-w8gIjfG6oNng4_BzbLGEKSQnVa9i_G74/edit?tab=t.0
3. terraform用とスイッチロール用のIAMを作成する（スイッチロール用）
   1. IAMからユーザーの作成でAdministrator権限で作成する。
   2. アクセスキー、シークレットキーを作成する
   3. 後のtfvarsにコピペする。
   4. ロールの追加を行う
   5. 信頼されたエンティティの選択でAWSアカウントを選択
   6. アカウントIDは 866026585491 (dev.laurensia) を選択
   7. 管理者権限を付与する。
   8. ロールの名前とAWSのアカウントIDが必要なのでメモしておくこと。
   9. 移行 system.laurensiaのアカウントからスイッチロールして作業すること。
4. キーペアを作成する
   1. EC2を開く
   2. キーペアを選択
   3. pemで作成すること
   4. 以下のドライブに保存する
   5. https://drive.google.com/drive/u/0/folders/1TpmLlUG6-yAx2vtLqRnOCWUdBk4sdVt1
   6. 自分の.sshフォルダにも保存する
5. プライベートVPCを作成する
   1. 時間がないので今回は省略
   2. VPCID, セキュリティグループID, サブネットID(今回はa,cを利用)
6. tfvarsを作成する
   1. 最新のmds-prod.terraform.tfvarsをコピーする
7. terraformの実行
   1. lockファイルの terraform.tfstateをリネームする
   2. terraform.exe plan -var-file ftech-prod.terraform.tfvars
   3. terraform.exe apply -var-file ftech-prod.terraform.tfvars
8. 移行は simple-ecs-on-ec2のREADME.md を参考に作業をすること。