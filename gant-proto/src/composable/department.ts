import {Api} from "@/api/axios";
import {PostDepartmentsRequest, Department} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {changeSort} from "@/utils/sort";
import {Emit} from "@/const/common";


// ユーザー一覧。特別ref系は必要ない。
export async function useDepartmentTable() {
    const list = ref<Department[]>([])
    const refresh = async () => {
        const resp = await Api.getDepartments()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()

    const updateOrder = async (index: number, direction: number) => {
        changeSort(list.value, index, direction)

        for (const v of list.value) {
            v.order = list.value.indexOf(v)
            // API直呼び出しは少し気持ち悪いが効率を考慮してこうする。
            await Api.postDepartmentsId(v.id!, {department: v})
        }
    }

    return {list, refresh, updateOrder}
}

// ユーザー追加・更新。
export async function useDepartment(departmentId?: number) {

    const department = ref<Department>({
        id: null,
        name: "",
        order: 0,
        created_at: undefined,
        updated_at: undefined
    })
    if (departmentId !== undefined) {
        const {data} = await Api.getDepartmentsId(departmentId)
        if (data.department != undefined) {
            department.value.id = data.department.id
            department.value.name = data.department.name
            department.value.order = data.department.order
            department.value.created_at = data.department.created_at
            department.value.updated_at = data.department.updated_at
        }
    }

    return {department}

}

export async function postDepartment(department: Department, order: number,emit: Emit) {
    department.order = order
    const req: PostDepartmentsRequest = {
        department: department
    }
    await Api.postDepartments(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postDepartmentById(department: Department, emit: Emit) {
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

export async function deleteDepartmentById(id: number, emit: Emit) {
    await Api.deleteDepartmentsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



