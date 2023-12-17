import {ScheduleAlert} from "@/api";
import {computed, ComputedRef, InjectionKey, Ref, ref} from "vue";

export const GLOBAL_SCHEDULE_ALERT_KEY = Symbol() as InjectionKey<GlobalScheduleAlert>

export type GlobalScheduleAlert = {
    tableIsOpen: Ref<boolean>,
    filterDelayDays: Ref<number | undefined>,
    filterProgressNumber: Ref<number | undefined>,
    filterFacility: Ref<number | undefined>,
    cScheduleAlert: ComputedRef<Map<number, ScheduleAlert[]>>
    facilityList: ComputedRef<{ code: number, name: string }[]>
}

// ユーザー一覧。特別ref系は必要ない。
export function useScheduleAlert(scheduleAlert: ScheduleAlert[]) {
    const tableIsOpen = ref<boolean>(false)

    const filterDelayDays = ref<number | undefined>(undefined)
    const filterProgressNumber = ref<number | undefined>(undefined)
    const filterFacility = ref<number | undefined>(undefined)

    const cScheduleAlert = computed(() => {
        const result: Map<number, ScheduleAlert[]> = new Map()
        let filtered: ScheduleAlert[]
        if (filterDelayDays.value != undefined) {
            filtered = scheduleAlert.filter(v => v.delay_days >= filterDelayDays.value!)
        } else {
            filtered = scheduleAlert
        }
        if (filterProgressNumber.value != undefined && filterProgressNumber.value > 0) {
            filtered = filtered.filter(v => v.progress_percent <= filterProgressNumber.value!)
        }
        if (filterFacility.value != undefined && filterFacility.value > -1) {
            filtered = filtered.filter(v => v.facility_id == filterFacility.value!)
        }

        filtered.forEach((v => {
            if (result.get(v.facility_id) == undefined) {
                result.set(v.facility_id, [])
            }
            result.get(v.facility_id)!.push(v)
        }))
        return result
    })

    const facilityList = computed(() => {
        const result: { code: number, name: string }[] = [
            {code: -1, name: ""}
        ]
        scheduleAlert.forEach(v => {
            if (!result.find(vv => vv.code == v.facility_id)) {
                result.push({code: v.facility_id, name: v.facility_name})
            }
        })
        return result
    })
    return {
        tableIsOpen, filterDelayDays, filterProgressNumber, filterFacility, cScheduleAlert, facilityList
    }
}
