import {Api} from "@/api/axios";
import {FeatureOption, User} from "@/api";

// NOTE: composableに入れるのは違う気がするが一旦ここ

let userInfo: User | undefined
let isSimulateUser: boolean | undefined
let featureOptions: FeatureOption[] | undefined // auth?に入れるのは微妙かもだけどタイミングが同じなのでいれる。

export async function loggedIn() {
    const {data: userData} = await Api.getUserInfo()
    const {data} = await Api.getFeatureOptions()
    userInfo = userData.user
    isSimulateUser = userData.isSimulateUser
    featureOptions = data.list
    return userData
}

export function getUserInfo() {
    if (userInfo == undefined) {
        // TODO: 作りとしてよくないが、画面に遷移している以上はユーザー情報は取得されているはず。
        loggedIn()
    }
    return {userInfo, isSimulateUser}
}

export function getFeatureOptions() {
    if (featureOptions == undefined) {
        loggedIn()
    }
    return featureOptions
}