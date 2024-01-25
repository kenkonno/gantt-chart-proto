import {onBeforeUnmount, onUnmounted, ref} from "vue";
import {allowed} from "@/composable/role";
import {globalFilterGetter, globalFilterMutation} from "@/utils/globalFilterState";

export type GanttFacilityHeader = {
    name: string,
    visible: boolean
}

export type DisplayType = "day" | "week" | "hour" | "month"

export function useGanttFacilityMenu() {

    // 初期値の取得
    const savedGanttHeader = globalFilterGetter.getGanttFacilityMenu()
    const savedViewType = globalFilterGetter.getViewType()

    // injectはsetupと同期的に呼び出す必要あり
    const GanttHeader = ref<GanttFacilityHeader[]>(savedGanttHeader)
    if (allowed('UPDATE_TICKET')) {
        if (!GanttHeader.value.find(v => v.name == "操作")) {
            GanttHeader.value.push({name: "操作", visible: false})
        }
    } else {
        const index = GanttHeader.value.findIndex(v => v.name == "操作")
        GanttHeader.value.splice(index)
    }
    const displayType = ref<DisplayType>(savedViewType)

    // フィルタ保存関連
    const safeFilter = () => {
        globalFilterMutation.updateGanttFacilityMenu(GanttHeader.value)
        globalFilterMutation.updateViewTypeFilter(displayType.value)
    }
    onBeforeUnmount(() => {
        safeFilter()
        window.removeEventListener("beforeunload", safeFilter)
    })
    window.addEventListener("beforeunload", safeFilter)
    return {
        GanttHeader,
        displayType,
    }
}
