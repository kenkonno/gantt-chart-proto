import {getUserInfo} from "@/composable/auth";
import {RoleType} from "@/const/common";

/**
 * RoleごとのTrue,Falseを返却していく。
 * セクション
 * ・全体設定
 * ・チケットの更新
 * ・進捗の更新
 * ・チケットの追加
 *
 */
export type Section =
    "ALL_SETTINGS" // 全体設定
    | "FACILITY_SETTINGS" // 設備設定
    | "UPDATE_TICKET" // チケットの更新
    | "UPDATE_PROGRESS" // 進捗の更新
    | "ADD_TICKET" // チケットの追加
    | "CHANGE_ROLE"// ロールの変更許可
    | "VIEW_PILEUPS" // 積み上げの表示
    | "VIEW_SCHEDULE_ALERT" // スケジュールアラーとの表示
    | "ALL_VIEW"    // 全体ビュー
    | "UPDATE_USER" // ユーザー編集の更新ボタンお押せるかどうか。ゲストユーザーように作成
    | "MENU" // Guest権限の時はメニュー系非表示
    | "FORCE_SIMULATE_USER" // シミュレーション操作ユーザー扱いにする。管理者のみ利用。
    | "CSV_UPLOAD"  // CSVアップロード全般
    | "UPDATE_MASTER" // 全体設定にかかるマスタ設定。管理者のみ利用可能とする。
    | "WORK_LOAD_WEIGHTING" // 重みづけ


// 大枠の制御のセット。とりあえずつけただけで使わないものが多いかも。
const Manager = {
    "ALL_SETTINGS": true,
    "FACILITY_SETTINGS": true,
    "UPDATE_TICKET": true,
    "UPDATE_PROGRESS": true,
    "ADD_TICKET": true,
    "CHANGE_ROLE": true,
    "VIEW_PILEUPS": true,
    "VIEW_SCHEDULE_ALERT": true,
    "ALL_VIEW": true,
    "UPDATE_USER": true,
    "MENU": true,
    "FORCE_SIMULATE_USER": false,
    "CSV_UPLOAD": true,
    "UPDATE_MASTER": false,
    "WORK_LOAD_WEIGHTING": true,
}
const Viewer = {
    "ALL_SETTINGS": false,
    "FACILITY_SETTINGS": false,
    "UPDATE_TICKET": false,
    "UPDATE_PROGRESS": false,
    "ADD_TICKET": false,
    "CHANGE_ROLE": false,
    "VIEW_PILEUPS": false,
    "VIEW_SCHEDULE_ALERT": true,
    "ALL_VIEW": true,
    "UPDATE_USER": true,
    "MENU": true,
    "FORCE_SIMULATE_USER": false,
    "CSV_UPLOAD": false,
    "UPDATE_MASTER": false,
    "WORK_LOAD_WEIGHTING": false,
}
const Worker = {
    "ALL_SETTINGS": false,
    "FACILITY_SETTINGS": false,
    "UPDATE_TICKET": false,
    "UPDATE_PROGRESS": true,
    "ADD_TICKET": false,
    "CHANGE_ROLE": false,
    "VIEW_PILEUPS": false,
    "VIEW_SCHEDULE_ALERT": true,
    "ALL_VIEW": true,
    "UPDATE_USER": true,
    "MENU": true,
    "FORCE_SIMULATE_USER": false,
    "CSV_UPLOAD": false,
    "UPDATE_MASTER": false,
    "WORK_LOAD_WEIGHTING": false,
}
const WorkerWithPileUps = {
    "ALL_SETTINGS": false,
    "FACILITY_SETTINGS": false,
    "UPDATE_TICKET": false,
    "UPDATE_PROGRESS": true,
    "ADD_TICKET": false,
    "CHANGE_ROLE": false,
    "VIEW_PILEUPS": true,
    "VIEW_SCHEDULE_ALERT": true,
    "ALL_VIEW": true,
    "UPDATE_USER": true,
    "MENU": true,
    "FORCE_SIMULATE_USER": false,
    "CSV_UPLOAD": false,
    "UPDATE_MASTER": false,
    "WORK_LOAD_WEIGHTING": false,
}
const Guest = {
    "ALL_SETTINGS": false,
    "FACILITY_SETTINGS": false,
    "UPDATE_TICKET": false,
    "UPDATE_PROGRESS": false,
    "ADD_TICKET": false,
    "CHANGE_ROLE": false,
    "VIEW_PILEUPS": false,
    "VIEW_SCHEDULE_ALERT": false,
    "ALL_VIEW": false,
    "UPDATE_USER": false,
    "MENU": false,
    "FORCE_SIMULATE_USER": false,
    "CSV_UPLOAD": false,
    "UPDATE_MASTER": false,
    "WORK_LOAD_WEIGHTING": false,
}

export function allowed(section: Section) {
    const {userInfo} = getUserInfo()
    if (userInfo?.role === RoleType.Admin) {
        return true
    }
    switch (userInfo?.role) {
        case RoleType.Manager:
            return Manager[section]
        case RoleType.Worker:
            return Worker[section]
        case RoleType.WorkerWithPileUps:
            return WorkerWithPileUps[section]
        case RoleType.Viewer:
            return Viewer[section]
        case RoleType.Guest:
            return Guest[section]
    }
    return false
}