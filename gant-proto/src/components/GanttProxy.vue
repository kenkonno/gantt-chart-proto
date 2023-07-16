<template>
  <g-gantt-chart
      :chart-start="chartStart"
      :chart-end="chartEnd"
      precision="day"
      :row-height="40"
      grid
      width="1800px"
      bar-start="beginDate"
      bar-end="endDate"
      :date-format="format"
      @click-bar="onClickBar($event.bar, $event.e, $event.datetime)"
      @mousedown-bar="onMousedownBar($event.bar, $event.e, $event.datetime)"
      @dblclick-bar="onMouseupBar($event.bar, $event.e, $event.datetime)"
      @mouseenter-bar="onMouseenterBar($event.bar, $event.e)"
      @mouseleave-bar="onMouseleaveBar($event.bar, $event.e)"
      @dragstart-bar="onDragstartBar($event.bar, $event.e)"
      @drag-bar="onDragBar($event.bar, $event.e)"
      @dragend-bar="onDragendBar($event.bar, $event.e)"
      @contextmenu-bar="onContextmenuBar($event.bar, $event.e, $event.datetime)"
  >
    <g-gantt-row v-for="bar in bars" :key="bar.ganttBarConfig.id" :bars="[bar]"/>
    <template #rows>
      <div class="g-gantt-row" v-for="row in rows" :key="row.bar.ganttBarConfig.id"
           style="height: 40px; display:flex; align-items: center;">
        <div style="white-space: nowrap;">{{ row.taskName }}</div>
        <div><input type="number" v-model="row.numberOfWorkers" style="width: 2rem"/></div>
        <div><input type="date" :value="row.workStartDate.format('YYYY-MM-DD')"
                    @input="updateWorkStartDate(row, $event.target.value)"/></div>
        <div><input type="date" :value="row.workEndDate.format('YYYY-MM-DD')"
                    @input="updateWorkEndDate(row, $event.target.value)"/></div>
      </div>
    </template>
  </g-gantt-chart>
</template>
<style scss scoped>
.g-gantt-row > div {
  padding: 0.5rem;
}
</style>

<script setup lang="ts">
import {ref} from "vue"

import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject, GGanttRow} from "@infectoone/vue-ganttastic";

const chartStart = ref("01.05.2023 00:00")
const chartEnd = ref("31.07.2023 00:00")
const format = ref("DD.MM.YYYY HH:mm")

type TaskName = string;
type NumberOfWorkers = number;
type WorkStartDate = Dayjs;
type WorkEndDate = Dayjs;
type RowID = string;

type Row = {
  bar: GanttBarObject
  taskName: TaskName
  numberOfWorkers: NumberOfWorkers
  workStartDate: WorkStartDate
  workEndDate: WorkEndDate
}

const newRow = (id: RowID, taskName: TaskName, numberOfWorkers: NumberOfWorkers, workStartDate: WorkStartDate, workEndDate: WorkEndDate) => {
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
  row.workEndDate = dayjs(newValue)
  row.bar.endDate = row.workEndDate.format(format.value)
}

const ganttDateToYMDDate = (date: string) => {
  const e = date.split(" ")[0].split(".")
  return `${e[2]}-${e[1]}-${e[0]}`
}


// ここからダミーデータ
const rows = ref<Row[]>([
  newRow("id-1", "作業1", 1, dayjs('2023-05-02'), dayjs('2023-05-03')),
  newRow("id-2", "作業2", 1, dayjs('2023-05-06'), dayjs('2023-05-12'))
])
const bars = rows.value.map(v => v.bar)

// ここからイベントフック
const onClickBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("click-bar", bar, e, datetime)
  // １日の始まりに合わせる操作
  bar.beginDate = bar.beginDate.substring(0, 11) + "00:00"
  bar.endDate = bar.endDate.substring(0, 11) + "00:00"
  // id をもとにvuejs側のデータも更新する
  const target = rows.value.find(v => v.bar.ganttBarConfig.id === bar.ganttBarConfig.id)
  if (!target) {
    console.error(`ID: ${bar.ganttBarConfig.id} が現在のガントに存在しません。`, rows)
  } else {
    target.workStartDate = dayjs(ganttDateToYMDDate(bar.beginDate))
    target.workEndDate = dayjs(ganttDateToYMDDate(bar.endDate))
  }

}

const onMousedownBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("mousedown-bar", bar, e, datetime)
}

const onMouseupBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("mouseup-bar", bar, e, datetime)
}

const onMouseenterBar = (bar: GanttBarObject, e: MouseEvent) => {
  console.log("mouseenter-bar", bar, e)
}

const onMouseleaveBar = (bar: GanttBarObject, e: MouseEvent) => {
  console.log("mouseleave-bar", bar, e)
}

const onDragstartBar = (bar: GanttBarObject, e: MouseEvent) => {
  console.log("dragstart-bar", bar, e)
}

const onDragBar = (bar: GanttBarObject, e: MouseEvent) => {
  console.log("drag-bar", bar, e)
}

const onDragendBar = (
    bar: GanttBarObject,
    e: MouseEvent,
) => {
  console.log("dragend-bar", bar, e)
}

const onContextmenuBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("contextmenu-bar", bar, e, datetime)
}
</script>
