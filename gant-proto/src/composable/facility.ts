import {Api} from "@/api/axios";
import {Facility, PostCopyFacilitysRequest, PostFacilitiesRequest} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {changeSort} from "@/utils/sort";
import {Emit, FacilityStatus, FacilityType} from "@/const/common";


// ユーザー一覧。特別ref系は必要ない。
export async function useFacilityTable() {
    const list = ref<Facility[]>([])
    const refresh = async () => {
        const resp = await Api.getFacilities()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()

    const updateOrder = async (index: number, direction: number) => {
        changeSort(list.value, index, direction)

        for (const v of list.value) {
            v.order = list.value.indexOf(v)
            // API直呼び出しは少し気持ち悪いが効率を考慮してこうする。
            await Api.postFacilitiesId(v.id!, {facility: v})
        }
    }

    return {list, refresh, updateOrder}
}

// ユーザー追加・更新。
export async function useFacility(facilityId?: number) {

    const facility = ref<Facility>({
        id: null,
        name: "",
        term_from: "",
        term_to: "",
        shipment_due_date: "",
        order: 0,
        created_at: undefined,
        updated_at: undefined,
        status: FacilityStatus.Enabled,
        type: FacilityType.Ordered
    })
    if (facilityId !== undefined) {
        const {data} = await Api.getFacilitiesId(facilityId)
        if (data.facility != undefined) {
            facility.value.id = data.facility.id
            facility.value.name = data.facility.name
            facility.value.term_from = data.facility.term_from.substring(0, 10)
            facility.value.term_to = data.facility.term_to.substring(0, 10)
            facility.value.shipment_due_date = data.facility.shipment_due_date.substring(0, 10)
            facility.value.order = data.facility.order
            facility.value.created_at = data.facility.created_at
            facility.value.updated_at = data.facility.updated_at
            facility.value.status = data.facility.status
            facility.value.type = data.facility.type
        }
    }

    return {facility}

}

export function validate(facility: Facility) {
    let isValid = true
    if (!facility.name) {
        toast.warning("名称は必須です")
        isValid = false
    }
    if (!facility.term_from) {
        toast.warning("開始日は必須です")
        isValid = false
    }
    if (!facility.term_to) {
        toast.warning("終了日は必須です")
        isValid = false
    }
    if (!facility.shipment_due_date) {
        toast.warning("出荷期日は必須です")
        isValid = false
    }
    return isValid
}

export async function postFacility(facility: Facility, order: number, emit: Emit) {
    facility.term_from = facility.term_from + "T00:00:00.00000+09:00"
    facility.term_to = facility.term_to + "T00:00:00.00000+09:00"
    facility.shipment_due_date = facility.shipment_due_date + "T00:00:00.00000+09:00"
    facility.order = order
    const req: PostFacilitiesRequest = {
        facility: facility
    }
    await Api.postFacilities(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
        emit('update')
    })
}

export async function copyFacility(facility: Facility, order: number, originalFacilityId: number, emit: Emit) {
    facility.term_from = facility.term_from + "T00:00:00.00000+09:00"
    facility.term_to = facility.term_to + "T00:00:00.00000+09:00"
    facility.shipment_due_date = facility.shipment_due_date + "T00:00:00.00000+09:00"
    delete (facility.created_at)
    delete (facility.updated_at)
    facility.order = order
    const req: PostCopyFacilitysRequest = {
        facility: facility,
        facilityId: originalFacilityId
    }
    await Api.postCopyFacilitys(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
        emit('update')
    })
}

export async function postFacilityById(facility: Facility, emit: Emit) {
    facility.term_from = facility.term_from + "T00:00:00.00000+09:00"
    facility.term_to = facility.term_to + "T00:00:00.00000+09:00"
    facility.shipment_due_date = facility.shipment_due_date + "T00:00:00.00000+09:00"
    const req: PostFacilitiesRequest = {
        facility: facility
    }
    if (facility.id != null) {
        await Api.postFacilitiesId(facility.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
            emit('update')
        })
    }
}

export async function deleteFacilityById(id: number, emit: Emit) {
    await Api.deleteFacilitiesId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
        emit('update')
    })
}



