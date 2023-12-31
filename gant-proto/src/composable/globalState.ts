/**
 * Global State
 * facilityに紐づかないものは直接
 * facilityに紐づくものは連想はれいつとして持つ。
 */
import {Department, Facility, Holiday, OperationSetting, Process, ScheduleAlert, Unit, User} from "@/api";
import {InjectionKey, ref, nextTick} from "vue";
import {Api} from "@/api/axios";
import {changeSort} from "@/utils/sort";
import {FacilityType} from "@/const/common";
import router from "@/router";

type GlobalState = {
    currentFacilityId: number,
    facilityList: Facility[]
    processList: Process[]
    departmentList: Department[]
    userList: User[]
    unitMap: { [key: number]: Unit[] }
    operationSettingMap: { [key: number]: OperationSetting[] }
    holidayMap: { [key: number]: Holiday[] }
    scheduleAlert: ScheduleAlert[]
    facilityTypes: string[]
    ganttFacilityRefresh: boolean
    ganttAllRefresh: boolean
    pileUpsRefresh: boolean
}

export const GLOBAL_STATE_KEY = Symbol() as InjectionKey<GlobalState>
export const GLOBAL_ACTION_KEY = Symbol() as InjectionKey<Actions>
export const GLOBAL_MUTATION_KEY = Symbol() as InjectionKey<Mutations>
export const GLOBAL_GETTER_KEY = Symbol() as InjectionKey<Getters>

interface Actions {
    refreshFacilityList: () => Promise<void>;
    refreshProcessList: () => Promise<void>;
    refreshDepartmentList: () => Promise<void>;
    refreshUserList: () => Promise<void>;
    refreshUnitMap: (facilityId: number) => Promise<void>;
    refreshHolidayMap: (facilityId: number) => Promise<void>;
    refreshOperationSettingMap: (facilityId: number) => Promise<void>;
    updateFacilityOrder: (index: number, direction: number) => Promise<void>;
    updateProcessOrder: (index: number, direction: number) => Promise<void>;
    updateDepartmentOrder: (index: number, direction: number) => Promise<void>;
    updateUnitOrder: (index: number, direction: number) => Promise<void>;

    getScheduleAlert(): Promise<void>;
}

interface Mutations {
    updateCurrentFacilityId: (id: number) => void;

    refreshGantt(id: number, moveFacilityView: boolean): void;

    setFacilityTypes(facilityType: string[]): void;

    refreshPileUpsRefresh(): void;

    refreshGanttAll(): void;

}

interface Getters {
    getDepartmentName: (departmentId: number) => string;
}


