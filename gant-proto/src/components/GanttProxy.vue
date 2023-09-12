<template>
  <div v-if="operationSettingList.length > 0" id="gantt-wrapper">
    <div class="action-menu d-flex">
      <div class="wrapper d-flex">
        <div class="justify-middle">
          <div>メニュー</div>
        </div>
      </div>
      <AccordionHorizontal class="justify-middle">
        <template v-slot:icon>
          <span class="material-symbols-outlined">menu_open</span>
        </template>
        <template v-slot:body>
          <input type="button" class="btn btn-sm btn-outline-dark" value="人日重視で設定する"
                 @click="setScheduleByPersonDayProxy()">
          <input type="button" class="btn btn-sm btn-outline-dark" value="スケジュール重視で設定する"
                 @click="setScheduleByFromTo(oldRows)">
        </template>
      </AccordionHorizontal>
      <AccordionHorizontal class="justify-middle">
        <template v-slot:icon>
          <span class="material-symbols-outlined">filter_list</span>
        </template>
        <template v-slot:body>
          <div class="justify-middle">
            <div class="filter">
              <label v-for="item in GanttHeader" :key="item" class="side-menu-cell">
                <input type="checkbox" v-model="item.visible"/>{{ item.name }}
              </label>
            </div>
          </div>
        </template>
      </AccordionHorizontal>
    </div>

    <!--  <input type="button" class="btn btn-sm btn-outline-dark" value="スケジュールをスライドする" @click="slideSchedule(oldRows)">-->
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
        color-scheme="creamy"

        class="gantt-chart-body"
        :highlighted-dates="holidays"
        sticky
    >
      <g-gantt-row v-for="bar in bars" :key="bar.ganttBarConfig.id" :bars="[bar]"/>
      <template #side-menu>
        <table class="side-menu">
          <thead class="side-menu-header">
          <tr>
            <th class="side-menu-cell"></th><!-- css hack min-height -->
            <th v-for="item in GanttHeader" :key="item" class="side-menu-cell" :class="{'d-none': !item.visible}">
              {{ item.name }}
            </th>
          </tr>
          </thead>
          <tbody>
          <template v-for="item in ganttChartGroup" :key="item.ganttGroup.id">
            <tr v-for="row in item.rows" :key="row.ticket.id">
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="GanttHeader[0].visible">{{ unitMap[item.ganttGroup.unit_id] }}</gantt-td>
              <gantt-td :visible="GanttHeader[1].visible">
                <select v-model="row.ticket.process_id" @change="updateTicket(row.ticket)">
                  <option v-for="item in processList" :key="item.id" :value="item.id">{{ item.name }}</option>
                </select>
              </gantt-td>
              <gantt-td :visible="GanttHeader[2].visible">
                <select v-model="row.ticket.department_id" @change="updateTicket(row.ticket)">
                  <option v-for="item in departmentList" :key="item.id" :value="item.id">{{ item.name }}</option>
                </select>
              </gantt-td>
              <gantt-td :visible="GanttHeader[3].visible" style="width: 14rem;">
                <UserMultiselect :userList="userList" :ticketUser="row.ticketUsers"
                                 @update="ticketUserUpdate(row.ticket.id ,$event)"></UserMultiselect>
              </gantt-td>
              <gantt-td :visible="GanttHeader[4].visible">
                <input type="date" v-model="row.ticket.limit_date" @change="updateTicket(row.ticket)"/>
              </gantt-td>
              <gantt-td :visible="GanttHeader[5].visible">
                <FormNumber class="small-numeric" v-model="row.ticket.estimate" @change="updateTicket(row.ticket)"/>
              </gantt-td>
              <gantt-td :visible="GanttHeader[6].visible">
                <FormNumber class="small-numeric" v-model="row.ticket.days_after" @change="updateTicket(row.ticket)"/>
              </gantt-td>
              <gantt-td :visible="GanttHeader[7].visible">
                <input type="date" v-model="row.ticket.start_date" @change="updateTicket(row.ticket)"/>
              </gantt-td>
              <gantt-td :visible="GanttHeader[8].visible">
                <input type="date" v-model="row.ticket.end_date" @change="updateTicket(row.ticket)"/>
              </gantt-td>
              <gantt-td :visible="GanttHeader[9].visible">
                <FormNumber class="middle-numeric" v-model="row.ticket.progress_percent"
                            @change="updateTicket(row.ticket)"/>
              </gantt-td>
            </tr>
            <tr>
              <td colspan="11">
                <button @click="addNewTicket(item.ganttGroup.id)" class="btn btn-outline-primary">{{
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
  </div>
  <div v-else>
    ユニットを追加してください。
  </div>
</template>
<style lang="scss" scoped>

.justify-middle {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.filter {
  label {
    display: inline-block;
    position: relative;
    padding-left: 1em;
  }

  label input {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    margin: auto;
  }
}


.action-menu {
  height: 32px;
  > div {
    border-right: 1px solid #eaeaea;
    padding: 0 5px;
  }
}

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

.gantt-chart-body {
  height: 100%;
  margin-top: -30px;
  padding-top: 30px;
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

  .side-menu-header {
    box-shadow: 0px 8px 4px -5px rgba(50, 50, 50, 0.5);
    background: rgb(255, 246, 240);
    color: rgb(84, 45, 5);
  }

  > tbody > tr > td:first-child {
    border: none;
  }

  .side-menu-header > tr > th:first-child {
    border: none;
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
import AccordionHorizontal from "@/components/accordionHorizontal/AccordionHorizontal.vue";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {useFacility} from "@/composable/facility";

type GanttProxyProps = {
  facilityId: number
}

const props = defineProps<GanttProxyProps>()

const {
  chartStart,
  chartEnd,
  format,
  footerLabels,
  ganttChartGroup,
  bars,
  GanttHeader,
  slideSchedule,
  setScheduleByPersonDay,
  setScheduleByFromTo,
  adjustBar,
  addNewTicket,
  updateTicket,
  ticketUserUpdate,
} = await useGantt(props.facilityId)

const {list: unitList, refresh: unitRefresh} = await useUnitTable(props.facilityId)
const {list: processList, refresh: processRefresh} = await useProcessTable()
const {list: departmentList, refresh: departmentRefresh} = await useDepartmentTable()
const {list: userList, refresh: userRefresh} = await useUserTable()
const {list: holidayList, refresh: holidayRefresh} = await useHolidayTable(props.facilityId)
const {list: operationSettingList, refresh: operationSettingRefresh} = await useOperationSettingTable(props.facilityId)

const holidays = holidayList.value.map(v => new Date(v.date))

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
