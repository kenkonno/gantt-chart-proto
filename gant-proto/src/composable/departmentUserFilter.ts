import {InjectionKey, Ref, ref} from "vue";

export const GLOBAL_DEPARTMENT_USER_FILTER_KEY = Symbol() as InjectionKey<GlobalDepartmentUserFilter>

export interface GlobalDepartmentUserFilter {
    selectedDepartment: Ref<number[]>,
    selectedUser: Ref<number[]>
}

export function useDepartmentUserFilter() {

    const selectedDepartment = ref<number[]>([])
    const selectedUser = ref<number[]>([])

    return {
        selectedDepartment,
        selectedUser
    }

}