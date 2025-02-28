import {Api} from "@/api/axios";
import {PostSimulationRequest, PutSimulationRequest, SimulationLock} from "@/api";
import {ref} from "vue";
import {Emit} from "@/const/common";
import Swal from "sweetalert2";

// ユーザー追加・更新。
export async function useSimulation() {

    const simulationLock = ref<SimulationLock>({
        simulationName: "",
        lockedAt: "",
        lockedBy: 0,
        status: ""
    })

    const refresh = async () => {
        const {data} = await Api.getSimulation()
        if (data.simulationLock != undefined) {
            simulationLock.value.simulationName = data.simulationLock.simulationName
            simulationLock.value.lockedAt = data.simulationLock.lockedAt
            simulationLock.value.lockedBy = data.simulationLock.lockedBy
            simulationLock.value.status = data.simulationLock.status
        }
    }
    await refresh()

    return {simulationLock, refresh}

}

export async function postSimulation(emit: Emit) {
    const req: PostSimulationRequest = {}
    await Api.postSimulation(req).then(() => {
        Swal.fire({
            title: '成功しました。画面をリロードします。',
            icon: 'success',
            showCancelButton: false,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: '閉じる',
        }).then(() => {
            window.location.reload()
        })
    }).finally(() => {
        emit('update')
    })
}

export type SimulationMode = "pending" | "resume" | "apply"

export async function putSimulation(mode: SimulationMode, emit: Emit) {
    const req: PutSimulationRequest = {
        mode: mode
    }
    await Api.putSimulation(req).then(() => {
        Swal.fire({
            title: '成功しました。画面をリロードします。',
            icon: 'success',
            showCancelButton: false,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: '閉じる',
        }).then(() => {
            window.location.reload()
        })
    }).finally(() => {
        emit('update')
    })
}

export async function deleteSimulation(emit: Emit) {
    await Api.deleteSimulation().then(() => {
        Swal.fire({
            title: '成功しました。画面をリロードします。',
            icon: 'success',
            showCancelButton: false,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: '閉じる',
        }).then(() => {
            window.location.reload()
        })
    }).finally(() => {
        emit('update')
    })
}



