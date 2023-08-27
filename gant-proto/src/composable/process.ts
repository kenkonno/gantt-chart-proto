import {Api} from "@/api/axios";
import {PostProcessesRequest, Process} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useProcessTable() {
    const list = ref<Process[]>([])
    const refresh = async () => {
        const resp = await Api.getProcesses()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
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
            process.value.created_at = data.process.created_at
            process.value.updated_at = data.process.updated_at
        }
    }

    return {process}

}

export async function postProcess(process: Process, emit: any) {
    const req: PostProcessesRequest = {
        process: process
    }
    await Api.postProcesses(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postProcessById(process: Process, emit: any) {
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

export async function deleteProcessById(id: number, emit: any) {
    await Api.deleteProcessesId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



