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
}

export function allowed(section: Section) {
    const role = getUserInfo()?.role
    if (role === RoleType.Admin) {
        return true
    }
    switch (role) {
        case RoleType.Manager:
            return Manager[section]
        case RoleType.Worker:
            return Worker[section]
        case RoleType.Viewer:
            return Viewer[section]
        case RoleType.Guest:
            return Guest[section]
    }
    return false
}