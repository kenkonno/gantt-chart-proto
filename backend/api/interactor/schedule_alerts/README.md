# 遅延通知

仕様が複雑なのできちんとまとめる。
処理ロジック自体はSQL一本で実行している。

# 遅延とは
進捗との乖離と定める。

完了している日付 = 営業日数 * 進捗
遅延日数 = 現在の日付 - 完了している日付

工数はとか人数とか関係なく営業日と進捗だけで計算する


OLD

------------

# 遅延とは
消化しきれていない工数を算出することと定める。

消化済み工数は 工数 * 進捗% で算出

未消化は 工数 - 消化済み工数 で算出

# 例
条件(通常の40h作業)
- 期間5日
- 工数40h
- 進捗50%
- 開始日：10日
- 終了日：14日
- 計測日：17日（次の月曜日）

結果
- 消化済み工数：20h
- 実際の進捗：11日（16h分は完了しているが、24h分は完了していないため）
- 遅延日数：3日

条件(薄く伸ばした20h作業)
- 期間5日
- 工数40h
- 進捗50%
- 開始日：10日
- 終了日：14日
- 計測日：17日（次の月曜日）

結果
- 消化済み工数：20h
- 実際の進捗：11日（16h分は完了しているが、24h分は完了していないため）
- 遅延日数：3日



# 進捗からの完了している日付の算出
色々考えたけど、単に遅れている日数が出たほうが直観的であるので、営業日換算で日数計算したほうがいいのか。
