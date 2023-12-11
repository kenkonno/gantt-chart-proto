import {Api} from "@/api/axios";
import {PostUnitsRequest, Unit} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {changeSort} from "@/utils/sort";
import {Emit} from "@/const/common";


// ユーザー一覧。特別ref系は必要ない。
export async function useUnitTable(facilityId: number) {
    const list = ref<Unit[]>([])
    const refresh = async () => {
        const resp = await Api.getUnits(facilityId)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    const updateOrder = async (index: number, direction: number) => {
        changeSort(list.value, index, direction)

        for (const v of list.value) {
            v.order = list.value.indexOf(v)
            // API直呼び出しは少し気持ち悪いが効率を考慮してこうする。
            await Api.postUnitsId(v.id!, {unit: v})
        }
    }

    return {list, refresh, updateOrder}
}

// ユーザー追加・更新。
export async function useUnit(unitId?: number) {

    const unit = ref<Unit>({
        id: null,
        name: "",
        facility_id: 0,
        order: 0,
        created_at: undefined,
        updated_at: undefined
    })
    if (unitId !== undefined) {
        const {data} = await Api.getUnitsId(unitId)
        if (data.unit != undefined) {
            unit.value.id = data.unit.id
            unit.value.name = data.unit.name
            unit.value.facility_id = data.unit.facility_id
            unit.value.order = data.unit.order
            unit.value.created_at = data.unit.created_at
            unit.value.updated_at = data.unit.updated_at
        }
    }

    return {unit}

}

export async function postUnit(unit: Unit, facilityId: number, order: number, emit: Emit) {
    unit.order = order
    unit.facility_id = facilityId
    const req: PostUnitsRequest = {
        unit: unit,
    }
    await Api.postUnits(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postUnitById(unit: Unit, facilityId: number, emit: Emit) {
    unit.facility_id = facilityId
    const req: PostUnitsRequest = {
        unit: unit,
    }
    if (unit.id != null) {
        await Api.postUnitsId(unit.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function deleteUnitById(id: number, emit: Emit) {
    await Api.deleteUnitsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



