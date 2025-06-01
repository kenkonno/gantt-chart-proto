import {Api} from "@/api/axios";
import {PostFacilityWorkSchedulesRequest, FacilityWorkSchedule} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useFacilityWorkScheduleTable(facilityId: number) {
    const list = ref<FacilityWorkSchedule[]>([])
    const refresh = async () => {
        const resp = await Api.getFacilityWorkSchedules(facilityId)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useFacilityWorkSchedule(facilityWorkScheduleId?: number) {

    const facilityWorkSchedule = ref<FacilityWorkSchedule>({
        id: null,
        facility_id: 0,
        date: "",
        type: "",
        created_at: undefined,
        updated_at: undefined
    })
    if (facilityWorkScheduleId !== undefined) {
        const {data} = await Api.getFacilityWorkSchedulesId(facilityWorkScheduleId)
        if (data.facilityWorkSchedule != undefined) {
            facilityWorkSchedule.value.id = data.facilityWorkSchedule.id
            facilityWorkSchedule.value.facility_id = data.facilityWorkSchedule.facility_id
            facilityWorkSchedule.value.date = data.facilityWorkSchedule.date.substring(0, 10)
            facilityWorkSchedule.value.type = data.facilityWorkSchedule.type
            facilityWorkSchedule.value.created_at = data.facilityWorkSchedule.created_at
            facilityWorkSchedule.value.updated_at = data.facilityWorkSchedule.updated_at
        }
    }

    return {facilityWorkSchedule}

}

export async function postFacilityWorkSchedule(facilityWorkSchedule: FacilityWorkSchedule, facilityId: number, emit: any) {
    facilityWorkSchedule.date = facilityWorkSchedule.date + "T00:00:00.00000+09:00"
    facilityWorkSchedule.facility_id = facilityId
    const req: PostFacilityWorkSchedulesRequest = {
        facilityWorkSchedule: facilityWorkSchedule
    }
    await Api.postFacilityWorkSchedules(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postFacilityWorkScheduleById(facilityWorkSchedule: FacilityWorkSchedule, facilityId: number, emit: any) {
    facilityWorkSchedule.date = facilityWorkSchedule.date + "T00:00:00.00000+09:00"
    facilityWorkSchedule.facility_id = facilityId
    const req: PostFacilityWorkSchedulesRequest = {
        facilityWorkSchedule: facilityWorkSchedule
    }
    if (facilityWorkSchedule.id != null) {
        await Api.postFacilityWorkSchedulesId(facilityWorkSchedule.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteFacilityWorkScheduleById(id: number, emit: any) {
    await Api.deleteFacilityWorkSchedulesId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



