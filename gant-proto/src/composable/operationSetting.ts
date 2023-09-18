import {Api} from "@/api/axios";
import {PostOperationSettingsRequest, OperationSetting} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー追加・更新。
export async function useOperationSettingTable(facilityId: number) {

    const list = ref<OperationSetting[]>([])
    const refresh = async () => {
        const resp = await Api.getOperationSettingsId(facilityId)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.operationSettings)
    }
    await refresh()
    return {list, refresh}
}

export async function postOperationSettingById(facilityId: number, operationSettings: OperationSetting[], emit: Emit) {
    const req: PostOperationSettingsRequest = {
        operationSettings: operationSettings
    }
    await Api.postOperationSettingsId(facilityId, req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function deleteOperationSettingById(id: number, emit: Emit) {
    await Api.deleteOperationSettingsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



