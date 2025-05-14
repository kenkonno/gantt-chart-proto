import dayjs from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {computed, inject, ref} from "vue";
import {Facility, GanttGroup, Ticket, User} from "@/api";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {endOfDay,} from "@/coreFunctions/manHourCalculation";
import {DAYJS_FORMAT} from "@/utils/day";
import {Api} from "@/api/axios";
import {GanttBarConfig} from "@infectoone/vue-ganttastic/lib_types/types";
import {round} from "@/utils/math";
import {ApiMode, DEFAULT_PROCESS_COLOR, FacilityStatus} from "@/const/common";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import {useMilestoneTable} from "@/composable/milestone";
import {AggregationAxis, DisplayType} from "@/composable/ganttFacilityMenu";
import {getUserInfo} from "@/composable/auth";

const BAR_COMPLETE_COLOR = "rgb(200 200 200)"

type GanttAllRow = {
    facility: Facility,
    users: User[],
    estimate: number,
    startDate: string,
    endDate: string,
    progress_percent: number,
    bars: GanttBarObject[]
}

export async function useGanttAll(aggregationAxis: AggregationAxis) {
    // injectはsetupと同期的に呼び出す必要あり
    const {facilityList, userList, processList, facilityTypes, holidayList} = inject(GLOBAL_STATE_KEY)!
    const {refreshFacilityList, refreshUserList, refreshProcessList, refreshHolidayList} = inject(GLOBAL_ACTION_KEY)!
    const {selectedDepartment, selectedUser} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!
    await refreshHolidayList()
    await refreshFacilityList()
    await refreshUserList()
    await refreshProcessList()

    const getGanttChartWidth = (displayType: DisplayType) => {
        // 1日30pxとして計算する
        return (dayjs(endDate).diff(dayjs(startDate), displayType) + 1) * 30 + "px"
    }

    const getProcessColor = (id?: number | null) => {
        if (id == null) {
            return DEFAULT_PROCESS_COLOR
        }
        return processList.find(v => v.id === id)?.color
    }

    const {data: allTickets} = await Api.getAllTickets()
    const {data: originalAllTickets} = await Api.getAllTickets(undefined, ApiMode.prod)
    const {data: allTicketUsers} = await Api.getTicketUsers([])

    let filteredFacilityList = facilityList.filter(v => v.status === FacilityStatus.Enabled)
    if (facilityTypes.length > 0) {
        filteredFacilityList = filteredFacilityList.filter(v => facilityTypes.includes(v.type))
    }

    const filteredAllTickets = allTickets
    const filteredAllTicketUsers = allTicketUsers
    // チケット、案件の絞り込みを実施する
    if (selectedDepartment.value != undefined) {
        filteredAllTickets.list = filteredAllTickets.list.filter(v => v.department_id == selectedDepartment.value)
        const ticketIds = filteredAllTickets.list.map(v => v.id)
        filteredAllTicketUsers.list = filteredAllTicketUsers.list.filter(v => ticketIds.includes(v.ticket_id))
    }
    if (selectedUser.value != undefined) {
        filteredAllTicketUsers.list = filteredAllTicketUsers.list.filter(v => v.user_id == selectedUser.value)
        const ticketIds = filteredAllTicketUsers.list.map(v => v.ticket_id)
        filteredAllTickets.list = filteredAllTickets.list.filter(v => ticketIds.includes(v.id!))
    }
    const hasFilter = () => {
        return selectedDepartment.value != undefined || selectedUser.value != undefined
    }
    // TODO: ここでチケットから 案件が絞り込めれば長楽 コードに違和感があるがAPI呼び出しをする
    if (hasFilter()) {
        const ganttGroupIds = Array.from(new Set(filteredAllTickets.list.map(v => v.gantt_group_id)))
        const facilityIds = await Promise.all(ganttGroupIds.map(async ganttGroupId => {
            const {data} = await Api.getGanttGroupsId(ganttGroupId)
            return data.ganttGroup?.facility_id
        }))
        const facilityIdFlat = facilityIds.flat()
        filteredFacilityList = filteredFacilityList.filter(v => facilityIdFlat.includes(v.id!))
    }

    // 全案件の最小
    const startDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_from > b.term_from ? 1 : -1
    }).shift()?.term_from.substring(0, 10) ?? dayjs(Date()).startOf('month').format("YYYYMMDD")
    // 全案件の最大
    const endDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_to > b.term_to ? 1 : -1
    }).pop()?.term_to.substring(0, 10) ?? dayjs(Date()).endOf('month').format("YYYYMMDD")

    async function getMilestones(facility: Facility, mode?: string) {
        return await (async () => {
            const result: GanttBarObject[] = []
            const {list} = await useMilestoneTable(facility.id!, mode)
            result.push(<GanttBarObject>{
                beginDate: dayjs(facility.shipment_due_date).format(DAYJS_FORMAT),
                endDate: endOfDay(facility.shipment_due_date),
                ganttBarConfig: {
                    bundle: "",
                    dragLimitLeft: 0,
                    dragLimitRight: 0,
                    hasHandles: false,
                    immobile: false,
                    pushOnOverlap: false,
                    id: 90000 + facility.id! + "", // TODO: IDのルールの決定
                    label: "出荷期日".substring(0, 1), // 工程名
                    progress: 0,
                    progressColor: BAR_COMPLETE_COLOR,
                    style: {
                        backgroundColor: "yellow",
                        padding: 0,
                        height: "90%",
                        opacity: 0.9,
                        borderRadius: "30px"
                    },
                }
            })
            result.push(
                ...list.value.map(v => {
                    return <GanttBarObject>{
                        beginDate: dayjs(v.date).format(DAYJS_FORMAT),
                        endDate: endOfDay(v.date),
                        ganttBarConfig: {
                            bundle: "",
                            dragLimitLeft: 0,
                            dragLimitRight: 0,
                            hasHandles: false,
                            immobile: false,
                            pushOnOverlap: false,
                            id: 80000 + facility.id! + "", // TODO: IDのルールの決定
                            label: v.description.substring(0, 1), // 工程名
                            progress: 0,
                            progressColor: BAR_COMPLETE_COLOR,
                            style: {
                                backgroundColor: "blue",
                                padding: 0,
                                height: "90%",
                                opacity: 0.9,
                                borderRadius: "30px",
                                color: "white"
                            },
                        }
                    }
                })
            )
            return result
        })();
    }

    const getFacilityInfos = async (ganttGroups: GanttGroup[], sourceTickets: Ticket[]) => {
        const tickets = sourceTickets.filter(v => ganttGroups.map(v => v.id!).includes(v.gantt_group_id))
        const estimate = tickets.reduce((p, c) => {
            if (c.estimate != null) {
                return p + c.estimate
            } else {
                return p
            }
        }, 0)
        // 進捗は 消化工数 / 全体工数 * 100
        const progress_percent = tickets.reduce((p, c) => {
            if (c.estimate != null && c.progress_percent != null) {
                return p + c.estimate * (c.progress_percent / 100)
            } else {
                return p
            }
        }, 0) / estimate * 100


        return {
            tickets,
            estimate,
            progress_percent
        }
    }

    const getProdBarsIfSimulation = async (facility: Facility) => {
        const {data} = await Api.getFacilitiesId(facility.id!, ApiMode.prod)
        const prodFacility = data.facility!

        const isSimulateUser = getUserInfo().isSimulateUser
        const bars: GanttBarObject[] = []
        if (isSimulateUser == true) {
            const {data: ganttGroups} = await Api.getGanttGroups(prodFacility.id!, ApiMode.prod)
            const milestones = await getMilestones(prodFacility, ApiMode.prod);
            const {
                tickets,
                progress_percent
            } = await getFacilityInfos(ganttGroups.list, originalAllTickets.list)
            bars.push(...(hasFilter() || aggregationAxis == 'process' ? createBars(tickets) : [
                createBar(progress_percent, prodFacility.name, prodFacility.id!, prodFacility.term_from, prodFacility.term_to),
                ...milestones
            ]))

        }
        return bars
    }

    // 案件ごとに行を作成する TODO: ここで呼び出しているAPIは一括で取得できるようにする。パフォーマンス改善。全件取得にオプションのパラメータでID指定できるようにする。
    const ganttAllRowPromise = filteredFacilityList.map(async facility => {
        // 案件に紐づくチケット一覧
        const {data: ganttGroups} = await Api.getGanttGroups(facility.id!)
        const {tickets, estimate, progress_percent} = await getFacilityInfos(ganttGroups.list, filteredAllTickets.list)

        const ticketUsers = filteredAllTicketUsers.list.filter(v => tickets.map(vv => vv.id!).includes(v.ticket_id))
        // 全てのチケットの予定工数を計上する
        const users: User[] = []
        {
            const r = ticketUsers.map(ticketUser => {
                return userList.find(v => v.id === ticketUser.user_id)!
            })
            users.push(...Array.from(new Set(r)))
        }

        const milestones = await getMilestones(facility);

        const bars = hasFilter() || aggregationAxis == 'process' ? createBars(tickets) : [
            createBar(progress_percent, facility.name, facility.id!, facility.term_from, facility.term_to),
            ...milestones
        ]
        const prodBars = await getProdBarsIfSimulation(facility)
        console.log("prodBars vs bars", prodBars, bars)
        if (prodBars.length > 0) {
            // 表示用のbarsにprodBarsに存在している者だけ追加する。
            const targets = prodBars.filter(v => {
                return bars.find(vv => vv.ganttBarConfig.id == v.ganttBarConfig.id) != undefined
            })
            bars.unshift(...targets.map(v => {
                const style = v.ganttBarConfig.style!
                style.opacity = 0.5
                v.ganttBarConfig.style = style
                return v
            }))
        }

        // ここのbarsが複数なので１つにして日付を最小最大にする。
        // TODO: CreateBarsでCSSをさらに変更させる
        return <GanttAllRow>{
            facility: facility,
            startDate: facility.term_from.substring(0, 10),
            endDate: facility.term_to.substring(0, 10),
            users: users,
            estimate: estimate,
            progress_percent: round(progress_percent),
            bars: bars,
        }
    })

    const createBar = (progressPercent: number, facilityName: string, facilityId: number, startDate: string, endDate: string) => {
        return <GanttBarObject>{
            beginDate: dayjs(startDate).format(DAYJS_FORMAT),
            endDate: endOfDay(endDate),
            ganttBarConfig: <GanttBarConfig>{
                bundle: "",
                dragLimitLeft: 0,
                dragLimitRight: 0,
                hasHandles: false,
                immobile: false,
                pushOnOverlap: false,
                id: facilityId.toString(),
                label: facilityName, // 工程名
                progress: progressPercent,
                progressColor: BAR_COMPLETE_COLOR,
                style: {backgroundColor: DEFAULT_PROCESS_COLOR},
            }
        }
    }

    // フィルタ済みの場合はこちらを利用する。
    const createBars = (tickets: Ticket[]) => {
        const bars: GanttBarObject[] = []
        bars.push(
            ...tickets.filter(v => v.process_id).map(ticket => {
                return <GanttBarObject>{
                    beginDate: dayjs(ticket.start_date!).format(DAYJS_FORMAT),
                    endDate: endOfDay(ticket!.end_date!),
                    ganttBarConfig: <GanttBarConfig>{
                        bundle: "",
                        dragLimitLeft: 0,
                        dragLimitRight: 0,
                        hasHandles: false,
                        id: ticket.id!.toString(),
                        immobile: false,
                        label: processList.find(v => v.id == ticket.process_id)!.name, // 工程名
                        progress: ticket.progress_percent,
                        progressColor: BAR_COMPLETE_COLOR,
                        pushOnOverlap: false,
                        style: {backgroundColor: getProcessColor(ticket.process_id)},
                    }
                }
            })
        )
        return bars
    }
    const ganttAllRow = await Promise.all(ganttAllRowPromise)

    const chartStart = ref(dayjs(startDate).format(DAYJS_FORMAT))
    const chartEnd = ref(dayjs(endDate).format(DAYJS_FORMAT))
    const holidaysAsDate = computed(() => (displayType: DisplayType) => {
        if (displayType === "week") {
            return Array.from([])
        } else {
            const r = holidayList.map(v => new Date(v.date))
            return Array.from(new Set(r))
        }
    })

    return {
        startDate,
        endDate,
        ganttAllRow,
        holidaysAsDate,
        getGanttChartWidth,
        tickets: filteredAllTickets.list,
        ticketUsers: filteredAllTicketUsers.list,
        holidayList,
        chartStart,
        chartEnd,
        hasFilter
    }
}
