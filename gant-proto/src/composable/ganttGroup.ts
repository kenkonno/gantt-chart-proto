import {Api} from "@/api/axios";
import {PostGanttGroupsRequest, GanttGroup} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useGanttGroupTable() {
    const list = ref<GanttGroup[]>([])
    const refresh = async () => {
        const resp = await Api.getGanttGroups()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useGanttGroup(ganttGroupId?: number) {

    const ganttGroup = ref<GanttGroup>({
        id: null,
        facility_id: 0,
        unit_id: 0,
        created_at: undefined,
        updated_at: undefined
    })
    if (ganttGroupId !== undefined) {
        const {data} = await Api.getGanttGroupsId(ganttGroupId)
        if (data.ganttGroup != undefined) {
            ganttGroup.value.id = data.ganttGroup.id
            ganttGroup.value.facility_id = data.ganttGroup.facility_id
            ganttGroup.value.unit_id = data.ganttGroup.unit_id
            ganttGroup.value.created_at = data.ganttGroup.created_at
            ganttGroup.value.updated_at = data.ganttGroup.updated_at
        }
    }

    return {ganttGroup}

}

export async function postGanttGroup(ganttGroup: GanttGroup, emit: any) {
    const req: PostGanttGroupsRequest = {
        ganttGroup: ganttGroup
    }
    await Api.postGanttGroups(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postGanttGroupById(ganttGroup: GanttGroup, emit: any) {
    const req: PostGanttGroupsRequest = {
        ganttGroup: ganttGroup
    }
    if (ganttGroup.id != null) {
        await Api.postGanttGroupsId(ganttGroup.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteGanttGroupById(id: number, emit: any) {
    await Api.deleteGanttGroupsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



