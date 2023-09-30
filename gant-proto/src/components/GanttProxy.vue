<template>
  <div v-if="getOperationList.length > 0" id="gantt-proxy-wrapper">
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
          <tippy content="（開始日 又は 日後）・工数・工程・人数 が入力されている行が対象です。">
            <input type="button" class="btn btn-sm btn-outline-dark" value="リスケ（工数h）重視"
                   @click="setScheduleByPersonDayProxy()">
          </tippy>
          <tippy content="開始日・終了日・工数・工程が入力されている行が対象です。担当所が入力されている行は無視されます。">
            <input type="button" class="btn btn-sm btn-outline-dark" value="リスケ(日付)重視"
                   @click="setScheduleByFromToProxy()">
          </tippy>
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
      <div class="d-flex justify-middle">
        <div class="form-check">
          <input class="form-check-input" type="radio" name="displayType" id="byDay" v-model="displayType" value="day">
          <label class="form-check-label" for="byDay">
            日毎
          </label>
        </div>
        <div class="form-check">
          <input class="form-check-input" type="radio" name="displayType" id="byWeek" v-model="displayType"
                 value="week">
          <label class="form-check-label" for="byWeek">
            週次
          </label>
        </div>
      </div>
    </div>

    <!--  <input type="button" class="btn btn-sm btn-outline-dark" value="スケジュールをスライドする" @click="slideSchedule(oldRows)">-->
    <div class="gantt-wrapper">
      <div class="d-flex overflow-x-scroll hide-scroll" ref="ganttWrapperElement">
        <g-gantt-chart
            :chart-start="chartStart"
            :chart-end="chartEnd"
            :precision="displayType"
            :row-height="40"
            grid
            :width="getGanttChartWidth"
            bar-start="beginDate"
            bar-end="endDate"
            :date-format="DAYJS_FORMAT"
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

            :highlighted-dates="getHolidaysForGantt"
            sticky
        >
          <template #side-menu>
            <table class="side-menu" ref="ganttSideMenuElement">
              <thead class="side-menu-header">
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th :colspan="GanttHeader.length">
                  「{{ facility.name }}」のスケジュール：{{ facility.term_from }}～{{ facility.term_to }}
                </th>
              </tr>
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th v-for="item in GanttHeader" :key="item" class="side-menu-cell" :class="{'d-none': !item.visible}">
                  {{ item.name }}
                </th>
              </tr>
              </thead>
              <tbody>
              <template v-for="item in ganttChartGroup" :key="item.ganttGroup.id">
                <tr v-for="(row, index) in item.rows" :key="row.ticket.id">
                  <td class="side-menu-cell"></td><!-- css hack min-height -->
                  <gantt-td :visible="GanttHeader[0].visible">{{ getUnitName(item.ganttGroup.unit_id) }}</gantt-td>
                  <gantt-td :visible="GanttHeader[1].visible">
                    <select v-model="row.ticket.process_id" @change="updateTicket(row.ticket)">
                      <option v-for="item in processList" :key="item.id" :value="item.id">{{ item.name }}</option>
                    </select>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[2].visible">
                    <select v-model="row.ticket.department_id" @change="updateDepartment(row.ticket)">
                      <option v-for="item in departmentList" :key="item.id" :value="item.id">{{ item.name }}</option>
                    </select>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[3].visible" style="min-width: 8rem;">
                    <UserMultiselect :userList="getUserListByDepartmentId(row.ticket.department_id)"
                                     :ticketUser="row.ticketUsers"
                                     @update:modelValue="ticketUserUpdate(row.ticket ,$event)"></UserMultiselect>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[4].visible">
                    <FormNumber class="small-numeric" v-model="row.ticket.number_of_worker"
                                @change="updateTicket(row.ticket)" :disabled="row.ticketUsers?.length > 0"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[5].visible">
                    <input type="date" v-model="row.ticket.limit_date" @change="updateTicket(row.ticket)"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[6].visible">
                    <FormNumber class="small-numeric" v-model="row.ticket.estimate" @change="updateTicket(row.ticket)"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[7].visible">
                    <FormNumber class="small-numeric" v-model="row.ticket.days_after"
                                @change="updateTicket(row.ticket)"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[8].visible">
                    <input type="date" v-model="row.ticket.start_date" @change="updateTicket(row.ticket)"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[9].visible">
                    <input type="date" v-model="row.ticket.end_date" @change="updateTicket(row.ticket)"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[10].visible">
                    <FormNumber class="middle-numeric" v-model="row.ticket.progress_percent"
                                @change="updateTicket(row.ticket)"/>
                  </gantt-td>
                  <gantt-td :visible="GanttHeader[11].visible">
                    <a href="#" @click="updateOrder(item.rows,index, -1)"><span class="material-symbols-outlined">arrow_upward</span></a>
                    <a href="#" @click="updateOrder(item.rows, index,1)"><span class="material-symbols-outlined">arrow_downward</span></a>
                    <a href="#" @click="deleteTicket(row.ticket)"><span class="material-symbols-outlined">delete</span></a>
                  </gantt-td>
                </tr>
                <tr>
                  <td :colspan="GanttHeader.length + 1">
                    <button @click="addNewTicket(item.ganttGroup.id)" class="btn btn-outline-primary">{{
                        getUnitName(item.ganttGroup.unit_id)
                      }}の工程を追加する
                    </button>
                  </td>
                </tr>
              </template>
              </tbody>
            </table>
          </template>
          <g-gantt-row v-for="bar in bars" :key="bar.ganttBarConfig.id" :bars="[bar]"/>
        </g-gantt-chart>
      </div>
      <!-- 山積み部分 -->
      <hr>
      <div class="d-flex overflow-x-scroll" ref="childGanttWrapperElement">
        <PileUps :chart-start="chartStart"
                 :chart-end="chartEnd"
                 :display-type="displayType"
                 :holidays="getHolidays"
                 :tickets="getTickets"
                 :ticket-users="ticketUserList"
                 :width="getGanttChartWidth"
                 :highlightedDates="getHolidaysForGantt"
                 :syncWidth="syncWidth"

        >
        </PileUps>
      </div>
    </div>
  </div>
  <div v-else>
    ユニットを追加してください。
  </div>
