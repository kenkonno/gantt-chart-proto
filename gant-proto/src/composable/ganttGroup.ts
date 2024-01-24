import {Api} from "@/api/axios";
import {PostGanttGroupsRequest, GanttGroup} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {Emit} from "@/const/common";


// ユーザー一覧。特別ref系は必要ない。
export async function useGanttGroupTable() {
    const list = ref<GanttGroup[]>([])
    const refresh = async (facilityId: number) => {
        const resp = await Api.getGanttGroups(facilityId)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
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

export async function postGanttGroup(ganttGroup: GanttGroup, emit: Emit) {
    const req: PostGanttGroupsRequest = {
        ganttGroup: ganttGroup
    }
    await Api.postGanttGroups(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postGanttGroupById(ganttGroup: GanttGroup, emit: Emit) {
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

export async function deleteGanttGroupById(id: number, emit: Emit) {
    await Api.deleteGanttGroupsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



