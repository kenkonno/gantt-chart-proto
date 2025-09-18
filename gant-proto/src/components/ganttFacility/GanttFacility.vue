<template>
  <div v-if="getOperationList.length > 0" id="gantt-proxy-wrapper">
    <div class="gantt-wrapper d-flex flex-column">
      <div class="gantt-facility-wrapper d-flex overflow-x-scroll" ref="ganttWrapperElement"
           :class="{'hide-scroll': allowed('VIEW_PILEUPS') && globalState.showPileUp, 'full-max-height': !allowed('VIEW_PILEUPS') || !globalState.showPileUp}">
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
            @bar-update="onBarUpdate($event.bar, $event.newValue)"
            @contextmenu-bar="onContextmenuBar($event.bar, $event.e, $event.datetime)"
            color-scheme="creamy"

            :highlighted-dates="getHolidaysForGantt(displayType)"
            sticky
            :display-today-line="true"
            @today-line-position-x="initScroll($event, ganttWrapperElement)"
            :mile-stone-list="milestones"
            :vertical-lines="milestoneVerticalLines"
            ref="gGanttChartRef"
        >
          <template #side-menu>
            <table class="side-menu" ref="ganttSideMenuElement">
              <thead class="side-menu-header">
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th :colspan="props.ganttFacilityHeader.length">
                  <span>「</span>
                  <span class="d-inline-block position-relative">{{ facility.name }}<green-check
                      v-if="facility.type === FacilityType.Ordered"/></span>
                  <span>」</span>
                  <span>のスケジュール：{{ facility.term_from }}～{{ facility.term_to }}</span>
                </th>
              </tr>
              <tr>
                <th class="side-menu-cell"></th><!-- css hack min-height -->
                <th v-for="item in props.ganttFacilityHeader" :key="item" class="side-menu-cell"
                    :class="{'d-none': !item.visible}">
                  <unit-toggle-button :is-open="isAllOpenUnit" @toggle="toggleAllUnitOpenProxy()"
                                      v-if="item.name == 'ユニット'"
                  />
                  <span class="align-middle">{{ item.name }}</span>
                </th>
              </tr>
              </thead>
              <tbody>
              <template v-for="item in ganttChartGroup" :key="item.ganttGroup.id">
                <template v-if="isOpenUnit(item.unitId)">
                  <tr v-for="(row, index) in item.rows" :key="row.ticket.id">
                    <td class="side-menu-cell"></td><!-- css hack min-height -->
                    <gantt-td :visible="props.ganttFacilityHeader[0].visible" class="text-start">
                      <template v-if="index === 0">
                        <unit-toggle-button :is-open="isOpenUnit(item.unitId)"
                                            @toggle="toggleUnitOpenProxy(item.unitId)"/>
                        <span class="align-middle">{{ getUnitName(item.ganttGroup.unit_id) }}</span>
                      </template>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[1].visible">
                      <select :value="row.ticket.process_id"
                              @change="mutation.setProcessId($event.target.value, row.ticket)"
                              :disabled="!allowed('UPDATE_TICKET')">
                        <option v-for="item in processList" :key="item.id" :value="item.id">{{ item.name }}</option>
                      </select>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[2].visible">
                      <select :value="row.ticket.department_id"
                              @change="mutation.setDepartmentId($event.target.value, row.ticket)"
                              :disabled="!allowed('UPDATE_TICKET')">
                        <option v-for="item in cDepartmentList" :key="item.id" :value="item.id">{{ item.name }}</option>
                      </select>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[3].visible" style="min-width: 8rem;">
                      <UserMultiselect
                          :userList="getUserListByDepartmentId(row.ticket.department_id, row.ticket.start_date, row.ticket.end_date)"
                          :ticketUser="row.ticketUsers"
                          :disabled="!allowed('UPDATE_TICKET')"
                          @update:modelValue="mutation.setTicketUser(row.ticket ,$event)"></UserMultiselect>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[4].visible">
                      <FormNumber class="small-numeric"
                                  :value="row.ticket.number_of_worker"
                                  :min="1"
                                  @change="mutation.setNumberOfWorker($event, row.ticket)"
                                  :disabled="row.ticketUsers?.length > 0 || !allowed('UPDATE_TICKET')"/>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[5].visible">
                      <FormNumber class="small-numeric"
                                  :value="row.ticket.estimate"
                                  @change="mutation.setEstimate($event, row.ticket)"
                                  :disabled="!allowed('UPDATE_TICKET')"/>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[6].visible">
                      <FormNumber class="small-numeric"
                                  :value="row.ticket.days_after"
                                  @change="mutation.setDaysAfter($event, row.ticket)"
                                  :disabled="!allowed('UPDATE_TICKET')"/>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[7].visible">
                      <input type="date"
                             :value="row.ticket.start_date"
                             @change="mutation.setStartDate($event.target.value, row.ticket)"
                             :disabled="!allowed('UPDATE_TICKET')"/>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[8].visible">
                      <input type="date"
                             :value="row.ticket.end_date"
                             @change="mutation.setEndDate($event.target.value, row.ticket)"
                             :disabled="!allowed('UPDATE_TICKET')"/>
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[9].visible">
                    <FormNumber class="middle-numeric"
                                  :value="row.ticket.progress_percent"
                                  @change="mutation.setProgressPercent($event, row.ticket)"
                                  :disabled="!allowed('UPDATE_PROGRESS')"
                                  :min=0 />
                    </gantt-td>
                    <gantt-td :visible="props.ganttFacilityHeader[10].visible" v-if="allowed('UPDATE_TICKET')">
                      <a href="#" @click.prevent="!isUpdateOrder && updateOrder(item.rows, index, -1)"
                         :class="{disabled: isUpdateOrder}"><span
                          class="material-symbols-outlined">arrow_upward</span></a>
                      <a href="#" @click.prevent="!isUpdateOrder && updateOrder(item.rows, index,1)"
                         :class="{disabled: isUpdateOrder}"><span
                          class="material-symbols-outlined">arrow_downward</span></a>
                      <a href="#" @click.prevent="deleteTicket(row.ticket)"><span
                          class="material-symbols-outlined">delete</span></a>
                    </gantt-td>
                  </tr>
                </template>
                <template v-else>
                  <GanttSideMenuByUnit
                      :unit-name="getUnitName(item.ganttGroup.unit_id)"
                      :user-list="globalState.userList"
                      :gantt-chart-group="item"
                      :gantt-facility-header="ganttFacilityHeader"
                      :unit-id="item.unitId"
                      @toggle-unit="toggleUnitOpenProxy($event)"
                  />
                </template>
                <tr :style="addTicketRowStyle()" v-if="!hasFilter && isOpenUnit(item.unitId)">
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
          <g-gantt-row v-for="(bar, index) in bars" :key="index" :bars="bar"
                       :class="getUnitCollapseClass(bar)"/>
        </g-gantt-chart>
      </div>
      <!-- 山積み部分 -->
      <hr v-if="globalState.pileUpsRefresh && allowed('VIEW_PILEUPS') && globalState.showPileUp"/>
      <div class="gantt-facility-pile-ups-wrapper d-flex overflow-x-scroll flex-grow-1" ref="childGanttWrapperElement">
        <PileUps
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
            :vertical-lines="milestoneVerticalLines"
            :milestone-vertical-lines="[]"
            @on-mounted="forceScroll"
            :defaultPileUps="defaultPileUps"
            :global-start-date="globalStartDate"
            :default-valid-user-index-map="defaultValidUserIndexMap"
            v-if="globalState.pileUpsRefresh && allowed('VIEW_PILEUPS') && globalState.showPileUp"
            :force-reload="forcePileUpReload"
        >
        </PileUps>
      </div>
    </div>
  </div>
  <div v-else>
    ユニットを追加してください。
  </div>
  <Suspense v-if="modalIsOpen">
    <default-modal title="工程詳細" @close-edit-modal="closeModalProxy" :full-height="true">
      <async-ticket-edit :id="modalTicketId" :unit-id="modalUnitId" :facility-id="currentFacilityId"
                         @close-edit-modal="closeTicketMemo"></async-ticket-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>
