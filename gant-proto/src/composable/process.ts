import {Api} from "@/api/axios";
import {PostProcessesRequest, Process} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {changeSort} from "@/utils/sort";


// ユーザー一覧。特別ref系は必要ない。
export async function useProcessTable() {
    const list = ref<Process[]>([])
    const refresh = async () => {
        const resp = await Api.getProcesses()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    const updateOrder = async (index: number, direction: number) => {
        changeSort(list.value, index, direction)

        for (const v of list.value) {
            v.order = list.value.indexOf(v)
            // API直呼び出しは少し気持ち悪いが効率を考慮してこうする。
            await Api.postProcessesId(v.id!, {process: v})
        }
    }

    return {list, refresh, updateOrder}
}

// ユーザー追加・更新。
export async function useProcess(processId?: number) {

    const process = ref<Process>({
        id: null,
        name: "",
        created_at: undefined,
        updated_at: undefined
    })
    if (processId !== undefined) {
        const {data} = await Api.getProcessesId(processId)
        if (data.process != undefined) {
            process.value.id = data.process.id
            process.value.name = data.process.name
            process.value.order = data.process.order
            process.value.created_at = data.process.created_at
            process.value.updated_at = data.process.updated_at
        }
    }

    return {process}

}

export async function postProcess(process: Process, order: number,emit: Emit) {
    process.order = order
    const req: PostProcessesRequest = {
        process: process
    }
    await Api.postProcesses(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postProcessById(process: Process, emit: Emit) {
    const req: PostProcessesRequest = {
        process: process
    }
    if (process.id != null) {
        await Api.postProcessesId(process.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteProcessById(id: number, emit: Emit) {
    await Api.deleteProcessesId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



