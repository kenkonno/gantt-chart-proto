import {Api} from "@/api/axios";
import {PostSimulationRequest, PutSimulationRequest, SimulationLock} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {changeSort} from "@/utils/sort";
import {Emit} from "@/const/common";

// ユーザー追加・更新。
export async function useSimulation() {

    const simulationLock = ref<SimulationLock>({
        simulationName: "",
        lockedAt: "",
        lockedBy: 0,
        status: ""
    })
    const {data} = await Api.getSimulation()
    if (data.simulationLock != undefined) {
        simulationLock.value.simulationName = data.simulationLock.simulationName
        simulationLock.value.lockedAt = data.simulationLock.lockedAt
        simulationLock.value.lockedBy = data.simulationLock.lockedBy
        simulationLock.value.status = data.simulationLock.status
    }

    return {simulationLock}

}

export async function postSimulation(emit: Emit) {
    const req: PostSimulationRequest = {}
    await Api.postSimulation(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export type SimulationMode = "pending" | "resume" | "apply"

export async function putSimulation(mode: SimulationMode, emit: Emit) {
    const req: PutSimulationRequest = {
        mode: mode
    }
    await Api.putSimulation(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function deleteSimulation(emit: Emit) {
    await Api.deleteSimulation().then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



