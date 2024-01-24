import {ref} from "vue";
import {allowed} from "@/composable/role";

export type GanttFacilityHeader = {
    name: string,
    visible: boolean
}

export type DisplayType = "day" | "week" | "hour" | "month"

export function useGanttFacilityMenu() {
    // injectはsetupと同期的に呼び出す必要あり
    const GanttHeader = ref<GanttFacilityHeader[]>([
        {name: "ユニット", visible: true},
        {name: "工程", visible: true},
        {name: "部署", visible: true},
        {name: "担当者", visible: true},
        {name: "人数", visible: false},
        {name: "期日", visible: false},
        {name: "工数(h)", visible: true},
        {name: "日後", visible: false},
        {name: "開始日", visible: false},
        {name: "終了日", visible: false},
        {name: "進捗", visible: true},
    ])
    if (allowed('UPDATE_TICKET')) {
        GanttHeader.value.push({name: "操作", visible: false})
    }
    const displayType = ref<DisplayType>("day")
    return {
        GanttHeader,
        displayType,
    }
}

