import {onBeforeUnmount, ref} from "vue";
import {globalFilterGetter, globalFilterMutation} from "@/utils/globalFilterState";
import {DisplayType} from "@/composable/ganttFacilityMenu";

export type Header = {
    name: string,
    visible: boolean
}

export function useGanttAllMenu() {

    const savedGanttAllMenu = globalFilterGetter.getGanttAllMenu()
    const savedViewType = globalFilterGetter.getViewType()

    const displayType = ref<DisplayType>(savedViewType)

    const GanttHeader = ref<Header[]>(savedGanttAllMenu)
    const safeFilter = () => {
        globalFilterMutation.updateGanttAllMenu(GanttHeader.value)
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


