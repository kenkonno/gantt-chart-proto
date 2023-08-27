import {Api} from "@/api/axios";
import {PostUnitsRequest, Unit} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function useUnitTable() {
    const list = ref<Unit[]>([])
    const refresh = async () => {
        const resp = await Api.getUnits()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useUnit(unitId?: number) {

    const unit = ref<Unit>({
        id: null,
        name: "",
        created_at: undefined,
        updated_at: undefined
    })
    if (unitId !== undefined) {
        const {data} = await Api.getUnitsId(unitId)
        if (data.unit != undefined) {
            unit.value.id = data.unit.id
            unit.value.name = data.unit.name
            unit.value.created_at = data.unit.created_at
            unit.value.updated_at = data.unit.updated_at
        }
    }

    return {unit}

}

export async function postUnit(unit: Unit, emit: any) {
    const req: PostUnitsRequest = {
        unit: unit
    }
    await Api.postUnits(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postUnitById(unit: Unit, emit: any) {
    const req: PostUnitsRequest = {
        unit: unit
    }
    if (unit.id != null) {
        await Api.postUnitsId(unit.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteUnitById(id: number, emit: any) {
    await Api.deleteUnitsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



