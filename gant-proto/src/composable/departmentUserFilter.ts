import {InjectionKey, Ref, ref} from "vue";

export const GLOBAL_DEPARTMENT_USER_FILTER_KEY = Symbol() as InjectionKey<GlobalDepartmentUserFilter>

export interface GlobalDepartmentUserFilter {
    selectedDepartment: Ref<number | undefined>,
    selectedUser: Ref<number | undefined>
}

export function useDepartmentUserFilter() {

    const selectedDepartment = ref<number | undefined>(undefined)
    const selectedUser = ref<number | undefined>(undefined)

    return {
        selectedDepartment,
        selectedUser
    }

}