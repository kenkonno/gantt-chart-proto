import {Api} from "@/api/axios";
import {getUserInfo} from "@/composable/auth";
import {computed} from "vue";

export const useIsSimulate = async () => {
    // NOTE: 正しい動きをしてるけど条件分岐が複雑。いいアイデアが思いついたらリファクタする。
    const {data} = await Api.getSimulation()
    const userInfo = getUserInfo()
    const simulateUser = userInfo.isSimulateUser


    const isSimulate = computed(() => {
        // 非シミュレーションユーザーは操作可能にする
        if (!simulateUser) return false

        // シミュレーション中かで判定する
        return data.simulationLock.status == 'in_progress';

    })

    const isViewOnly = computed(() => {
        // シミュレーション中であればシミュレーションユーザーのみ編集可能
        return isSimulate.value;
    })


    const isSimulateUser = computed(() => {
        // シミュレーション操作者ならば警告を表示しない
        if (simulateUser) {
            return false
        }
        // シミュレーション中かで判定する
        return data.simulationLock.status == 'in_progress';

    })
    const isViewOnlyUser = computed(() => {
        // シミュレーション中で無ければ編集可能
        if (!isSimulateUser.value) {
            return false
        }
        // シミュレーション中であればシミュレーションユーザーのみ編集可能
        return !simulateUser;

    })


    return {
        isSimulate,
        isViewOnly,
        isSimulateUser,
        isViewOnlyUser
    }

}

