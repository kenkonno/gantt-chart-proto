# 設備に関するInteractor

特筆すべきはステータスによって絞り込みをすること。

Status：有効・無効・完了

Type：受注済み・非受注

チケット取得に関連するものは 無効・完了は除外する。

ビューの選択が受注済みの場合は受注済みのみ、非受注の場合は受注済みも含む

- 設備一覧
  - すべて取得
- 設備のドロップダウン
  - 無効を除く、Typeは状態に従う
    - APIとして一覧の者と同じなのでfilterで対応する
- 設備を選択したときの状態
  - ID指定なので特に気にしなくてよい気がする
- 山積みの場合
  - 無効・完了を除く、Typeは状態に従う
  - all-ticketsAPIに修正をする
    - facility-type, facility-statusの条件を追加
  - ticket-usersはall-ticketに紐づくものを取得しているので気にしなくてOK
  - pieUpsApiを修正する（こんなの作ってたんだ）
    - 同じく facility-type, facility-status の条件を追加する
