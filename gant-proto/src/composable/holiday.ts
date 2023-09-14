import {Api} from "@/api/axios";
import {PostHolidaysRequest, Holiday} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useHolidayTable(facilityId: number) {
    const list = ref<Holiday[]>([])
    const refresh = async () => {
        const resp = await Api.getHolidays(facilityId)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useHoliday(holidayId?: number) {

    const holiday = ref<Holiday>({
        id: null,
        facility_id: 0,
        name: "",
        date: "",
        created_at: undefined,
        updated_at: undefined
    })
    if (holidayId !== undefined) {
        const {data} = await Api.getHolidaysId(holidayId)
        if (data.holiday != undefined) {
            holiday.value.id = data.holiday.id
            holiday.value.facility_id = data.holiday.facility_id
            holiday.value.name = data.holiday.name
            holiday.value.date = data.holiday.date.substring(0,10)
            holiday.value.created_at = data.holiday.created_at
            holiday.value.updated_at = data.holiday.updated_at
        }
    }

    return {holiday}

}

export async function postHoliday(holiday: Holiday, facilityId: number, emit: Emit) {
    holiday.date = holiday.date + "T00:00:00.00000+09:00"
    holiday.facility_id = facilityId
    const req: PostHolidaysRequest = {
        holiday: holiday,
    }
    await Api.postHolidays(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postHolidayById(holiday: Holiday, facilityId: number, emit: Emit) {
    holiday.facility_id = facilityId
    const req: PostHolidaysRequest = {
        holiday: holiday,
    }
    if (holiday.id != null) {
        await Api.postHolidaysId(holiday.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteHolidayById(id: number, emit: Emit) {
    await Api.deleteHolidaysId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



