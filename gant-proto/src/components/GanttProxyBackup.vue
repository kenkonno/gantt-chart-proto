<template>
  <input type="button" value="人日重視で設定する" @click="setScheduleByPersonDay(rows)">
  <input type="button" value="スケジュール重視で設定する" @click="setScheduleByFromTo(rows)">
  <input type="button" value="スケジュールをスライドする" @click="slideSchedule(rows)">
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

      :footer-labels="footerLabels"
      sticky
  >
    <g-gantt-row v-for="bar in bars" :key="bar.ganttBarConfig.id" :bars="[bar]"/>
    <template #table-header>
      table header
    </template>
    <template #rows>
      <div class="row-wrapper">
        <div class="g-gantt-row" v-for="row in rows" :key="row.bar.ganttBarConfig.id"
             style="height: 40px; align-items: center;">
          <div><input type="text" v-model="row.taskName" /></div>
          <div><input type="number" v-model="row.numberOfWorkers" style="width: 2rem"/>人</div>
          <div><input type="number" v-model="row.estimatePersonDay" style="width: 2rem"/>人日</div>
          <div><input type="date" :value="row.workStartDate.format('YYYY-MM-DD')"
                      @input="updateWorkStartDate(row, $event.target.value)"/></div>
          <div><input type="date" :value="row.workEndDate.format('YYYY-MM-DD')"
                      @input="updateWorkEndDate(row, $event.target.value)"/></div>
          <div><a href="#" @click="deleteRow(row.bar.ganttBarConfig.id)">削除</a></div>
        </div>
        <div>
          <input type="button" value="行の追加" @click="addRow()">
        </div>
      </div>
    </template>
  </g-gantt-chart>
</template>
<style lang="scss" scoped>
.row-wrapper {
  display: table;
  > .g-gantt-row {
    display: table-row;

    > div {
      padding: 0.5rem;
      white-space: nowrap;
      display: table-cell;
    }
  }

}
</style>

<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {useGanttFacility} from "@/composable/ganttFacility";

const {
  rows,
  bars,
  chartStart,
  chartEnd,
  format,
  footerLabels,
  updateWorkStartDate,
  updateWorkEndDate,
  setScheduleByPersonDay,
  setScheduleByFromTo,
  adjustBar,
  slideSchedule,
  addRow,
  deleteRow
} = useGanttFacility()


// ここからイベントフック
const onClickBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("click-bar", bar, e, datetime)
  adjustBar(bar)

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
