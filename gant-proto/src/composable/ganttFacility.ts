import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {computed, inject, ref, toValue} from "vue";
import {GanttGroup, Holiday, OperationSetting, Ticket, TicketUser} from "@/api";
import {Api} from "@/api/axios";
import {useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";
import {useTicketUserTable} from "@/composable/ticketUser";
import {useFacility} from "@/composable/facility";
import {changeSort} from "@/utils/sort";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {
    adjustStartDateByHolidays,
    dayBetween,
    endOfDay,
    ganttDateToYMDDate,
    getEndDateByRequiredBusinessDay,
    getNumberOfBusinessDays
} from "@/coreFunctions/manHourCalculation";
import {DAYJS_FORMAT} from "@/utils/day";

type GanttChartGroup = {
    ganttGroup: GanttGroup // TODO: 結局ganttRowにganttGroup設定しているから設計としては微妙っぽい。
    rows: GanttRow[]
}

export type GanttRow = {
    bar: GanttBarObject
    ganttGroup?: GanttGroup
    ticket?: Ticket
    ticketUsers?: TicketUser[]
}

type Header = {
    name: string,
    visible: boolean
}

export type DisplayType = "day" | "week" | "hour" | "month"

const BAR_NORMAL_COLOR = "rgb(147 206 255)"
const BAR_COMPLETE_COLOR = "rgb(76 255 18)"
const BAR_DANGER_COLOR = "rgb(255 89 89)"


function getScheduledOperatingHours(operationSettings: OperationSetting[], row: GanttRow) {
    return operationSettings.filter(operationSetting => {
        return operationSetting.facility_id === row.ganttGroup?.facility_id &&
            operationSetting.unit_id === row.ganttGroup?.unit_id
    }).reduce((accumulateValue, currentValue) => {
        const workHours = currentValue.workHours.find(v => v.process_id === row.ticket?.process_id)
        return accumulateValue + workHours!.work_hour!
    }, 0);
}

export async function useGanttFacility() {
    // injectはsetupと同期的に呼び出す必要あり
    const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
    const {userList, processList, departmentList, holidayMap, operationSettingMap, unitMap} = inject(GLOBAL_STATE_KEY)!

    const GanttHeader = ref<Header[]>([
        {name: "ユニット", visible: true},
        {name: "工程", visible: true},
        {name: "部署", visible: true},
        {name: "担当者", visible: true},
        {name: "人数", visible: false},
        {name: "期日", visible: false},
        {name: "工数(h)", visible: true},
        {name: "日後", visible: false},
        {name: "開始日", visible: false},
        {name: "終了日", visible: false},
        {name: "進捗", visible: true},
        {name: "操作", visible: false},
    ])
    const displayType = ref<DisplayType>("day")
    const {facility} = await useFacility(currentFacilityId)
    const chartStart = ref(dayjs(facility.value.term_from).format(DAYJS_FORMAT))
    const chartEnd = ref(dayjs(facility.value.term_to).format(DAYJS_FORMAT))

    const getUnitName = computed(() => (id: number) => {
        return unitMap[currentFacilityId].find(v => v.id === id)?.name
    })
    const getDepartmentName = computed(() => (id: number) => {
        return departmentList.find(v => v.id === id)?.name
    })
    const getProcessName = (id: number) => {
        return processList.find(v => v.id === id)?.name
    }
    const getOperationList = computed(() => {
        return operationSettingMap[currentFacilityId]
    })
    const getHolidaysForGantt = computed(() => {
        if (displayType.value === "day") {
            return holidayMap[currentFacilityId].map(v => new Date(v.date))
        } else {
            return []
        }
    })
    const getHolidays = computed(() => {
        return holidayMap[currentFacilityId]
    })
    const getGanttChartWidth = computed<string>(() => {
        // 1日30pxとして計算する
        return (dayjs(facility.value.term_to).diff(dayjs(facility.value.term_from), displayType.value) + 1) * 30 + "px"
    })


    // とりあえず何も考えずにAPIからガントチャート表示に必要なオブジェクトを作る
    const {list: ganttGroupList, refresh: ganttGroupRefresh} = await useGanttGroupTable()
    const {list: ticketList, refresh: ticketRefresh} = await useTicketTable()
    const {list: ticketUserList, refresh: ticketUserRefresh} = await useTicketUserTable()
    await ganttGroupRefresh(currentFacilityId)
    const ganttGroupIds = ganttGroupList.value.map(v => v.id!)
    await ticketRefresh(ganttGroupIds)
    const ticketIds = ticketList.value.map(v => v.id!)
    await ticketUserRefresh(ticketIds)
    // ここからガントチャートに渡すオブジェクトを作成する
    const ganttChartGroup = ref<GanttChartGroup[]>([])
    const refreshLocalGantt = () => {
        ganttChartGroup.value.length = 0
        ganttGroupList.value.forEach(ganttGroup => {
            ganttChartGroup.value.push(
                {
                    ganttGroup: ganttGroup,
                    rows: ticketList.value.filter(ticket => ticket.gantt_group_id === ganttGroup.id)
                        .sort(v => v.order)
                        .map(ticket => ticketToGanttRow(ticket, ticketUserList.value, ganttGroup))
                }
            )
        })
    }
    // computedだとドラッグ関連が上手くいかない為やはりオブジェクトとして双方向に同期をとるようにする。
    const bars = ref<GanttBarObject[]>([])
    // ガントチャート描画用のオブジェクトの生成
    const refreshBars = () => {
        bars.value.length = 0;
        ganttChartGroup.value.forEach(unit => {
            const emptyRow: GanttBarObject = {
                beginDate: "",
                endDate: "",
                ganttBarConfig: {
                    id: Math.random().toString(),
                    style: {backgroundColor: BAR_NORMAL_COLOR}
                }
            }
            bars.value.push(...unit.rows.map(task => <GanttBarObject>{
                beginDate: dayjs(task.ticket!.start_date!).format(DAYJS_FORMAT),
                endDate: endOfDay(task.ticket!.end_date!),
                ganttGroupId: unit.ganttGroup.id,
                ganttBarConfig: {
                    hasHandles: true,
                    id: task.ticket?.id!.toString(),
                    label: getProcessName(task.ticket?.process_id == null ? -1 : task.ticket?.process_id),
                    style: {backgroundColor: BAR_NORMAL_COLOR},
                    progress: task.ticket?.progress_percent,
                    progressColor: BAR_COMPLETE_COLOR
                }
            }))
            // ボタン用の空行を追加する
            bars.value.push(emptyRow)
        })
    }
    refreshLocalGantt()
    refreshBars()

    const updateOrder = async (ganttRows: GanttRow[], index: number, direction: number) => {
        changeSort(ganttRows, index, direction)
        // ガントチャートのオブジェクトと同期をとる
        changeSort(bars.value, index, direction)
        for (const v of ganttRows) {
            v.ticket!.order = ganttRows.indexOf(v)
            await updateTicket(v.ticket!)
        }
    }

    // 部署の変更。
    const updateDepartment = async (ticket: Ticket) => {
        // 担当者をすべて外す。
        const ticketUsers = ticketUserList.value.filter(v => v.ticket_id === ticket.id)
        if (ticketUsers != null) {
            ticketUsers.length = 0
        }
        await ticketUserUpdate(ticket, [])
    }
    const getUserListByDepartmentId = (departmentId?: number) => {
        if (departmentId == null) {
            return []
        }
        return userList.filter(v => v.department_id === departmentId)
    }

    const deleteTicket = async (ticket: Ticket) => {
        // TODO: confirmを作る
        await Api.deleteTicketsId(ticket.id!)
        // チケット情報をリフレッシュする
        await ticketRefresh(ganttGroupIds)
        refreshLocalGantt()
        refreshBars()
    }

    // DBへのストア及びローカルのガントに情報を反映する
    const updateTicket = async (ticket: Ticket) => {
        const reqTicket = Object.assign({}, ticket)
        if (reqTicket.start_date != undefined) {
            reqTicket.start_date = ticket.start_date + "T00:00:00.00000+09:00"
        }
        if (reqTicket.end_date != undefined) {
            reqTicket.end_date = ticket.end_date + "T00:00:00.00000+09:00"
        }
        if (reqTicket.limit_date != undefined) {
            reqTicket.limit_date = ticket.limit_date + "T00:00:00.00000+09:00"
        }
        await Api.postTicketsId(ticket.id!, {ticket: reqTicket})
        reflectTicketToGantt(reqTicket)
    }

    // DBへのストア及びローカルのガントに情報を反映する
    const addNewTicket = async (ganttGroupId: number) => {
        const order = ticketList.value.filter(v => v.gantt_group_id === ganttGroupId).length + 1
        const newTicket: Ticket = {
            created_at: undefined,
            days_after: undefined,
            department_id: undefined,
            end_date: undefined,
            estimate: undefined,
            gantt_group_id: ganttGroupId,
            id: undefined,
            limit_date: undefined,
            order: order,
            process_id: undefined,
            progress_percent: undefined,
            start_date: undefined,
            updated_at: 0
        }
        const {data} = await Api.postTickets({ticket: newTicket})
        ticketList.value.push(data.ticket!)
        refreshLocalGantt()
        refreshBars()
    }
    // DBへのストア及びローカルのガントに情報を反映する
    const ticketUserUpdate = async (ticket: Ticket, userIds: number[]) => {
        const {data} = await Api.postTicketUsers({ticketId: ticket.id!, userIds: userIds})
        // ticketUserList から TicketId をつけなおす
        const newTicketUserList = ticketUserList.value.filter(v => v.ticket_id !== ticket.id!)
        newTicketUserList.push(...data.ticketUsers)
        ticketUserList.value.length = 0
        ticketUserList.value.push(...newTicketUserList)

        // 人数を更新する
        ticket.number_of_worker = data.ticketUsers.length
        await updateTicket(ticket)
        // TODO: refreshLocalGanttは重くなるので性能対策が必要かも
        refreshLocalGantt()
    }

    const adjustBar = async (bar: GanttBarObject) => {
        // id をもとに VueJs 側のデータも更新する
        const targetGanttGroup = ganttChartGroup.value.find(v => v.ganttGroup.id === bar.ganttGroupId)
        if (!targetGanttGroup) {
            console.error(`Unit ID: ${bar.ganttGroupId} が現在のガントに存在しません。`, ganttChartGroup)
        } else {
            const targetTicket = targetGanttGroup.rows.find(v => v.ticket?.id!.toString() === bar.ganttBarConfig.id)
            if (!targetTicket) {
                console.error(`ID: ${bar.ganttBarConfig.id} が現在のガントに存在しません。`, ganttChartGroup)
            } else {
                targetTicket.ticket!.start_date = ganttDateToYMDDate(bar.beginDate)
                targetTicket.ticket!.end_date = ganttDateToYMDDate(bar.endDate)
                await updateTicket(targetTicket.ticket!)
            }
        }
    }

    // モデル情報をガントチャート（bars）に反映させる
    const reflectTicketToGantt = (ticket: Ticket) => {
        const targetTicket = bars.value.find(v => v.ganttBarConfig.id! === ticket.id!.toString())
        if (!targetTicket) {
            console.error(`TicketID: ${ticket.id} が現在のガントに存在しません。`, bars)
        } else {
            // パフォーマンスのためにガントチャートに反映すべきものは特別にここで記述する
            targetTicket.beginDate = dayjs(ticket.start_date!).format(DAYJS_FORMAT)
            targetTicket.endDate = endOfDay(ticket.end_date!)
            if (ticket.progress_percent) {
                targetTicket.ganttBarConfig.progress = ticket.progress_percent
            }
            targetTicket.ganttBarConfig.label = getProcessName(ticket.process_id == null ? -1 : ticket.process_id)
        }
    }

    /**********************************************************
     *                   リスケジュール関連
     **********************************************************/
    /**
     * setScheduleByPersonDay
     * リスケ（工数h）重視
     * 稼働予定時間ベースでスケジュールを引き直す。
     *
     * @param rows
     * 色々考えたけど難しい（答えがない）のでマジでシンプルに１行ずつ処理するようにする。
     *
     * 基本的な仕様
     * ・人日重視で期間を変更する。
     * ・開始日・工数が設定されていない行は無視する
     * ・稼働の消化は担当者の設定順とする。
     * とりあえず作らないと話にならないので一旦作る
     */
    const setScheduleByPersonDay = (rows: GanttRow[]) => {
        const holidays = toValue(getHolidays)
        const operationSettings = toValue(getOperationList)
        let prevEndDate: Dayjs
        rows.forEach(row => {
            let startDate: Dayjs
            // 開始日・工数・工程・担当者が未設定の場合はスキップする
            if (!row.ticket?.start_date || !row.ticket?.estimate || !row.ticketUsers?.length || !row.ticket?.process_id) {
                if (row.ticket?.end_date != null) {
                    prevEndDate = adjustStartDateByHolidays(dayjs(row.ticket.end_date).add(1, 'day'), holidays)
                }
                return;
            }
            if (row.ticket.days_after != null && row.ticket.days_after! >= 0 && prevEndDate != null) {
                startDate = prevEndDate.add(row.ticket.days_after!, 'day') // 現在日を最後に１日足しているため。
            } else {
                startDate = dayjs(row.ticket!.start_date)
            }
            // 開始日が祝日だった場合ずらす
            startDate = adjustStartDateByHolidays(startDate, holidays)

            // 工程から総稼働予定時間を取得する
            const scheduledOperatingHours = getScheduledOperatingHours(operationSettings, row) * row.ticketUsers.length
            // 必要日数。少数が出たら最小の整数にする。例：二人で5人日の場合 3日必要 5/2 = 2.5 つまり少数が出たら整数に足す
            const numberOfDaysRequired = Math.ceil(row.ticket.estimate! / scheduledOperatingHours)
            // 開始日を設定する
            row.ticket.start_date = ganttDateToYMDDate(startDate.format(DAYJS_FORMAT))
            // 終了日を決定する、祝日が含まれている場合終了日をずらす。
            row.ticket.end_date = ganttDateToYMDDate(
                getEndDateByRequiredBusinessDay(startDate, numberOfDaysRequired, holidays)
                    .add(-1, 'minute').format(DAYJS_FORMAT)
            )
            // ガントに反映する
            reflectTicketToGantt(row.ticket)
            updateTicket(row.ticket)
            // 前回の終了日を設定する
            prevEndDate = adjustStartDateByHolidays(dayjs(row.ticket.end_date).add(1, 'day'), holidays)
        })
    }

    /**
     * リスケ（日付）重視
     * 開始日・終了日重視で人数を設定する
     * 期間で満了できるように人数を割り当てる
     * ※担当者が入っていたとしても人数は必ず入っているので、人数列で参照する。
     *
     * @param rows
     */
    const setScheduleByFromTo = (rows: GanttRow[]) => {
        // 必要人数の算出。スケジュールを満了できる最少の人数を計算する。
        // 担当者が未定の時に人数を計算する
        // 必要日数 = 工数 / 労働時間
        // 必要人数 = 必要日数 / 期間日数 の数値は上に寄せる
        const holidays = toValue(getHolidays)
        const operationSettings = toValue(getOperationList)
        let prevEndDate: Dayjs
        rows.forEach(row => {
            if (row.ticket && row.ticketUsers?.length == 0 &&
                row.ticket?.estimate && row.ticket?.estimate > 0 &&
                row.ticket?.process_id && row.ticket?.process_id > 0 &&
                row.ticket?.start_date && row.ticket?.end_date) {
                // 日後の処理
                if (row.ticket.days_after != null && row.ticket.days_after >= 0 && prevEndDate != null) {
                    // 現在設定されている営業日を計算する
                    let dayjsStartDate = dayjs(row.ticket.start_date)
                    let dayjsEndDate = dayjs(row.ticket.end_date)
                    const numberOfRequiredBusinessDays = getNumberOfBusinessDays(dayjsStartDate, dayjsEndDate, holidays)
                    // 開始日を日後分ずらした日付に設定する
                    dayjsStartDate = prevEndDate.add(row.ticket.days_after, 'days')
                    // 終了日を設定されている営業日分確保する
                    dayjsEndDate = getEndDateByRequiredBusinessDay(dayjsStartDate, numberOfRequiredBusinessDays, holidays)
                    console.log("##############", numberOfRequiredBusinessDays)
                    // チケットに反映させる
                    row.ticket.start_date = ganttDateToYMDDate(dayjsStartDate.format(DAYJS_FORMAT))
                    row.ticket.end_date = ganttDateToYMDDate(dayjsEndDate.format(DAYJS_FORMAT))
                }

                // 工数(h)
                const estimate = row.ticket.estimate
                // 労働予定時間 ユニット & 工程 に紐づく稼働予定時間を取得
                console.log(operationSettings)
                const scheduledOperatingHours = getScheduledOperatingHours(operationSettings, row)
                // 必要日数
                const requiredNumberOfDays = estimate / scheduledOperatingHours
                // 日数（期間）
                const numberOfDays = getNumberOfDays(row, holidays)
                // 必要人数
                row.ticket.number_of_worker = Math.ceil(requiredNumberOfDays / numberOfDays)
                // 前回の終了日を設定する
                updateTicket(row.ticket)
            }
            prevEndDate = adjustStartDateByHolidays(dayjs(row.ticket?.end_date).add(1, 'day'), holidays)
        })
    }
    return {
        GanttHeader,
        bars,
        chartEnd,
        chartStart,
        displayType,
        facility,
        getHolidaysForGantt,
        ganttChartGroup,
        getDepartmentName,
        getGanttChartWidth,
        getUnitName,
        getOperationList,
        getHolidays,
        addNewTicket,
        adjustBar,
        deleteTicket,
        getUserListByDepartmentId,
        setScheduleByFromTo,
        setScheduleByPersonDay,
        ticketUserUpdate,
        updateDepartment,
        updateOrder,
        updateTicket,
    }
}


/**
 * 祝日を除く期間の日数を返却する
 * @param row
 * @param holidays
 */
function getNumberOfDays(row: GanttRow, holidays: Holiday[]) {
    const dayjsStartDate = dayjs(row.ticket?.start_date)
    const dayjsEndDate = dayjs(row.ticket?.end_date)
    const numberOfHolidays = holidays.filter(v => {
        return dayBetween(dayjs(v.date), dayjsStartDate, dayjsEndDate)
    })
    const numberOfDays = dayjsEndDate.diff(dayjsStartDate, 'day') + 1
    return numberOfDays - numberOfHolidays.length
}


const ticketToGanttRow = (ticket: Ticket, ticketUserList: TicketUser[], ganttGroup: GanttGroup) => {
    if (ticket.start_date != undefined) {
        ticket.start_date = ticket.start_date.substring(0, 10)
    }
    if (ticket.end_date != undefined) {
        ticket.end_date = ticket.end_date.substring(0, 10)
    }
    if (ticket.limit_date != undefined) {
        ticket.limit_date = ticket.limit_date.substring(0, 10)
    }

    const result: GanttRow = {
        bar: {
            beginDate: dayjs(ticket.start_date!).format(DAYJS_FORMAT),
            endDate: dayjs(ticket.end_date!).format(DAYJS_FORMAT),
            ganttBarConfig: {
                hasHandles: true,
                id: ticket.id!.toString(), // TODO: まあIDが数字でもいっか
                label: "", // TODO: コメントをラベルにする
            },
        },
        ticket: ticket,
        ticketUsers: ticketUserList.filter(v => v.ticket_id === ticket.id),
        ganttGroup: ganttGroup
    }
    return result
}


