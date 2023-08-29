import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {Ref, ref, UnwrapRef, watch} from "vue";
import {GanttGroup, Ticket, TicketUser, Unit} from "@/api";
import {Api} from "@/api/axios";
import {useGanttGroup, useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";
import {useTicketUserTable} from "@/composable/ticketUser";

type TaskName = string;
type NumberOfWorkers = number;
type WorkStartDate = Dayjs;
type WorkEndDate = Dayjs;
type RowID = string;
type EstimatePersonDay = number;

const chartStart = ref("01.05.2023 00:00")
const chartEnd = ref("31.07.2023 00:00")
const format = ref("DD.MM.YYYY HH:mm")

type GanttChartGroup = {
    ganttGroup: GanttGroup
    rows: GanttRow[]
}

type GanttRow = {
    bar: GanttBarObject
    ganttGroup?: GanttGroup
    ticket?: Ticket
    ticketUsers?: TicketUser[]
    // TODO: ここから下は後で消す
    taskName: TaskName
    numberOfWorkers: NumberOfWorkers
    workStartDate: WorkStartDate
    workEndDate: WorkEndDate
    estimatePersonDay: EstimatePersonDay
}

export async function useGantt(facilityId: number) {
    // TODO: マスタ編集系と同期がとれていない。global store 戦略のほうが良いかも。

    // とりあえず何も考えずにAPIからガントチャート表示に必要なオブジェクトを作る
    const {list: ganttGroupList, refresh: ganttGroupRefresh} = await useGanttGroupTable()
    const {list: ticketList, refresh: ticketRefresh} = await useTicketTable()
    const {list: ticketUserList, refresh: ticketUserRefresh} = await useTicketUserTable()
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
                    rows: ticketList.value.filter(ticket => ticket.gantt_group_id === ganttGroup.id).sort(v => v.order).map(ticket => ticketToGanttRow(ticket, ticketUserList.value))
                }
            )
        })
    }
    refreshLocalGantt()
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
    }

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
    }
    const ticketUserUpdate = async (ticketId: number, userIds: number[]) => {
        const {data} = await Api.postTicketUsers({ticketId: ticketId, userIds: userIds})
        // ticketUserList から TicketId をつけなおす
        const newTicketUserList = ticketUserList.value.filter(v => v.ticket_id !== ticketId)
        newTicketUserList.push(...data.ticketUsers)
        ticketUserList.value.length = 0
        ticketUserList.value.push(...newTicketUserList)
        refreshLocalGantt()

    }

    // ####################### ここから下は昔の資産 ###############################

    // ここからダミーデータ
    const rows = ref<GanttRow[]>([
        newRow("id-1", "作業1", 1, 5, dayjs('2023-05-02 00:00'), dayjs('2023-05-06 23:59')),
        newRow("id-2", "作業2", 1, 6, dayjs('2023-05-07 00:00'), dayjs('2023-05-12 23:59')),
        newRow("id-3", "作業3", 1, 3, dayjs('2023-05-13 00:00'), dayjs('2023-05-15 23:59')),
        newRow("id-4", "作業4", 1, 2, dayjs('2023-05-16 00:00'), dayjs('2023-05-17 23:59')),
        newRow("id-5", "作業5", 2, 7, dayjs('2023-05-18 00:00'), dayjs('2023-05-24 23:59')),
        newRow("id-6", "作業6", 2, 10, dayjs('2023-05-25 00:00'), dayjs('2023-06-03 23:59'))
    ])
    const bars = rows.value.map(v => v.bar)
    const footerLabels = ref<string[]>([])
    // 積み上げの更新
    footerLabels.value.push(...calculateStuckPersons(rows.value))

    // リアクティブな関数群
    const adjustBar = (bar: GanttBarObject) => {
        // id をもとにvuejs側のデータも更新する
        const target = rows.value.find(v => v.bar.ganttBarConfig.id === bar.ganttBarConfig.id)
        if (!target) {
            console.error(`ID: ${bar.ganttBarConfig.id} が現在のガントに存在しません。`, rows)
        } else {
            // target.workStartDate = dayjs(ganttDateToYMDDate(bar.beginDate))
            // target.workEndDate = dayjs(ganttDateToYMDDate(bar.endDate))
        }
    }
    const addRow = () => {
        rows.value.push(newRow(
            Math.random().toString(),
            "新規タスク",
            1,
            1,
            dayjs(ganttDateToYMDDate(chartStart.value)),
            dayjs(ganttDateToYMDDate(chartStart.value)).add(3, 'day'),
        ))
    }

    const deleteRow = (id: RowID) => {
        console.log("#deleterow")
        rows.value.forEach((row, i) => {
            if (row.bar.ganttBarConfig.id === id) {
                rows.value.splice(i, 1)
            }
        })
    }

    watch(rows, () => {
        console.log("WATCH")
        footerLabels.value.splice(0)
        footerLabels.value.push(...calculateStuckPersons(rows.value))

        // rowsの変更を barsに反映させる。色々散らばっているが、ここでは時刻関係以外とする。
        bars.splice(0)
        bars.push(...rows.value.map(v => {
            v.bar.ganttBarConfig.label = v.taskName
            return v.bar
        }))

    }, {deep: true})

    return {
        rows,
        bars,
        footerLabels,
        chartStart,
        chartEnd,
        format,
        ganttChartGroup,
        updateWorkStartDate,
        updateWorkEndDate,
        setScheduleByPersonDay,
        setScheduleByFromTo,
        adjustBar,
        slideSchedule,
        addRow,
        deleteRow,
        addNewTicket,
        updateTicket,
        ticketUserUpdate
    }
}

