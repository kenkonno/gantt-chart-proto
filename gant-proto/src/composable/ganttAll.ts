import dayjs from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {computed, inject, ref} from "vue";
import {Facility, Holiday, Ticket, User} from "@/api";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {
    endOfDay,
} from "@/coreFunctions/manHourCalculation";
import {DAYJS_FORMAT} from "@/utils/day";
import {Api} from "@/api/axios";
import {GanttBarConfig} from "@infectoone/vue-ganttastic/lib_types/types";
import {round} from "@/utils/math";
import {DEFAULT_PROCESS_COLOR, FacilityStatus} from "@/const/common";
import {DisplayType, Header} from "@/composable/ganttAllMenu";

const BAR_COMPLETE_COLOR = "rgb(76 255 18)"

type GanttAllRow = {
    facility: Facility,
    users: User[],
    estimate: number,
    startDate: string,
    endDate: string,
    progress_percent: number,
    bars: GanttBarObject[]
}

export async function useGanttAll() {
    // injectはsetupと同期的に呼び出す必要あり
    const {facilityList, userList, processList, facilityTypes} = inject(GLOBAL_STATE_KEY)!
    const {refreshFacilityList, refreshUserList, refreshProcessList} = inject(GLOBAL_ACTION_KEY)!
    await refreshFacilityList()
    await refreshUserList()
    await refreshProcessList()

    const holidays: Holiday[] = []

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
    const {data: allTicketUsers} = await Api.getTicketUsers(allTickets.list.map(v => v.id!))

    let filteredFacilityList = facilityList.filter(v => v.status === FacilityStatus.Enabled)
    if(facilityTypes.length > 0 ) {
        filteredFacilityList = filteredFacilityList.filter(v => facilityTypes.includes(v.type) )
    }

    // 全設備の最小
    const startDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_from > b.term_from ? 1 : -1
    }).shift()!.term_from.substring(0, 10)
    // 全設備の最大
    const endDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_to > b.term_to ? 1 : -1
    }).pop()!.term_to.substring(0, 10)

    // 設備ごとに行を作成する
    const ganttAllRowPromise = filteredFacilityList.map(async facility => {
        // 設備に紐づくチケット一覧
        const {data: ganttGroups} = await Api.getGanttGroups(facility.id!)
        const {data} = await Api.getHolidays(facility.id!)
        holidays.push(...data.list)
        const tickets = allTickets.list.filter(v => ganttGroups.list.map(v => v.id!).includes(v.gantt_group_id))
        const ticketUsers = allTicketUsers.list.filter(v => tickets.map(vv => vv.id!).includes(v.ticket_id))
        // 全てのチケットの予定工数を計上する
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

        const users: User[] = []
        {
            const r = ticketUsers.map(ticketUser => {
                return userList.find(v => v.id === ticketUser.user_id)!
            })
            users.push(...Array.from(new Set(r)))
        }

        // ここのbarsが複数なので１つにして日付を最小最大にする。
        return <GanttAllRow>{
            facility: facility,
            startDate: facility.term_from.substring(0, 10),
            endDate: facility.term_to.substring(0, 10),
            users: users,
            estimate: estimate,
            progress_percent: round(progress_percent),
            bars: [createBar(progress_percent, facility.name, facility.id!, facility.term_from, facility.term_to)],
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
                id: facilityId.toString(),
                immobile: false,
                label: facilityName, // 工程名
                progress: progressPercent,
                progressColor: BAR_COMPLETE_COLOR,
                pushOnOverlap: false,
                style: {backgroundColor: DEFAULT_PROCESS_COLOR},
            }
        }
    }

    // １本に各工程を混ぜるのは未使用になった。一旦残しておく。
    const createBars = (tickets: Ticket[]) => {
        const bars: GanttBarObject[] = []
        bars.push(
            ...tickets.filter(v => v.process_id ).map(ticket => {
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
    const holidaysAsDate: Date[] = []
    {
        const r = holidays.map(v => new Date(v.date))
        holidaysAsDate.push(...Array.from(new Set(r)))
    }

    return {
        startDate,
        endDate,
        ganttAllRow,
        holidaysAsDate,
        getGanttChartWidth,
        tickets: allTickets.list,
        ticketUsers: allTicketUsers.list,
        holidays,
        chartStart,
        chartEnd
    }
}