export const useGlobalState = async () => {
    const globalState = ref<GlobalState>({
        currentFacilityId: -1,
        departmentList: [],
        facilityList: [],
        processList: [],
        holidayMap: {},
        operationSettingMap: {},
        unitMap: {},
        userList: [],
        scheduleAlert: [],
        facilityTypes: [FacilityType.Ordered], // TODO: 初期値を記憶するようにする
        ganttFacilityRefresh: true, // refreshさせるだけのフラグ
        ganttAllRefresh: true,
        pileUpsRefresh: true,
    })
    const init = async () => {
        await actions.refreshFacilityList()
        await actions.refreshProcessList()
        await actions.refreshDepartmentList()
        await actions.refreshUserList()
    }

    const actions: Actions = {
        // refresh系
        refreshFacilityList: async () => {
            const resp = await Api.getFacilities()
            globalState.value.facilityList.splice(0, globalState.value.facilityList.length)
            globalState.value.facilityList.push(...resp.data.list)
        }, refreshProcessList: async () => {
            const resp = await Api.getProcesses()
            globalState.value.processList.splice(0, globalState.value.processList.length)
            globalState.value.processList.push(...resp.data.list)
        }, refreshDepartmentList: async () => {
            const resp = await Api.getDepartments()
            globalState.value.departmentList.splice(0, globalState.value.departmentList.length)
            globalState.value.departmentList.push(...resp.data.list)
        }, refreshUserList: async () => {
            const resp = await Api.getUsers()
            globalState.value.userList.splice(0, globalState.value.userList.length)
            globalState.value.userList.push(...resp.data.list)
        }, refreshUnitMap: async (facilityId: number) => {
            const resp = await Api.getUnits(facilityId)
            if (globalState.value.unitMap[facilityId] == null) globalState.value.unitMap[facilityId] = []
            globalState.value.unitMap[facilityId].splice(0, globalState.value.unitMap[facilityId].length)
            globalState.value.unitMap[facilityId].push(...resp.data.list)

        }, refreshHolidayMap: async (facilityId: number) => {
            const resp = await Api.getHolidays(facilityId)
            if (globalState.value.holidayMap[facilityId] == null) globalState.value.holidayMap[facilityId] = []
            globalState.value.holidayMap[facilityId].splice(0, globalState.value.holidayMap[facilityId].length)
            globalState.value.holidayMap[facilityId].push(...resp.data.list)

        }, refreshOperationSettingMap: async (facilityId: number) => {
            const resp = await Api.getOperationSettingsId(facilityId)
            if (globalState.value.operationSettingMap[facilityId] == null) globalState.value.operationSettingMap[facilityId] = []
            globalState.value.operationSettingMap[facilityId].splice(0, globalState.value.operationSettingMap[facilityId].length)
            globalState.value.operationSettingMap[facilityId].push(...resp.data.operationSettings)
        },
        // Order系
        updateFacilityOrder: async (index: number, direction: number) => {
            changeSort(globalState.value.facilityList, index, direction)
            for (const v of globalState.value.facilityList) {
                v.order = globalState.value.facilityList.indexOf(v)
                await Api.postFacilitiesId(v.id!, {facility: v})
            }
        }, updateProcessOrder: async (index: number, direction: number) => {
            changeSort(globalState.value.processList, index, direction)
            for (const v of globalState.value.processList) {
                v.order = globalState.value.processList.indexOf(v)
                await Api.postProcessesId(v.id!, {process: v})
            }
        }, updateDepartmentOrder: async (index: number, direction: number) => {
            changeSort(globalState.value.departmentList, index, direction)
            for (const v of globalState.value.departmentList) {
                v.order = globalState.value.departmentList.indexOf(v)
                await Api.postDepartmentsId(v.id!, {department: v})
            }
        }, updateUnitOrder: async (index: number, direction: number) => {
            changeSort(globalState.value.unitMap[globalState.value.currentFacilityId], index, direction)
            for (const v of globalState.value.unitMap[globalState.value.currentFacilityId]) {
                v.order = globalState.value.unitMap[globalState.value.currentFacilityId].indexOf(v)
                await Api.postUnitsId(v.id!, {unit: v})
            }
        }, getScheduleAlert: async () => {
            const resp = await Api.getScheduleAlerts()
            globalState.value.scheduleAlert.length = 0
            globalState.value.scheduleAlert.push(...resp.data.list)
        }
    }

    const mutations: Mutations = {
        updateCurrentFacilityId: (id: number) => {
            console.log("CHANGE FACILITY ID")
            globalState.value.currentFacilityId = id
        },
        refreshGantt: async (facilityId: number, moveFacilityView = false) => {
            globalState.value.currentFacilityId = facilityId
            globalState.value.ganttFacilityRefresh = false
            // facility紐づくデータを初期化する
            await actions.refreshHolidayMap(facilityId)
            await actions.refreshUnitMap(facilityId)
            await actions.refreshOperationSettingMap(facilityId)
            nextTick(() => {
                if (moveFacilityView) {
                    router.push("/")
                }
                globalState.value.ganttFacilityRefresh = true
            })
        },
        refreshPileUpsRefresh: async () => {
            globalState.value.pileUpsRefresh = false
            nextTick(() => {
                globalState.value.pileUpsRefresh = true
            })
        },
        refreshGanttAll: () => {
            globalState.value.ganttAllRefresh = false
            nextTick(() => {
                globalState.value.ganttAllRefresh = true
            })
        },
        setFacilityTypes: (facilityTypes: string[]) => {
            globalState.value.facilityTypes.length = 0
            globalState.value.facilityTypes = facilityTypes
        }
    }

    const getters: Getters = {
        getDepartmentName: (departmentId: number) => {
            return globalState.value.departmentList.find(v => v.id === departmentId)!.name
        }
    }

    await init()
    return {globalState, actions, mutations, getters}
}