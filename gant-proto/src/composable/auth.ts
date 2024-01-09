import {Api} from "@/api/axios";
import {User} from "@/api";

// NOTE: composableに入れるのは違う気がするが一旦ここ

let userInfo: User | undefined

export async function loggedIn() {
    // TODO: COOKIE NAME
    const {data} = await Api.getUserInfo()
    userInfo = data.user
    return data
}

export function getUserInfo() {
    if (userInfo == undefined) {
        // TODO: 作りとしてよくないが、画面に遷移している以上はユーザー情報は取得されているはず。
        loggedIn()
    }
    return userInfo
}