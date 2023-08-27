import {Api} from "@/api/axios";
import {PostUsersRequest, User} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useUserTable() {
    const list = ref<User[]>([])
    const refresh = async () => {
        const resp = await Api.getUsers()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useUser(userId?: number) {

    const user = ref<User>({
        id: null,
        name: "",
        password: "",
        email: "",
        created_at: undefined,
        updated_at: undefined
    })
    if (userId !== undefined) {
        const {data} = await Api.getUsersId(userId)
        if (data.user != undefined) {
            user.value.id = data.user.id
            user.value.name = data.user.name
            user.value.password = data.user.password
            user.value.email = data.user.email
            user.value.created_at = data.user.created_at
            user.value.updated_at = data.user.updated_at
        }
    }

    return {user}

}

export async function postUser(user: User, emit: any) {
    const req: PostUsersRequest = {
        user: user
    }
    await Api.postUsers(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postUserById(user: User, emit: any) {
    const req: PostUsersRequest = {
        user: user
    }
    if (user.id != null) {
        await Api.postUsersId(user.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteUserById(id: number, emit: any) {
    await Api.deleteUsersId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



