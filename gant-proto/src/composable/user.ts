import {Api} from "@/api/axios";
import {PostUploadUsersCsvFileRequest, PostUsersRequest, User} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {Emit, RoleType} from "@/const/common";
import Swal from "sweetalert2";


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
        lastName: "",
        firstName: "",
        department_id: 0,
        limit_of_operation: 8.0,
        password: "",
        email: "",
        role: RoleType.Viewer,
        created_at: undefined,
        updated_at: undefined,
        password_reset: false,
        employment_start_date: "",
        employment_end_date: undefined,

    })
    if (userId !== undefined) {
        const {data} = await Api.getUsersId(userId)
        if (data.user != undefined) {
            user.value.id = data.user.id
            user.value.department_id = data.user.department_id
            user.value.lastName = data.user.lastName
            user.value.firstName = data.user.firstName
            user.value.limit_of_operation = data.user.limit_of_operation
            user.value.password = data.user.password
            user.value.email = data.user.email
            user.value.role = data.user.role
            user.value.created_at = data.user.created_at
            user.value.updated_at = data.user.updated_at
            user.value.password_reset = data.user.password_reset
            user.value.employment_start_date = data.user.employment_start_date.substring(0, 10)
            user.value.employment_end_date = data.user.employment_end_date?.substring(0, 10)
        }
    }

    return {user}

}

export function validate(user: User, validatePassword = false) {
    let isValid = true
    if (!user.firstName) {
        toast.warning("姓は必須です")
        isValid = false
    }
    if (!user.lastName) {
        toast.warning("名は必須です")
        isValid = false
    }
    if (!user.password && validatePassword) {
        toast.warning("Passwordは必須です")
        isValid = false
    }
    if (!user.email) {
        toast.warning("Emailは必須です")
        isValid = false
    }
    if (!user.email) {
        toast.warning("在籍期間(開始)は必須です")
        isValid = false
    }
    return isValid
}

export async function confirm(user: User) {
    const {data} = await Api.getDetectWorkOutsideEmploymentPeriods(user.id!, user.employment_start_date, user.employment_end_date || "")
    if (data.list.length > 0) {
        const result = await Swal.fire({
            title: '在籍期間外になる工程が存在します。',
            text: "在籍期間外になる工程からは担当者から除外されますが更新をしますか？",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: '実行',
            cancelButtonText: 'キャンセル'
        });
        return result.isConfirmed
    }
    return true
}


export async function postUser(user: User, emit: Emit) {
    user.employment_start_date = user.employment_start_date + "T00:00:00.00000+09:00"
    if (user.employment_end_date) {
        user.employment_end_date = user.employment_end_date + "T00:00:00.00000+09:00"
    } else {
        user.employment_end_date = undefined
    }
    const req: PostUsersRequest = {
        user: user
    }

    await Api.postUsers(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postUserById(user: User, emit: Emit) {
    const clone = {...user}
    clone.employment_start_date = clone.employment_start_date + "T00:00:00.00000+09:00"
    if (clone.employment_end_date) {
        clone.employment_end_date = clone.employment_end_date + "T00:00:00.00000+09:00"
    } else {
        clone.employment_end_date = undefined
    }
    const req: PostUsersRequest = {
        user: clone
    }

    // 在籍期間外の確認
    if (!await confirm(clone)) {
        return
    }

    if (clone.id != null) {
        await Api.postUsersId(clone.id, req).then(() => {
            toast("成功しました。")
            emit('closeEditModal')
        })
    }
}

export async function deleteUserById(id: number, emit: Emit) {
    await Api.deleteUsersId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}


export async function postUploadUsersCsvFile(csv: File) {
    console.log(csv)

    await Api.postUploadUsersCsvFile(csv, {
        headers: {
            'Content-Type': 'application/octet-stream',
        }
    }).then(() => {
        toast("成功しました。")
    })
}
