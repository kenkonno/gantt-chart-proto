<template>
  <div v-if="getOperationList.length > 0" id="gantt-proxy-wrapper">
    <div class="gantt-wrapper">
      <div class="d-flex overflow-x-scroll hide-scroll" ref="ganttWrapperElement">
        <g-gantt-chart
            :chart-start="chartStart"
            :chart-end="chartEnd"
            :precision="displayType"
            :row-height="40"
            grid
            :width="getGanttChartWidth(displayType)"
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

            :highlighted-dates="getHolidaysForGantt(displayType)"
            sticky
            :display-today-line="true"
            @today-line-position-x="initScroll($event, ganttWrapperElement)"
        >
          <template #side-menu>
            <table class="side-menu" ref="ganttSideMenuElement">
              <thead class="side-menu-header">
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th :colspan="props.ganttFacilityHeader.length">
                  「{{ facility.name }}」のスケジュール：{{ facility.term_from }}～{{ facility.term_to }}
                </th>
              </tr>
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th v-for="item in props.ganttFacilityHeader" :key="item" class="side-menu-cell"
                    :class="{'d-none': !item.visible}">
                  {{ item.name }}
                </th>
              </tr>
              </thead>
              <tbody>
              <template v-for="item in ganttChartGroup" :key="item.ganttGroup.id">
                <tr v-for="(row, index) in item.rows" :key="row.ticket.id">
                  <td class="side-menu-cell"></td><!-- css hack min-height -->
                  <gantt-td :visible="props.ganttFacilityHeader[0].visible">{{
                      getUnitName(item.ganttGroup.unit_id)
                    }}
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[1].visible">
                    <select v-model="row.ticket.process_id" @change="updateTicket(row.ticket)"
                            :disabled="!allowed('UPDATE_TICKET')">
                      <option v-for="item in processList" :key="item.id" :value="item.id">{{ item.name }}</option>
                    </select>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[2].visible">
                    <select v-model="row.ticket.department_id" @change="updateDepartment(row.ticket)"
                            :disabled="!allowed('UPDATE_TICKET')">
                      <option v-for="item in departmentList" :key="item.id" :value="item.id">{{ item.name }}</option>
                    </select>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[3].visible" style="min-width: 8rem;">
                    <UserMultiselect :userList="getUserListByDepartmentId(row.ticket.department_id)"
                                     :ticketUser="row.ticketUsers"
                                     :disabled="!allowed('UPDATE_TICKET')"
                                     @update:modelValue="ticketUserUpdate(row.ticket ,$event)"></UserMultiselect>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[4].visible">
                    <FormNumber class="small-numeric" v-model="row.ticket.number_of_worker"
                                @change="updateTicket(row.ticket)"
                                :disabled="row.ticketUsers?.length > 0 || !allowed('UPDATE_TICKET')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[5].visible">
                    <input type="date" v-model="row.ticket.limit_date" @change="updateTicket(row.ticket)"
                           :disabled="!allowed('UPDATE_TICKET')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[6].visible">
                    <FormNumber class="small-numeric" v-model="row.ticket.estimate" @change="updateTicket(row.ticket)"
                                :disabled="!allowed('UPDATE_TICKET')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[7].visible">
                    <FormNumber class="small-numeric" v-model="row.ticket.days_after"
                                @change="updateTicket(row.ticket)"
                                :disabled="!allowed('UPDATE_TICKET')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[8].visible">
                    <input type="date" v-model="row.ticket.start_date" @change="updateTicket(row.ticket)"
                           :disabled="!allowed('UPDATE_TICKET')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[9].visible">
                    <input type="date" v-model="row.ticket.end_date" @change="updateTicket(row.ticket)"
                           :disabled="!allowed('UPDATE_TICKET')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[10].visible">
                    <FormNumber class="middle-numeric" v-model="row.ticket.progress_percent"
                                @change="updateTicket(row.ticket)"
                                :disabled="!allowed('UPDATE_PROGRESS')"/>
                  </gantt-td>
                  <gantt-td :visible="props.ganttFacilityHeader[11].visible" v-if="allowed('UPDATE_TICKET')">
                    <a href="#" @click.prevent="updateOrder(item.rows, index, -1)"><span
                        class="material-symbols-outlined">arrow_upward</span></a>
                    <a href="#" @click.prevent="updateOrder(item.rows, index,1)"><span
                        class="material-symbols-outlined">arrow_downward</span></a>
                    <a href="#" @click.prevent="deleteTicket(row.ticket)"><span
                        class="material-symbols-outlined">delete</span></a>
                  </gantt-td>
                </tr>
                <tr v-if="!hasFilter && allowed('UPDATE_TICKET')">
                  <td :colspan="props.ganttFacilityHeader.length + 1">
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
        <PileUps
            v-if="globalState.pileUpsRefresh"
            :chart-start="chartStart"
            :chart-end="chartEnd"
            :display-type="displayType"
            :holidays="getHolidays"
            :tickets="getTickets"
            :ticket-users="ticketUserList"
            :width="getGanttChartWidth(displayType)"
            :highlightedDates="getHolidaysForGantt(displayType)"
            :syncWidth="syncWidth"
            :current-facility-id="currentFacilityId"
            @on-mounted="forceScroll"
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
@import '@/assets/gantt-override.scss';
</style>
<style lang="scss" scoped>
@import '@/assets/gantt';
</style>

<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {useGanttFacility} from "@/composable/ganttFacility";
import FormNumber from "@/components/form/FormNumber.vue";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {inject, nextTick, ref, watch} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {DAYJS_FORMAT} from "@/utils/day";
import PileUps from "@/components/pileUps/PileUps.vue";
import {useSyncWidthAndScroll} from "@/composable/syncWidth";
import {initScroll} from "@/utils/initScroll";
import {DisplayType, GanttFacilityHeader} from "@/composable/ganttFacilityMenu";
import {allowed} from "@/composable/role";

type GanttProxyProps = {
  ganttFacilityHeader: GanttFacilityHeader[],
  displayType: DisplayType,
}
const props = defineProps<GanttProxyProps>()
const globalState = inject(GLOBAL_STATE_KEY)!
const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {processList, departmentList} = inject(GLOBAL_STATE_KEY)!
const {
  bars,
  chartEnd,
  chartStart,
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
  hasFilter
} = await useGanttFacility()

// リスケ関連 親から呼び出される
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
defineExpose({setScheduleByPersonDayProxy, setScheduleByFromToProxy})
// スクロールや大きさの同期系
const ganttSideMenuElement = ref<HTMLDivElement>()
const ganttWrapperElement = ref<HTMLDivElement>()
const childGanttWrapperElement = ref<HTMLDivElement>()

const {
  syncWidth,
  resizeSyncWidth,
  forceScroll
} = useSyncWidthAndScroll(ganttSideMenuElement, ganttWrapperElement, childGanttWrapperElement)

watch(props.ganttFacilityHeader, () => {
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
