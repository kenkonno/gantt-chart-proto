import dayjs from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {computed, inject, ref} from "vue";
import {Facility, Holiday, Ticket, User} from "@/api";
import {useGanttGroupTable} from "@/composable/ganttGroup";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {
    endOfDay,
} from "@/coreFunctions/manHourCalculation";
import {DAYJS_FORMAT} from "@/utils/day";
import {Api} from "@/api/axios";
import {GanttBarConfig} from "@infectoone/vue-ganttastic/lib_types/types";
import {round} from "@/utils/math";

type Header = {
    name: string,
    visible: boolean
}

export type DisplayType = "day" | "week" | "hour" | "month"

const BAR_NORMAL_COLOR = "rgb(147 206 255)"
const BAR_COMPLETE_COLOR = "rgb(76 255 18)"
const BAR_DANGER_COLOR = "rgb(255 89 89)"

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
    const {facilityList, userList, processList} = inject(GLOBAL_STATE_KEY)!
    const {refreshFacilityList, refreshUserList, refreshProcessList} = inject(GLOBAL_ACTION_KEY)!
    await refreshFacilityList()
    await refreshUserList()
    await refreshProcessList()

    const displayType = ref<DisplayType>("week")

    const GanttHeader = ref<Header[]>([
        {name: "設備名", visible: true},
        {name: "担当者", visible: false},
        {name: "開始日", visible: true},
        {name: "終了日", visible: true},
        {name: "工数(h)", visible: true},
        {name: "進捗", visible: true},
    ])
    const holidays: Holiday[] = []

    const getGanttChartWidth = computed<string>(() => {
        // 1日30pxとして計算する
        return (dayjs(endDate).diff(dayjs(startDate), displayType.value) + 1) * 30 + "px"
    })

    const {data: allTickets} = await Api.getAllTickets()
    const {data: allTicketUsers} = await Api.getTicketUsers(allTickets.list.map(v => v.id!))
    // 全設備の最小
    const startDate: string = facilityList.slice().sort((a, b) => {
        return a.term_from > b.term_from ? 1 : -1
    }).shift()!.term_from.substring(0, 10)
    // 全設備の最大
    const endDate: string = facilityList.slice().sort((a, b) => {
        return a.term_to > b.term_to ? 1 : -1
    }).pop()!.term_to.substring(0, 10)

    // 設備ごとに行を作成する
    const ganttAllRowPromise = facilityList.map(async facility => {
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

        return <GanttAllRow>{
            facility: facility,
            startDate: facility.term_from.substring(0, 10),
            endDate: facility.term_to.substring(0, 10),
            users: users,
            estimate: estimate,
            progress_percent: round(progress_percent),
            bars: createBars(tickets),
        }
    })
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
                        style: {backgroundColor: BAR_NORMAL_COLOR},
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
        GanttHeader,
        startDate,
        endDate,
        ganttAllRow,
        holidaysAsDate,
        displayType,
        getGanttChartWidth,
        tickets: allTickets.list,
        ticketUsers: allTicketUsers.list,
        holidays,
        chartStart,
        chartEnd
    }
}