<style>
@import '@/assets/gantt-override.scss';
</style>
<style lang="scss" scoped>
@use '@/assets/gantt.scss';
</style>

<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {useGanttFacility} from "@/composable/ganttFacility";
import FormNumber from "@/components/form/FormNumber.vue";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {computed, inject, nextTick, ref, watch} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {DAYJS_FORMAT} from "@/utils/day";
import PileUps from "@/components/pileUps/PileUps.vue";
import {useSyncScrollY, useSyncWidthAndScroll} from "@/composable/syncWidth";
import {initScroll} from "@/utils/initScroll";
import {DisplayType, GanttFacilityHeader} from "@/composable/ganttFacilityMenu";
import {allowed} from "@/composable/role";
import {Department, Ticket} from "@/api";
import {useModalWithId} from "@/composable/modalWIthId";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import AsyncTicketEdit from "@/components/ticket/AsyncTicketEdit.vue";
import dayjs from "dayjs";
import GreenCheck from "@/components/icon/GreenCheck.vue";
import {FacilityType} from "@/const/common";
import {getDefaultPileUps} from "@/composable/pileUps";
import {postTicketMemoById} from "@/composable/ticket";
import UnitToggleButton from "@/components/ganttFacility/UnitToggleButton.vue";
import GanttSideMenuByUnit from "@/components/ganttFacility/GanttSideMenuByUnit.vue";

type GanttProxyProps = {
  ganttFacilityHeader: GanttFacilityHeader[],
  displayType: DisplayType,
  ticketDailyWeightMode: boolean,
}
const props = defineProps<GanttProxyProps>()
const globalState = inject(GLOBAL_STATE_KEY)!
const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {processList, departmentList, facilityTypes} = inject(GLOBAL_STATE_KEY)!
const cDepartmentList = computed(() => {
  const result: Department[] = [{created_at: "", id: undefined, name: "", order: 0, updated_at: 0}]
  result.push(...departmentList)
  return result
})