</template>
<style>
.g-gantt-chart {
  flex-shrink: 0;
}
</style>
<style lang="scss" scoped>
@import '@/assets/gantt.scss';
</style>

<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {useGanttFacility} from "@/composable/ganttFacility";
import FormNumber from "@/components/form/FormNumber.vue";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import AccordionHorizontal from "@/components/accordionHorizontal/AccordionHorizontal.vue";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {inject, nextTick, ref, watch} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {DAYJS_FORMAT} from "@/utils/day";
import PileUps from "@/components/pileUps/PileUps.vue";
import {Tippy} from "vue-tippy";
import {useSyncWidthAndScroll} from "@/composable/syncWidth";

type GanttProxyProps = {
  facilityId: number
}

defineProps<GanttProxyProps>()
const {processList, departmentList} = inject(GLOBAL_STATE_KEY)!

const {
  GanttHeader,
  bars,
  chartEnd,
  chartStart,
  displayType,
  facility,
  getHolidaysForGantt,
  ganttChartGroup,
  getGanttChartWidth,
  getUnitName,
  getOperationList,
  getHolidays,
  getTickets,
  ticketUserList,
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
} = await useGanttFacility()

// リスケ関連
const setScheduleByPersonDayProxy = () => {
  ganttChartGroup.value.forEach(v => {
    setScheduleByPersonDay(v.rows)
  })
}

const setScheduleByFromToProxy = () => {
  ganttChartGroup.value.forEach(v => {
    setScheduleByFromTo(v.rows)
  })
}


// スクロールや大きさの同期系
const ganttSideMenuElement = ref<HTMLDivElement>()
const ganttWrapperElement = ref<HTMLDivElement>()
const childGanttWrapperElement = ref<HTMLDivElement>()

const {
  syncWidth,
  resizeSyncWidth
} = useSyncWidthAndScroll(ganttSideMenuElement, ganttWrapperElement, childGanttWrapperElement)

watch(GanttHeader, () => {
  nextTick(resizeSyncWidth)
}, {deep: true})

// ここからイベントフック
const onClickBar = async (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("click-bar", bar, e, datetime)
  await adjustBar(bar)
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
