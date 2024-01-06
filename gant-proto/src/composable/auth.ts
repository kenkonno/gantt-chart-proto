import {Api} from "@/api/axios";

// NOTE: composableに入れるのは違う気がするが一旦ここ

export async function loggedIn() {
    // TODO: COOKIE NAME
    const {data} = await Api.getUserInfo()
    return data
}