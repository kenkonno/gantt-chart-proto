import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {ref} from "vue";

type TaskName = string;
type NumberOfWorkers = number;
type WorkStartDate = Dayjs;
type WorkEndDate = Dayjs;
type RowID = string;
type EstimatePersonDay = number;

const chartStart = ref("01.05.2023 00:00")
const chartEnd = ref("31.07.2023 00:00")
const format = ref("DD.MM.YYYY HH:mm")


type Row = {
    bar: GanttBarObject
    taskName: TaskName
    numberOfWorkers: NumberOfWorkers
    workStartDate: WorkStartDate
    workEndDate: WorkEndDate
    estimatePersonDay: EstimatePersonDay
}

export function useGantt() {
// ここからダミーデータ
    const rows = ref<Row[]>([
        newRow("id-1", "作業1", 1, 5, dayjs('2023-05-02 00:00'), dayjs('2023-05-06 23:59')),
        newRow("id-2", "作業2", 1, 6, dayjs('2023-05-07 00:00'), dayjs('2023-05-12 23:59')),
        newRow("id-3", "作業3", 1, 3, dayjs('2023-05-13 00:00'), dayjs('2023-05-15 23:59')),
        newRow("id-4", "作業4", 1, 2, dayjs('2023-05-16 00:00'), dayjs('2023-05-17 23:59')),
        newRow("id-5", "作業5", 2, 7, dayjs('2023-05-18 00:00'), dayjs('2023-05-24 23:59')),
        newRow("id-6", "作業6", 2, 10, dayjs('2023-05-25 00:00'), dayjs('2023-06-03 23:59'))
    ])
    const bars = rows.value.map(v => v.bar)
    const adjustBar = (bar: GanttBarObject) => {
        // id をもとにvuejs側のデータも更新する
        const target = rows.value.find(v => v.bar.ganttBarConfig.id === bar.ganttBarConfig.id)
        if (!target) {
            console.error(`ID: ${bar.ganttBarConfig.id} が現在のガントに存在しません。`, rows)
        } else {
            target.workStartDate = dayjs(ganttDateToYMDDate(bar.beginDate))
            target.workEndDate = dayjs(ganttDateToYMDDate(bar.endDate))
        }
    }

    return {
        rows,
        bars,
        chartStart,
        chartEnd,
        format,
        updateWorkStartDate,
        updateWorkEndDate,
        setScheduleByPersonDay,
        setScheduleByFromTo,
        adjustBar
    }
}

// 人日重視でスケジュールを設定する
// 人日*作業人数で期間を決定する
const setScheduleByPersonDay = (rows: Row[]) => {
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
const setScheduleByFromTo = (rows: Row[]) => {
    rows.forEach(row => {
        // 必要人数の算出。スケジュールを満了できる最少の人数を計算する。
        // 必要人数 = 見積人日 / 期間
        const period = (row.workEndDate.diff(row.workStartDate, 'day', false) + 1)
        console.log("### setScheduleByFromTo ############")
        console.log("### setScheduleByFromTo", row.estimatePersonDay)
        console.log("### setScheduleByFromTo", period)
        console.log("### setScheduleByFromTo", row.estimatePersonDay / period)
        row.numberOfWorkers = Math.ceil(row.estimatePersonDay /  period)
    })
}
const reflectGantt = (row: Row) => {
    row.bar.beginDate = row.workStartDate.format(format.value)
    row.bar.endDate = row.workEndDate.format(format.value)
}


const newRow = (id: RowID, taskName: TaskName, numberOfWorkers: NumberOfWorkers, estimatePersonDay: EstimatePersonDay, workStartDate: WorkStartDate, workEndDate: WorkEndDate) => {
    const result: Row = {
        bar: {
            beginDate: workStartDate.format(format.value),
            endDate: workEndDate.format(format.value),
            ganttBarConfig: {
                hasHandles: true,
                id: id,
                label: taskName,
            }
        },
        taskName: taskName,
        numberOfWorkers: numberOfWorkers,
        workStartDate: workStartDate,
        workEndDate: workEndDate,
        estimatePersonDay: estimatePersonDay
    }
    return result
}

// ここからデータ更新
const updateWorkStartDate = (row: Row, newValue: string) => {
    console.log(row, newValue)
    row.workStartDate = dayjs(newValue)
    row.bar.beginDate = row.workStartDate.format(format.value)
}
const updateWorkEndDate = (row: Row, newValue: string) => {
    console.log(row, newValue)
    row.workEndDate = dayjs(newValue).set('minute',59).set('hour',23)
    row.bar.endDate = row.workEndDate.format(format.value)
}


const ganttDateToYMDDate = (date: string) => {
    const e = date.split(" ")[0].split(".")
    return `${e[2]}-${e[1]}-${e[0]}`
}


