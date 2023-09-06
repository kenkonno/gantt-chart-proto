<template>
  <input type="button" value="人日重視で設定する" @click="setScheduleByPersonDayProxy()">
  <input type="button" value="スケジュール重視で設定する" @click="setScheduleByFromTo(oldRows)">
  <input type="button" value="スケジュールをスライドする" @click="slideSchedule(oldRows)">
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
    <template #side-menu>
      <table class="side-menu">
        <thead class="side-menu-header">
        <tr>
          <th class="side-menu-cell"></th><!-- css hack min-height -->
          <th class="side-menu-cell">ユニット</th>
          <th class="side-menu-cell">工程</th>
          <th class="side-menu-cell">部署</th>
          <th class="side-menu-cell">担当者</th>
          <th class="side-menu-cell">期日</th>
          <th class="side-menu-cell">工数</th>
          <th class="side-menu-cell">日後</th>
          <th class="side-menu-cell">開始日</th>
          <th class="side-menu-cell">終了日</th>
          <th class="side-menu-cell">進捗</th>
        </tr>
        </thead>
        <tbody>
        <template v-for="item in ganttChartGroup" :key="item.ganttGroup.id">
          <tr v-for="row in item.rows" :key="row.ticket.id">
            <td class="side-menu-cell"></td><!-- css hack min-height -->
            <td class="side-menu-cell">{{ unitMap[item.ganttGroup.unit_id] }}</td>
            <td class="side-menu-cell">
              <select v-model="row.ticket.process_id" @change="updateTicket(row.ticket)">
                <option v-for="item in processList" :key="item.id" :value="item.id">{{ item.name }}</option>
              </select>
            </td>
            <td class="side-menu-cell">
              <select v-model="row.ticket.department_id" @change="updateTicket(row.ticket)">
                <option v-for="item in departmentList" :key="item.id" :value="item.id">{{ item.name }}</option>
              </select>
            </td>
            <td class="side-menu-cell" style="width: 14rem;">
              <UserMultiselect :userList="userList" :ticketUser="row.ticketUsers"
                               @update="ticketUserUpdate(row.ticket.id ,$event)"></UserMultiselect>
            </td>
            <td class="side-menu-cell"><input type="date" v-model="row.ticket.limit_date"
                                              @change="updateTicket(row.ticket)"/></td>
            <td class="side-menu-cell">
              <FormNumber class="small-numeric" v-model="row.ticket.estimate" @change="updateTicket(row.ticket)"/>
            </td>
            <td class="side-menu-cell">
              <FormNumber class="small-numeric" v-model="row.ticket.days_after" @change="updateTicket(row.ticket)"/>
            </td>
            <td class="side-menu-cell"><input type="date" v-model="row.ticket.start_date"
                                              @change="updateTicket(row.ticket)"/></td>
            <td class="side-menu-cell"><input type="date" v-model="row.ticket.end_date"
                                              @change="updateTicket(row.ticket)"/></td>
            <td class="side-menu-cell">
              <FormNumber class="middle-numeric" v-model="row.ticket.progress_percent"
                          @change="updateTicket(row.ticket)"/>
            </td>
          </tr>
          <tr>
            <td colspan="11">
              <button @click="addNewTicket(item.ganttGroup.id)">{{
                  unitMap[item.ganttGroup.unit_id]
                }}の工程を追加する
              </button>
            </td>
          </tr>
        </template>
        </tbody>
      </table>
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

.side-menu {
  white-space: nowrap;
  text-align: center;

  td, th {
    border: solid 1px #eaeaea;
    box-sizing: border-box;
  }

  tr {
    height: 40px;
  }

  .side-menu-header > tr > th:first-child {
    display: block;
    float: left;
    content: "";
    min-height: 75px;
    height: 8vh;
  }

  .small-numeric {
    width: 3rem;
  }

  .middle-numeric {
    width: 4rem;
  }
}

</style>

<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {useGantt} from "@/composable/gantt";
import {useUnitTable} from "@/composable/unit";
import {useProcessTable} from "@/composable/process";
import {useDepartment, useDepartmentTable} from "@/composable/department";
import {useUserTable} from "@/composable/user";
import FormNumber from "@/components/form/FormNumber.vue";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import {useHolidayTable} from "@/composable/holiday";
import {useOperationSettingTable} from "@/composable/operationSetting";

const {
  chartStart,
  chartEnd,
  format,
  footerLabels,
  setScheduleByPersonDay,
  setScheduleByFromTo,
  adjustBar,
  slideSchedule,
  ganttChartGroup,
  addNewTicket,
  updateTicket,
  ticketUserUpdate,
  bars,
} = await useGantt(10) // TODO: facilityId

const {list: unitList, refresh: unitRefresh} = await useUnitTable(10) // TODO: facilityId
const {list: processList, refresh: processRefresh} = await useProcessTable()
const {list: departmentList, refresh: departmentRefresh} = await useDepartmentTable()
const {list: userList, refresh: userRefresh} = await useUserTable()
const {list: holidayList, refresh: holidayRefresh} = await useHolidayTable(10) // TODO: facilityId
const {list: operationSettingList, refresh: operationSettingRefresh} = await useOperationSettingTable(10) // TODO: facilityId

const setScheduleByPersonDayProxy = () => {
  ganttChartGroup.value.forEach(v => {
    setScheduleByPersonDay(v.rows, holidayList.value, operationSettingList.value)
  })
}

const unitMap: { [x: number]: string; } = {}
unitList.value.forEach(v => {
  unitMap[v.id!] = v.name
})
const processMap: { [x: number]: string; } = {}
processList.value.forEach(v => {
  processMap[v.id!] = v.name
})
const departmentMap: { [x: number]: string; } = {}
departmentList.value.forEach(v => {
  departmentMap[v.id!] = v.name
})
const userMap: { [x: number]: string; } = {}
userList.value.forEach(v => {
  userMap[v.id!] = v.name
})
const userOptions = userList.value.map(v => {
  return {id: v.id!, name: v.name!}
})

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