const forcePileUpReload = ref(false)

const ret = await Promise.all([
  getDefaultPileUps(currentFacilityId, "day", false, facilityTypes),
  useGanttFacility(computed(() => props.ticketDailyWeightMode))
])

const {
  globalStartDate,
  defaultPileUps,
  defaultValidUserIndexMap,
} = ret[0]

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
  updateOrder,
  refreshTicketMemo,
  hasFilter,
  milestones,
  mutation,
  updateTicket,
  getUnitIdByTicketId,
  isUpdateOrder,
  isOpenUnit,
  isAllOpenUnit,
  toggleUnitOpen,
  toggleAllUnitOpen,
  getUnitCollapseClass,
  updateTicketDailyWeight,
} = ret[1]


const milestoneVerticalLines = computed(() => {
  return milestones.value.map(v => {
    let color = "blue"
    if (v.description == "出荷期日") {
      color = "yellow"
    }
    return {
      color: color,
      date: dayjs(v.date)
    }
  })
})

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
const gGanttChartRef = ref<HTMLDivElement>() // ガントチャート本体
const ganttSideMenuElement = ref<HTMLDivElement>() // サイドメニュー
const ganttWrapperElement = ref<HTMLDivElement>() // ガントチャートのラッパー
const childGanttWrapperElement = ref<HTMLDivElement>() // 積み上げ

const {
  syncWidth,
  resizeSyncWidth,
  forceScroll
} = useSyncWidthAndScroll(ganttSideMenuElement, ganttWrapperElement, childGanttWrapperElement, computed(() => globalState.showPileUp))// NOTE: GanttAll.vue二も変更が必要

useSyncScrollY(gGanttChartRef, gGanttChartRef)

watch(props.ganttFacilityHeader, () => {
  nextTick(resizeSyncWidth)
}, {deep: true})

// チケット更新モーダル
const {modalIsOpen, id: modalTicketId, closeEditModal} = useModalWithId()
const modalUnitId = ref(0)
const emit = defineEmits(["update"])

const closeModalProxy = async () => {
  emit("update")
  closeEditModal()
}
const closeTicketMemo = async (ticket: Ticket, userIds: number[]) => {
  // NOTE: 本当はAsyncEditMemo側で完結させたかったが、ガントチャートに反映させる都合上更新はこっちの方で実行するようにする。
  try {
    const result = await postTicketMemoById(ticket.id, ticket.memo, ticket.updated_at)
    ticket.updated_at = result?.updated_at
    await updateTicket(ticket)
    await mutation.setTicketUser(ticket, userIds)
    await refreshTicketMemo(ticket.id)
    closeEditModal()
  } catch (e) {
    console.log(e)
  }
}


const openTicketDetail = (ticketId: number, unitId: number) => {
  // NOTE: simulationモードの本チャンのチケットは文字列を含むのでNaNになる。コードとしてはよろしくない
  if (Number.isNaN(ticketId)) return
  if (!allowed('UPDATE_PROGRESS')) return
  modalTicketId.value = ticketId
  modalUnitId.value = unitId
  modalIsOpen.value = true
}

const addTicketRowStyle = () => {
  // 更新できる人の場合はフィルタを考慮する。
  if (allowed('UPDATE_TICKET')) {
    return {visibility: !hasFilter.value ? 'visible' : 'hidden'}
  } else {
    // 更新できない人はそもそも非表示とする。
    return {display: 'none'}
  }
}

const toggleAllUnitOpenProxy = () => {
  toggleAllUnitOpen()
  nextTick(resizeSyncWidth)
}

const toggleUnitOpenProxy = async (unitId: number) => {
  console.log("toggle-unit", unitId)
  await toggleUnitOpen(unitId)
  await nextTick(resizeSyncWidth)
}

let isDragged = false;

// ここからイベントフック
const onClickBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("click-bar", bar, e, datetime)
  if (!bar.editable) {
    if (!isDragged) {
      openTicketDetail(Number(bar.ganttBarConfig.id), getUnitIdByTicketId(Number(bar.ganttBarConfig.id)))
    }
    isDragged = false
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
  isDragged = true
  console.log("drag-bar", bar, e)
}

const onDragendBar = async (
    bar: GanttBarObject,
    e: MouseEvent,
) => {
  await adjustBar(bar)
  console.log("dragend-bar", bar, e)
}

const onContextmenuBar = (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("contextmenu-bar", bar, e, datetime)
}

const onBarUpdate = async (bar: GanttBarObject, workHour: number | undefined) => {
  console.log("bar-update", bar, workHour)
  // TODO: 本当は良くないけど文字列でデータを持たせる
  console.log("bar-update", bar.ganttBarConfig)
  const [ticketId, date] = bar.ganttBarConfig.id.split("@")
  await updateTicketDailyWeight(+ticketId, date, workHour)
  forcePileUpReload.value = !forcePileUpReload.value
}

</script>
