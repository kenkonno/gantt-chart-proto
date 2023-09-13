import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {ref} from "vue";
import {GanttGroup, Holiday, OperationSetting, Ticket, TicketUser} from "@/api";
import {Api} from "@/api/axios";
import {useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";
import {useTicketUserTable} from "@/composable/ticketUser";
import {useProcessTable} from "@/composable/process";
import {useFacility} from "@/composable/facility";

const format = ref("DD.MM.YYYY HH:mm")

type GanttChartGroup = {
    ganttGroup: GanttGroup // TODO: 結局ganttRowにganttGroup設定しているから設計としては微妙っぽい。
    rows: GanttRow[]
}

type GanttRow = {
    bar: GanttBarObject
    ganttGroup?: GanttGroup
    ticket?: Ticket
    ticketUsers?: TicketUser[]
}

const DEFAULT_PERSON_DAY = 8

type Header = {
    name: string,
    visible: boolean
}

const BAR_NORMAL_COLOR = "rgb(147 206 255)"
const BAR_COMPLETE_COLOR = "rgb(76 255 18)"
const BAR_DANGER_COLOR = "rgb(255 89 89);"

function getScheduledOperatingHours(operationSettings: OperationSetting[], row: GanttRow) {
    return operationSettings.filter(operationSetting => {
        return operationSetting.facility_id === row.ganttGroup?.facility_id &&
            operationSetting.unit_id === row.ganttGroup?.unit_id
    }).reduce((accumulateValue, currentValue) => {
        const workHours = currentValue.workHours.find(v => v.process_id === row.ticket?.process_id)
        return accumulateValue + workHours!.work_hour!
    }, 0);
}

export async function useGantt(facilityId: number) {
    const GanttHeader = ref<Header[]>([
        {name: "ユニット", visible: true},
        {name: "工程", visible: true},
        {name: "部署", visible: true},
        {name: "担当者", visible: true},
        {name: "人数", visible: true},
        {name: "期日", visible: true},
        {name: "工数(h)", visible: true},
        {name: "日後", visible: true},
        {name: "開始日", visible: true},
        {name: "終了日", visible: true},
        {name: "進捗", visible: true},
    ])
    const {facility} = await useFacility(facilityId)
    const chartStart = ref(dayjs(facility.value.term_from).format(format.value))
    const chartEnd = ref(dayjs(facility.value.term_to).format(format.value))

    // とりあえず何も考えずにAPIからガントチャート表示に必要なオブジェクトを作る
    const {list: ganttGroupList, refresh: ganttGroupRefresh} = await useGanttGroupTable()
    const {list: ticketList, refresh: ticketRefresh} = await useTicketTable()
    const {list: ticketUserList, refresh: ticketUserRefresh} = await useTicketUserTable()
    const {list: processList, refresh: processRefresh} = await useProcessTable()
    const processMap: { [x: number]: string; } = {}
    processList.value.forEach(v => {
        processMap[v.id!] = v.name
    })
    await ganttGroupRefresh(facilityId)
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
                beginDate: dayjs(task.ticket!.start_date!).format(format.value),
                endDate: endOfDay(task.ticket!.end_date!),
                ganttGroupId: unit.ganttGroup.id,
                ganttBarConfig: {
                    hasHandles: true,
                    id: task.ticket?.id!.toString(),
                    label: processMap[task.ticket?.process_id == null ? -1 : task.ticket?.process_id],
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

    const adjustBar = (bar: GanttBarObject) => {
        // id をもとにvuejs側のデータも更新する
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
                updateTicket(targetTicket.ticket!)
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
            targetTicket.beginDate = dayjs(ticket.start_date!).format(format.value)
            targetTicket.endDate = endOfDay(ticket.end_date!)
            if (ticket.progress_percent) {
                targetTicket.ganttBarConfig.progress = ticket.progress_percent
            }
            targetTicket.ganttBarConfig.label = processMap[ticket.process_id == null ? -1 : ticket.process_id]
        }
    }

    /**
     * setScheduleByPersonDay
     * リスケ（工数h）重視
     * 稼働予定時間ベースでスケジュールを引き直す。
     *
     * @param rows
     * @param holidays
     * @param operationSettings
     * 色々考えたけど難しい（答えがない）のでマジでシンプルに１行ずつ処理するようにする。
     *
     * 基本的な仕様
     * ・人日重視で期間を変更する。
     * ・開始日・工数が設定されていない行は無視する
     * ・稼働の消化は担当者の設定順とする。
     * とりあえず作らないと話にならないので一旦作る
     */
    const setScheduleByPersonDay = (rows: GanttRow[], holidays: Holiday[], operationSettings: OperationSetting[]) => {
        let prevEndDate: Dayjs
        rows.forEach(row => {
            let startDate: Dayjs
            // 開始日・工数・担当者が未設定の場合はスキップする
            if (!row.ticket?.start_date || !row.ticket?.estimate || !row.ticketUsers?.length) {
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
            row.ticket.start_date = ganttDateToYMDDate(startDate.format(format.value))
            // 終了日を決定する、祝日が含まれている場合終了日をずらす。
            row.ticket.end_date = ganttDateToYMDDate(
                getEndDateByRequiredBusinessDay(startDate, numberOfDaysRequired, holidays)
                    .add(-1, 'minute').format(format.value)
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
     * @param holidays
     * @param operationSettings
     */
    const setScheduleByFromTo = (rows: GanttRow[], holidays: Holiday[], operationSettings: OperationSetting[]) => {
        // 必要人数の算出。スケジュールを満了できる最少の人数を計算する。
        // 担当者が未定の時に人数を計算する
        // 必要日数 = 工数 / 労働時間
        // 必要人数 = 必要日数 / 期間日数 の数値は上に寄せる
        let prevEndDate: Dayjs
        rows.forEach(row => {
            if (row.ticket && row.ticketUsers?.length == 0 &&
                row.ticket?.estimate && row.ticket?.estimate > 0 &&
                row.ticket?.start_date && row.ticket?.end_date) {
                // 日後の処理
                if (row.ticket.days_after != null && row.ticket.days_after >= 0 && prevEndDate != null) {
                    // 現在設定されている営業日を計算する
                    let dayjsStartDate = dayjs(row.ticket.start_date)
                    let dayjsEndDate = dayjs(row.ticket.end_date)
                    const requiredNumberOfBusinessDays = getNumberOfBusinessDays(dayjsStartDate, dayjsEndDate, holidays)
                    // 開始日を日後分ずらした日付に設定する
                    dayjsStartDate = prevEndDate.add(row.ticket.days_after, 'days')
                    // 終了日を設定されている営業日分確保する
                    dayjsEndDate = getEndDateByRequiredBusinessDay(dayjsStartDate, requiredNumberOfBusinessDays, holidays)
                    // チケットに反映させる
                    row.ticket.start_date = ganttDateToYMDDate(dayjsStartDate.format(format.value))
                    row.ticket.end_date = ganttDateToYMDDate(dayjsEndDate.format(format.value))
                }

                // 工数(h)
                const estimate = row.ticket.estimate
                // 労働予定時間 ユニット & 工程 に紐づく稼働予定時間を取得
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

    // ####################### ここから下は昔の資産 ###############################
    const footerLabels = ref<string[]>([])
    // 積み上げの更新
    // footerLabels.value.push(...calculateStuckPersons(oldRows.value))

    //
    // watch(oldRows, () => {
    //     console.log("WATCH")
    //     footerLabels.value.splice(0)
    //     footerLabels.value.push(...calculateStuckPersons(oldRows.value))
    //
    //     // rowsの変更を barsに反映させる。色々散らばっているが、ここでは時刻関係以外とする。
    //     oldBars.splice(0)
    //     oldBars.push(...oldRows.value.map(v => {
    //         v.bar.ganttBarConfig.label = v.taskName
    //         return v.bar
    //     }))
    //
    // }, {deep: true})

    return {
        chartStart,
        chartEnd,
        ganttChartGroup,
        bars,
        footerLabels,
        format,
        GanttHeader,
        setScheduleByPersonDay,
        setScheduleByFromTo,
        adjustBar,
        slideSchedule,
        addNewTicket,
        updateTicket,
        ticketUserUpdate,
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


// スケジュールをきれいにスライドさせる（重複をなくす）
const slideSchedule = (rows: GanttRow[]) => {
    // let prevEndDate: Dayjs
    // rows.forEach((row, i) => {
    //     // 1回目は無視
    //     if (i > 0) {
    //         // 期間を計算する
    //         const duration = row.workEndDate.diff(row.workStartDate, 'day')
    //         row.workStartDate = prevEndDate.add(1, 'day').set('hour', 0).set('minute', 0)
    //         row.workEndDate = row.workStartDate.add(duration, 'day').set('hour', 23).set('minute', 59)
    //         reflectGantt(row)
    //     }
    //     prevEndDate = row.workEndDate
    // })
}

// 人数の積み上げを行う
const calculateStuckPersons = (rows: GanttRow[]) => {
    // console.log("calculateStuckPersons start")
    //
    // const result: number[] = []
    // // 開始時刻
    // let currentDate = dayjs(ganttDateToYMDDate(chartStart.value))
    // const endDate = dayjs(ganttDateToYMDDate(chartEnd.value))
    // // 結果セットの初期化
    // while (currentDate.isBefore(endDate)) {
    //     currentDate = currentDate.add(1, 'day')
    //     result.push(0)
    // }
    //
    // let seq = 0
    // currentDate = dayjs(ganttDateToYMDDate(chartStart.value))
    // while (currentDate.isBefore(endDate)) {
    //     // 開始日が同じ行を探す
    //     const targets = rows.filter(v => v.workStartDate.isSame(currentDate))
    //     if (targets.length > 0) {
    //         targets.forEach(target => {
    //             let innerSeq = seq
    //             let innerCurrentDate = currentDate
    //             let estimatePersonDay = target.estimatePersonDay
    //             // 作業人数
    //             const numberOfWorkers = target.numberOfWorkers
    //
    //             while (estimatePersonDay > 0) {
    //                 estimatePersonDay = estimatePersonDay - numberOfWorkers
    //                 if (estimatePersonDay < 0) {
    //                     result[innerSeq] += numberOfWorkers + estimatePersonDay
    //                 } else {
    //                     result[innerSeq] += numberOfWorkers
    //                 }
    //                 innerCurrentDate = innerCurrentDate.add(1, 'day')
    //                 innerSeq++
    //             }
    //         })
    //     }
    //     currentDate = currentDate.add(1, 'day')
    //     seq++
    // }
    //
    // return result.map(v => v === 0 ? '' : v.toString())
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
            beginDate: dayjs(ticket.start_date!).format(format.value),
            endDate: dayjs(ticket.end_date!).format(format.value),
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

const ganttDateToYMDDate = (date: string) => {
    const e = date.split(" ")[0].split(".")
    return `${e[2]}-${e[1]}-${e[0]}`
}


const endOfDay = (dateString: string) => {
    return dayjs(dateString).add(1, 'day').add(-1, 'minute').format(format.value)
}
/**
 * 開始日・終了日に祝日が含まれていれば終了日をその分ずらす
 * @param startDate
 * @param endDate
 * @param holidays
 */
const getNumberOfBusinessDays = (startDate: Dayjs, endDate: Dayjs, holidays: Holiday[]) => {
    // 差がない場合は1営業日になる
    const result = endDate.diff(startDate,'days') + 1
    // 開始日と終了日の内、祝日を除いた日数を返す
    const includes = holidays.filter(holiday => {
        const dayjsHoliday = dayjs(holiday.date)
        return dayBetween(dayjsHoliday, startDate, endDate)
    })
    return result + includes.length
}

const getEndDateByRequiredBusinessDay = (startDate: Dayjs, requiredNumberOfBusinessDays: number, holidays: Holiday[]) => {
    let currentDate = startDate
    const dayjsHolidays = holidays.map(v => dayjs(v.date))
    // 1営業日だとしたらstartDateを返せばよい
    while (requiredNumberOfBusinessDays > 1) {
        currentDate = currentDate.add(1, 'day')
        if (dayjsHolidays.some(v => v.isSame(currentDate))) continue
        requiredNumberOfBusinessDays--
    }
    return currentDate.endOf('day')
}

const dayBetween = (day: Dayjs, form: Dayjs, to: Dayjs) => {
    return (form.isSame(day) || to.isSame(day)) ||
        (day.isAfter(form) && day.isBefore(to))
}
const endCheckHoliday = (holiday: Dayjs, endDate: Dayjs) => {
    return holiday.isBefore(endDate) && endDate.isBefore(holiday.add(1, 'day'))
}

/**
 * 開始日が祝日だった場合開始日をずらす
 * @param startDate
 * @param holidays
 */
const adjustStartDateByHolidays = (startDate: Dayjs, holidays: Holiday[]) => {
    let result = startDate
    let endCheck = holidays.find(v => dayjs(v.date).isSame(result))
    while (endCheck) {
        result = result.add(1, "day")
        endCheck = holidays.find(v => dayjs(v.date).isSame(result))
    }
    return result
}