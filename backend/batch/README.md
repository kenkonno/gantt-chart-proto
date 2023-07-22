# Batch

--- 
 - up postgresql
   - docker-compose up -d postgres
  - build batch
   - docker-compose run batch bash -c "cd batch && go build" --rm batch
     - ビルド
   - docker-compose run batch bash -c "cd batch && ./batch MigrationV1" --rm batch
   - docker-compose run batch bash -c "cd batch && ./batch AddStreamer [loginNames]" --rm batch
   - docker-compose run batch bash -c "cd batch && ./batch AddPreferredStreamerExecute [loginNames] [gameID]" --rm batch
     - 優先ストリーマーの追加。雑談：509658 p_tilda
   - docker-compose run batch bash -c "cd batch && ./batch UpdateFamousGameByTwitchTop" --rm batch
     - twitchの人気ゲームから有名ゲームフラグの更新（有名ゲームフラグはクリップの取得に使用）
   - docker-compose run batch bash -c "cd batch && ./batch UpdateGamesByFamousStreamers" --rm batch
     - ストリーマーから有名ゲームフラグの更新
   - docker-compose run batch bash -c "cd batch && ./batch UpdateClipByFamousGames" --rm batch
     - クリップの更新 
   - docker-compose run batch bash -c "cd batch && ./batch UpdateClipByPreferredStreamers" --rm batch
     - 優先ストリーマーのクリップ更新
   - docker-compose run batch bash -c "cd batch && ./batch GameClipReport" --rm batch
     - クリップレポートの更新
   - docker-compose run batch bash -c "cd batch && ./batch GameClipReportPreferredStreamers" --rm batch
     - 優先ストリーマーのレポート更新
   - docker-compose run batch bash -c "cd batch && ./batch UpdateStreamers" --rm batch
     - streamerの情報を更新する
--- 

Twitch Helix Api を利用してひたすらデータを収集する

# Commands

--- 

update_clip_by_famous_games
 - clips
 - games
 - streamers
 - videos

AddStreamer [login]

# Tips

 - 時間の区切り
   - 朝 05:00 - 18:00
   - 夜 18:00 - 05:00

 - DIVISION
   - [1|2|3][0|1] 
   - 2桁目：Daily, Weekly, 優先ストリーマー
   - 1桁目：朝・夜

-- エロ系配信者を集めるバッチも作りたい でもclipがおわってるんだよなー

# DB Prod to local
1. localのTABLEはTRUNCATEしておく（超注意ssh_tunnnelしていないこと
2. data gripのpg_dumpを使う（注意：ssh_tunnelはdatagrip側の設定でしておくこと
3. 右クリックする箇所はpublic

# V2の考え方

V1でダメだった点の振り返り

 - master_updaterが終わってた
 - twitchの仕様的に結局clipだけでよかった

課題を解決させるには

 - videoの検索はしない
 - clipの検索だけにする
 - streamerの検索はほどほどにする

具体的には

game と streamer の選別が重要
もはや手動メンテでいいのでは？ いいAPIないかな
IRCチャット監視すればいける？いったん手動でいい気がしてきた

なんかgameごとにclipとってくるでいい気がしてきた。マイナー配信者には悪いけどそんなもんでしょうが。
クリップ自体はたくさんとってきて同じ配信者は沢山出さないようにすればいいかも。
gamesからclipを１００件取得するようにする
gamesの追加のトリガーが発生しないから、有名ストリーマーだけで新しいgameがないかマスターを作るか。

 - gamesからclipsの取得
 - 有名ストリーマーからgamesの取得
 - streamersにフラグを追加する

実装の流れ

 - games から clips の取得
 - clips から存在しないstreamersがいたらstreamersの更新
 - 有名ストリーマーからgamesの更新

問題を分けましょう


 - 日付指定でclipをとってくるバッチを作る
 - clipをとってくるgameマスタを更新する (追加した)
   - gameマスタを手動で指定する (追加した)
 - gameマスタの更新
   - 有名ストリーマーからフラグ指定で設定する (追加した)

つまりストリーマーから始まる物語
フラグの更新ロジックは一旦後で考えよう

Get Top Games Apiでもしかして全部解決したのかも・・・
まあ有名ストリーマーが配信してたら有名なゲームでしょ

ふと思ったけどストリーマーを更新するバッチがないと意味ないね。
まあ一旦、初期データとして投入するか
と思ったけど clip を更新するときに存在しないストリーマーは追加するからそれでいいのかな

有名ゲームからclipsをとってくる
 - ストリーマーもおのずとくっついてくる
 - 月一くらいでそのゲームジャンルのTOP20くらいで足きりすればいい感じになるのでは？


バッチの実行順を考える

 - UpdateFamousGamesByTwitchTop
   - 世界的に人気なゲームの決定をする（このを実行する前に全部famous=falseにする）
 - UpdateGamesByFamousStreamers
   - 取りこぼしがないように日本の人気ストリーマーからgamesを更新する
 - UpdateClipByFamousGames（05:00～23:59日の間に実行、引数なしだと１日前を指定する
   - 人気ゲームをもとにClipを更新する
ということは、ゲーム→ストリーマー→clipとなるので自動でいけそう。
月次バッチで先月のクリップから人気ストリーマーを更新するのも自動でいける？まあ後で考えよう

タイムライン
 - gamesの更新（朝５時）
 - Clipsだけの更新（夕方５時）

TODO:本番からデータコピーする仕組みを作る（開発にずらすぎるから

・明日はAPIがモデル変えたのでちゃんと動くか確認
・famous streamers の設定変更
・フロントを完全にclip型に変更する
　・この辺からデザイン回りで色々迷うかもね。

・個別検索とエロ系ドッチ先？
　・たぶんエロ系、個別検索したい衝動にかられないから
　・優先ストリーマーとエロ系機能を先に作る（クリップレポートに無理やり乗せる方法を検討する）
　・優先ストリーマーは作った

・お気に入りのクリップを探すを作って、今日の人気クリップとどうやって分けるかを検討する。
