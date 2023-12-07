import {ref} from "vue";

export type Header = {
    name: string,
    visible: boolean
}

export type DisplayType = "day" | "week" | "hour" | "month"

export function useGanttAllMenu() {

    const displayType = ref<DisplayType>("week")

    const GanttHeader = ref<Header[]>([
        {name: "設備名", visible: true},
        {name: "担当者", visible: false},
        {name: "開始日", visible: true},
        {name: "終了日", visible: true},
        {name: "工数(h)", visible: true},
        {name: "進捗", visible: true},
    ])

    return {
        GanttHeader,
        displayType,
    }
}


