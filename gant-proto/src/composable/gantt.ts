import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {computed, ref, watch} from "vue";
import {GanttGroup, Ticket, TicketUser} from "@/api";
import {Api} from "@/api/axios";
import {useGanttGroupTable} from "@/composable/ganttGroup";
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
    }
    // DBへのストア及びローカルのガントに情報を反映する
    const ticketUserUpdate = async (ticketId: number, userIds: number[]) => {
        const {data} = await Api.postTicketUsers({ticketId: ticketId, userIds: userIds})
        // ticketUserList から TicketId をつけなおす
        const newTicketUserList = ticketUserList.value.filter(v => v.ticket_id !== ticketId)
        newTicketUserList.push(...data.ticketUsers)
        ticketUserList.value.length = 0
        ticketUserList.value.push(...newTicketUserList)
        // TODO: refreshLocalGanttは重くなるので性能対策が必要かも
        refreshLocalGantt()
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
                ganttBarConfig: {id: Math.random().toString()}
            }
            bars.value.push(...unit.rows.map(task => <GanttBarObject>{
                beginDate: dayjs(task.ticket!.start_date!).format(format.value),
                endDate: dayjs(task.ticket!.end_date!).format(format.value),
                ganttGroupId: unit.ganttGroup.id,
                ganttBarConfig: {
                    hasHandles: true,
                    id: task.ticket?.id!.toString(),
                    label: "",
                }
            }))
            // ボタン用の空行を追加する
            bars.value.push(emptyRow)
        })
    }
    refreshBars()

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
            targetTicket.beginDate = dayjs(ticket.start_date!).format(format.value)
            targetTicket.endDate = dayjs(ticket.end_date!).format(format.value)
        }
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
        setScheduleByPersonDay,
        setScheduleByFromTo,
        adjustBar,
        slideSchedule,
        addNewTicket,
        updateTicket,
        ticketUserUpdate,
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
        taskName: "taskName",
        numberOfWorkers: 1,
        workStartDate: dayjs(),
        workEndDate: dayjs(),
        estimatePersonDay: 1,
    }
    return result
}

const ganttDateToYMDDate = (date: string) => {
    const e = date.split(" ")[0].split(".")
    return `${e[2]}-${e[1]}-${e[0]}`
}


