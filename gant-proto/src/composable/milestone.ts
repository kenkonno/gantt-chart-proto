import {Api} from "@/api/axios";
import {PostMilestonesRequest, Milestone} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useMilestoneTable(facilityId: number, mode?: string) {
    const list = ref<Milestone[]>([])
    const refresh = async () => {
        const resp = await Api.getMilestones(facilityId, mode)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useMilestone(facilityId: number, milestoneId?: number) {

    const milestone = ref<Milestone>({
        id: null,
        facility_id: facilityId,
        date: "",
        description: "",
        order: 0,
        created_at: undefined,
        updated_at: undefined
    })
    if (milestoneId !== undefined) {
        const {data} = await Api.getMilestonesId(milestoneId)
        if (data.milestone != undefined) {
            milestone.value.id = data.milestone.id
            milestone.value.facility_id = data.milestone.facility_id
            milestone.value.date = data.milestone.date.substring(0, 10)
            milestone.value.description = data.milestone.description
            milestone.value.order = data.milestone.order
            milestone.value.created_at = data.milestone.created_at
            milestone.value.updated_at = data.milestone.updated_at
        }
    }

    return {milestone}

}
export function validate(milestone: Milestone) {
    let isValid = true
    if (!milestone.description) {
        toast.warning("説明は必須です")
        isValid = false
    }
    if (!milestone.date) {
        toast.warning("日付は必須です")
        isValid = false
    }
    return isValid
}
export async function postMilestone(milestone: Milestone, order: number, emit: any) {
    milestone.date = milestone.date + "T00:00:00.00000+09:00"
    milestone.order = order
    const req: PostMilestonesRequest = {
        milestone: milestone
    }
    await Api.postMilestones(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
        emit('update')
    })
}

export async function postMilestoneById(milestone: Milestone, emit: any) {
    milestone.date = milestone.date + "T00:00:00.00000+09:00"
    const req: PostMilestonesRequest = {
        milestone: milestone
    }
    if (milestone.id != null) {
        await Api.postMilestonesId(milestone.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
            emit('update')
        })
    }
}

export async function deleteMilestoneById(id: number, emit: any) {
    await Api.deleteMilestonesId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
        emit('update')
    })
}



