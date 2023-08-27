import {Api} from "@/api/axios";
import {PostDepartmentsRequest, Department} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useDepartmentTable() {
    const list = ref<Department[]>([])
    const refresh = async () => {
        const resp = await Api.getDepartments()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useDepartment(departmentId?: number) {

    const department = ref<Department>({
        id: null,
        name: "",
        created_at: undefined,
        updated_at: undefined
    })
    if (departmentId !== undefined) {
        const {data} = await Api.getDepartmentsId(departmentId)
        if (data.department != undefined) {
            department.value.id = data.department.id
            department.value.name = data.department.name
            department.value.created_at = data.department.created_at
            department.value.updated_at = data.department.updated_at
        }
    }

    return {department}

}

export async function postDepartment(department: Department, emit: any) {
    const req: PostDepartmentsRequest = {
        department: department
    }
    await Api.postDepartments(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postDepartmentById(department: Department, emit: any) {
    const req: PostDepartmentsRequest = {
        department: department
    }
    if (department.id != null) {
        await Api.postDepartmentsId(department.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteDepartmentById(id: number, emit: any) {
    await Api.deleteDepartmentsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



