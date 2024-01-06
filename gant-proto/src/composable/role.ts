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
export type Section = "ALL_SETTINGS" | "FACILITY_SETTINGS" | "UPDATE_TICKET" | "UPDATE_PROGRESS" | "ADD_TICKET"


// 大枠の制御のセット。とりあえずつけただけで使わないものが多いかも。
const Manager = {
    "ALL_SETTINGS": true,
    "FACILITY_SETTINGS": true,
    "UPDATE_TICKET": true,
    "UPDATE_PROGRESS": true,
    "ADD_TICKET": true,
}
const Viewer = {
    "ALL_SETTINGS": false,
    "FACILITY_SETTINGS": false,
    "UPDATE_TICKET": false,
    "UPDATE_PROGRESS": false,
    "ADD_TICKET": false,

}
const Worker = {
    "ALL_SETTINGS": false,
    "FACILITY_SETTINGS": false,
    "UPDATE_TICKET": false,
    "UPDATE_PROGRESS": true,
    "ADD_TICKET": false,
}

export function allowed(section: Section) {
    const role = getUserInfo()?.role
    console.log(role)
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
    }
    return false
}