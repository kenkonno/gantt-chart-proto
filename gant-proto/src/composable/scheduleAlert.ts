import {Api} from "@/api/axios";
import {ScheduleAlert} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useScheduleAlertTable() {
    const list = ref<ScheduleAlert[]>([])
    const refresh = async () => {
        const resp = await Api.getScheduleAlerts()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}
