<template>
  <div v-if="getOperationList().length > 0" id="gantt-proxy-wrapper">
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
          <input type="button" class="btn btn-sm btn-outline-dark" value="リスケ（工数h）重視"
                 @click="setScheduleByPersonDayProxy()">
          <input type="button" class="btn btn-sm btn-outline-dark" value="リスケ(日付)重視"
                 @click="setScheduleByFromToProxy()">
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
          <input class="form-check-input" type="radio" name="displayType" id="byWeek" v-model="displayType" value="week">
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

            :highlighted-dates="getHolidaysForGantt"
            sticky
        >
          <template #side-menu>
            <table class="side-menu" ref="ganttSideMenuElement">
              <thead class="side-menu-header">
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th :colspan="GanttHeader.length">
                  「{{ facility.name }}」のスケジュール：{{facility.term_from}}～{{facility.term_to}}
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
                    <UserMultiselect :userList="getUserList(row.ticket.department_id)" :ticketUser="row.ticketUsers"
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
        <g-gantt-chart
            :chart-start="chartStart"
            :chart-end="chartEnd"
            :precision="displayType"
            :row-height="40"
            grid
            :width="getGanttChartWidth"
            bar-start="beginDate"
            bar-end="endDate"
            :date-format="format"
            color-scheme="creamy"
            :hide-timeaxis="true"

            :highlighted-dates="getHolidaysForGantt"
            sticky
        >
          <template #side-menu>
            <table class="side-menu" :style="syncWidth">
              <tbody>
              <tr v-for="item in pileUpsByPerson" :key="item.user.id">
                <td class="side-menu-cell"></td><!-- css hack min-height -->
                <gantt-td :visible="true">{{ item.user.name }}</gantt-td>
              </tr>

              <tr v-for="item in pileUpsByDepartment" :key="item.departmentId">
                <td class="side-menu-cell"></td><!-- css hack min-height -->
                <gantt-td :visible="true">{{ getDepartmentName(item.departmentId) }}</gantt-td>
              </tr>
              </tbody>
            </table>
          </template>
          <g-gantt-label-row v-for="item in pileUpsByPerson" :key="item.user.id"
                             :labels="item.labels.map(v => v === 0 ? '' : round(v).toString())"></g-gantt-label-row>
          <g-gantt-label-row v-for="item in pileUpsByDepartment" :key="item.departmentId"
                             :labels="item.users.map(v => v.length === 0 ? '' : v.length.toString())"></g-gantt-label-row>
        </g-gantt-chart>
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

.gantt-wrapper {
  width: 100%;
  overflow-x: scroll;
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

  .side-menu-header > tr {
    height: 37.5px;
  }
  .side-menu-header > tr > th:first-child {
    border: none;
    display: block;
    float: left;
    content: "";
    min-height: 37.5px;
    height: 4vh;
  }

  .small-numeric {
    width: 3rem;
  }

  .middle-numeric {
    width: 4rem;
  }
}
.hide-scroll {
  scrollbar-width: none; /*Firefox対応のスクロールバー非表示コード*/
  -ms-overflow-style: none;/*Internet Explore対応のスクロールバー非表示コード*/
}
.hide-scroll::-webkit-scrollbar {
  display: none; /*Google Chrome、Safari、Microsoft Edge対応のスクロールバー非表示コード*/
}

</style>

<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttLabelRow, GGanttRow} from "@infectoone/vue-ganttastic";
import {useGantt} from "@/composable/gantt";
import FormNumber from "@/components/form/FormNumber.vue";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import AccordionHorizontal from "@/components/accordionHorizontal/AccordionHorizontal.vue";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {inject, nextTick, onMounted, onUnmounted, ref, watch} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {round} from "@/utils/math";

type GanttProxyProps = {
  facilityId: number
}

const props = defineProps<GanttProxyProps>()
const {processList, departmentList} = inject(GLOBAL_STATE_KEY)!

const {
  chartStart,
  chartEnd,
  format,
  ganttChartGroup,
  bars,
  GanttHeader,
  pileUpsByPerson,
  pileUpsByDepartment,
  getUnitName,
  getDepartmentName,
  facility,
  getGanttChartWidth,
  displayType,
  setScheduleByPersonDay,
  setScheduleByFromTo,
  adjustBar,
  addNewTicket,
  updateTicket,
  ticketUserUpdate,
  updateOrder,
  getOperationList,
  getHolidaysForGantt,
  updateDepartment,
  getUserList,
  deleteTicket,

} = await useGantt()

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
const syncWidth = ref<CSSStyleDeclaration>()
const ganttWrapperElement = ref<HTMLDivElement>()
const childGanttWrapperElement = ref<HTMLDivElement>()
const resizeSyncWidth = () => {
  const parentWidth = ganttSideMenuElement.value?.clientWidth
  syncWidth.value = {width: parentWidth + "px", overflow: 'scroll'}
  console.log("resizeSyncWidth", parentWidth)
}
watch(GanttHeader, () => {
  nextTick(resizeSyncWidth)
}, {deep: true})

// 一旦積み上げの更新を全体ウォッチにしておく。 パフォーマンス的には何かしら考慮が必要
// watch(ganttChartGroup, async () => {
//   await refreshPileUps()
// },{deep: true})

onMounted(() => {
  resizeSyncWidth()
  nextTick(resizeSyncWidth) // たまに上手くいかないので念のため
  ganttWrapperElement.value?.addEventListener("scroll", (event) => {
    childGanttWrapperElement.value?.scrollTo(event.srcElement.scrollLeft, 0)
  })
  childGanttWrapperElement.value?.addEventListener("scroll", (event) => {
    ganttWrapperElement.value?.scrollTo(event.srcElement.scrollLeft, 0)
  })
})
onUnmounted(() => {
  ganttWrapperElement.value?.removeEventListener("scroll")
  childGanttWrapperElement.value?.removeEventListener("scroll")
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