// 人日重視でスケジュールを設定する
// 人日*作業人数で期間を決定する
const setScheduleByPersonDay = (rows: GanttRow[]) => {
    let currentDate = rows[0].workStartDate
    rows.forEach(row => {
        // 必要日数。少数が出たら最小の整数にする。例：二人で5人日の場合 3日必要 5/2 = 2.5 つまり少数が出たら整数に足す
        const numberOfDaysRequired = Math.ceil(row.estimatePersonDay / row.numberOfWorkers)
        // 開始日を設定する
        row.workStartDate = currentDate
        // 終了日を決定する
        row.workEndDate = row.workStartDate.add(numberOfDaysRequired, 'day').add(-1, 'minute')
        // ガントに反映する
        reflectGantt(row)
        // 現在日付を設定する
        currentDate = row.workEndDate.add(1, 'minute')
    })
}
// 開始日・終了日重視で人数を設定する
// 期間で満了できるように人数を割り当てる
const setScheduleByFromTo = (rows: GanttRow[]) => {
    rows.forEach(row => {
        // 必要人数の算出。スケジュールを満了できる最少の人数を計算する。
        // 必要人数 = 見積人日 / 期間
        const period = (row.workEndDate.diff(row.workStartDate, 'day', false) + 1)
        row.numberOfWorkers = Math.ceil(row.estimatePersonDay / period)
    })
}
// スケジュールをきれいにスライドさせる（重複をなくす）
const slideSchedule = (rows: GanttRow[]) => {
    let prevEndDate: Dayjs
    rows.forEach((row, i) => {
        // 1回目は無視
        if (i > 0) {
            // 期間を計算する
            const duration = row.workEndDate.diff(row.workStartDate, 'day')
            row.workStartDate = prevEndDate.add(1, 'day').set('hour', 0).set('minute', 0)
            row.workEndDate = row.workStartDate.add(duration, 'day').set('hour', 23).set('minute', 59)
            reflectGantt(row)
        }
        prevEndDate = row.workEndDate
    })
}

// 人数の積み上げを行う
const calculateStuckPersons = (rows: GanttRow[]) => {
    console.log("calculateStuckPersons start")

    const result: number[] = []
    // 開始時刻
    let currentDate = dayjs(ganttDateToYMDDate(chartStart.value))
    const endDate = dayjs(ganttDateToYMDDate(chartEnd.value))
    // 結果セットの初期化
    while (currentDate.isBefore(endDate)) {
        currentDate = currentDate.add(1, 'day')
        result.push(0)
    }

    let seq = 0
    currentDate = dayjs(ganttDateToYMDDate(chartStart.value))
    while (currentDate.isBefore(endDate)) {
        // 開始日が同じ行を探す
        const targets = rows.filter(v => v.workStartDate.isSame(currentDate))
        if (targets.length > 0) {
            targets.forEach(target => {
                let innerSeq = seq
                let innerCurrentDate = currentDate
                let estimatePersonDay = target.estimatePersonDay
                // 作業人数
                const numberOfWorkers = target.numberOfWorkers

                while (estimatePersonDay > 0) {
                    estimatePersonDay = estimatePersonDay - numberOfWorkers
                    if (estimatePersonDay < 0) {
                        result[innerSeq] += numberOfWorkers + estimatePersonDay
                    } else {
                        result[innerSeq] += numberOfWorkers
                    }
                    innerCurrentDate = innerCurrentDate.add(1, 'day')
                    innerSeq++
                }
            })
        }
        currentDate = currentDate.add(1, 'day')
        seq++
    }

    return result.map(v => v === 0 ? '' : v.toString())
}

const reflectGantt = (row: GanttRow) => {
    row.bar.beginDate = row.workStartDate.format(format.value)
    row.bar.endDate = row.workEndDate.format(format.value)
}


const newRow = (id: RowID, taskName: TaskName, numberOfWorkers: NumberOfWorkers, estimatePersonDay: EstimatePersonDay, workStartDate: WorkStartDate, workEndDate: WorkEndDate) => {
    const result: GanttRow = {
        bar: {
            beginDate: workStartDate.format(format.value),
            endDate: workEndDate.format(format.value),
            ganttBarConfig: {
                hasHandles: true,
                id: id,
                label: taskName,
            },
        },
        taskName: taskName,
        numberOfWorkers: numberOfWorkers,
        workStartDate: workStartDate,
        workEndDate: workEndDate,
        estimatePersonDay: estimatePersonDay,
    }
    return result
}

const ticketToGanttRow = (ticket: Ticket, ticketUserList: TicketUser[]) => {
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
            beginDate: ticket.start_date!,
            endDate: ticket.end_date!,
            ganttBarConfig: {
                hasHandles: true,
                id: ticket.id!.toString(), // TODO: まあIDが数字でもいっか
                label: "", // TODO: コメントをラベルにする
            },
        },
        ticket: ticket,
        ticketUsers: ticketUserList.filter(v => v.ticket_id === ticket.id),
        taskName: "taskName",
        numberOfWorkers: 1,
        workStartDate: dayjs(),
        workEndDate: dayjs(),
        estimatePersonDay: 1,
    }
    return result
}

// ここからデータ更新
const updateWorkStartDate = (row: GanttRow, newValue: string) => {
    console.log(row, newValue)
    row.workStartDate = dayjs(newValue)
    row.bar.beginDate = row.workStartDate.format(format.value)
}
const updateWorkEndDate = (row: GanttRow, newValue: string) => {
    console.log(row, newValue)
    row.workEndDate = dayjs(newValue).set('minute', 59).set('hour', 23)
    row.bar.endDate = row.workEndDate.format(format.value)
}


const ganttDateToYMDDate = (date: string) => {
    const e = date.split(" ")[0].split(".")
    return `${e[2]}-${e[1]}-${e[0]}`
}


