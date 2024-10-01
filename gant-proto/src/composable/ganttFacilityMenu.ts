import {onBeforeUnmount, ref} from "vue";
import {allowed} from "@/composable/role";
import {globalFilterGetter, globalFilterMutation} from "@/utils/globalFilterState";

export type GanttFacilityHeader = {
    name: string,
    visible: boolean
}

export type DisplayType = "day" | "week" | "hour" | "month"
export type AggregationAxis = "process" | "facility"

export function useGanttFacilityMenu() {

    // 初期値の取得
    const savedGanttHeader = globalFilterGetter.getGanttFacilityMenu()
    const savedViewType = globalFilterGetter.getViewType()

    // injectはsetupと同期的に呼び出す必要あり
    const GanttHeader = ref<GanttFacilityHeader[]>(savedGanttHeader)
    if (!allowed('UPDATE_TICKET')) {
        const index = GanttHeader.value.findIndex(v => v.name == "操作")
        if (index >= 0) {
            GanttHeader.value.splice(index)
        }
    } else {
        // NOTE: 権限によってstorageから消したり増えたりするものをきちんと対応する必要がある。今回は全て末尾なのでこれでOK
        const index = GanttHeader.value.findIndex(v => v.name == "操作")
        if (index === -1) {
            GanttHeader.value.push({name: "操作", visible: false})
        }
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
