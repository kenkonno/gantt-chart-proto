<template>
  <g-gantt-chart
      :chart-start="chartStart"
      :chart-end="chartEnd"
      :precision="displayType"
      :row-height="40"
      grid
      :width="width"
      bar-start="beginDate"
      bar-end="endDate"
      :date-format="DAYJS_FORMAT"
      color-scheme="creamy"
      :hide-timeaxis="true"
      :highlighted-dates="highlightedDates"
      :vertical-lines="milestoneVerticalLines"
      sticky
      :display-today-line="true"
      ref="gGanttChartRef"
  >
    <template #side-menu>
      <table class="side-menu" :style="syncWidth">
        <tbody>
          <template v-for="pileUp in cPileUps" :key="pileUp.departmentId">
            <!-- 部署 -->
            <tr>
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="true">
                <div class="pileUp-title">
                  <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, '')" v-if="!pileUp.display">add</span>
                  <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, '')" v-if="pileUp.display">remove</span>
                  {{getDepartmentName(pileUp.departmentId)}}
                </div>
              </gantt-td>
            </tr>
            <template v-if="pileUp.display">
              <!-- アサイン済み -->
              <tr>
                <td class="side-menu-cell"></td><!-- css hack min-height -->
                <gantt-td :visible="true">
                  <div class="pileUp-title child">
                    <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, 'assignedUser')" v-if="!pileUp.assignedUser.display">add</span>
                    <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, 'assignedUser')" v-if="pileUp.assignedUser.display">remove</span>
                    アサイン済み
                  </div>
                </gantt-td>
              </tr>
              <!-- アサイン済みユーザー -->
              <template v-if="pileUp.assignedUser.display">
                <tr v-for="item in pileUp.assignedUser.users" :key="item.user.id">
                  <td class="side-menu-cell"></td><!-- css hack min-height -->
                  <gantt-td :visible="true">
                    <div class="pileUp-title child-child">
                      {{item.user.lastName}} {{item.user.firstName}}
                    </div>
                  </gantt-td>
                </tr>
              </template>
              <!-- 未アサイン -->
              <tr>
                <td class="side-menu-cell"></td><!-- css hack min-height -->
                <gantt-td :visible="true">
                  <div class="pileUp-title child">
                    <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, 'unAssignedPileUp')" v-if="!pileUp.unAssignedPileUp.display">add</span>
                    <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, 'unAssignedPileUp')" v-if="pileUp.unAssignedPileUp.display">remove</span>
                    未アサイン
                  </div>
                </gantt-td>
              </tr>
              <!-- 未アサイン - 案件 -->
              <template v-if="pileUp.unAssignedPileUp.display">
                <tr v-for="item in pileUp.unAssignedPileUp.facilities" :key="item.facilityId">
                  <td class="side-menu-cell"></td><!-- css hack min-height -->
                  <gantt-td :visible="true">
                    <div class="pileUp-title child-child">
                      {{getFacilityName(item.facilityId)}}
                    </div>
                  </gantt-td>
                </tr>
              </template>
              <!-- 未確定 -->
              <tr v-if="displayPrepared(pileUp.display)">
                <td class="side-menu-cell"></td><!-- css hack min-height -->
                <gantt-td :visible="true">
                  <div class="pileUp-title child">
                    <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, 'noOrdersReceivedPileUp')" v-if="!pileUp.noOrdersReceivedPileUp.display">add</span>
                    <span class="material-symbols-outlined pointer" @click="toggleDisplay(pileUp.departmentId, 'noOrdersReceivedPileUp')" v-if="pileUp.noOrdersReceivedPileUp.display">remove</span>
                    未確定
                  </div>
                </gantt-td>
              </tr>
              <!-- 未確定 - 案件 -->
              <template v-if="displayPrepared(pileUp.noOrdersReceivedPileUp.display)">
                <tr v-for="item in pileUp.noOrdersReceivedPileUp.facilities" :key="item.facilityId">
                  <td class="side-menu-cell"></td><!-- css hack min-height -->
                  <gantt-td :visible="true">
                    <div class="pileUp-title child-child">
                      {{getFacilityName(item.facilityId)}}
                    </div>
                  </gantt-td>
                </tr>
              </template>
            </template>
          </template>
        </tbody>
      </table>
    </template>
    <template v-for="pileUp in cPileUps" :key="pileUp.departmentId">
      <!-- 部署 -->
      <g-gantt-label-row :labels="pileUpsLabelFormat(pileUp.labels, displayType)" :styles="pileUp.styles"></g-gantt-label-row>
      <!-- アサイン済み -->
      <g-gantt-label-row :labels="pileUpsLabelFormat(pileUp.assignedUser.labels, displayType)" :styles="pileUp.assignedUser.styles" v-if="pileUp.display"></g-gantt-label-row>
      <!-- アサイン済みユーザー -->
      <template  v-if="pileUp.assignedUser.display && pileUp.display">
        <g-gantt-label-row v-for="item in pileUp.assignedUser.users" :key="item.user.id" :labels="pileUpsLabelFormat(item.labels, displayType)" :styles="item.styles"></g-gantt-label-row>
      </template>
      <!-- 未アサイン -->
      <g-gantt-label-row :labels="pileUpsLabelFormat(pileUp.unAssignedPileUp.labels, displayType)" :styles="pileUp.unAssignedPileUp.styles" v-if="pileUp.display"></g-gantt-label-row>
      <!-- 未アサイン - 案件 -->
      <template v-if="pileUp.unAssignedPileUp.display && pileUp.display">
        <g-gantt-label-row v-for="item in pileUp.unAssignedPileUp.facilities" :key="item.facilityId" :labels="pileUpsLabelFormat(item.labels, displayType)" :styles="item.styles"></g-gantt-label-row>
      </template>
      <!-- 未確定 -->
      <g-gantt-label-row :labels="pileUpsLabelFormat(pileUp.noOrdersReceivedPileUp.labels, displayType)" :styles="pileUp.noOrdersReceivedPileUp.styles"
                         v-if="displayPrepared(pileUp.display)"></g-gantt-label-row>
      <!-- 未確定 - 案件 -->
      <template v-if="displayPrepared(pileUp.noOrdersReceivedPileUp.display) && pileUp.display">
        <g-gantt-label-row v-for="item in pileUp.noOrdersReceivedPileUp.facilities" :key="item.facilityId" :labels="pileUpsLabelFormat(item.labels, displayType)" :styles="item.styles"></g-gantt-label-row>
      </template>
    </template>
  </g-gantt-chart>
</template>
<style lang="scss" scoped>
@import '@/assets/gantt.scss';

.pileUp-title {
  width: 50%;
  text-align: left;
  margin: auto;
}
.child {
  padding-left: 16px;
}
.child-child {
  padding-left: 32px;
}


</style>
<script setup lang="ts">
import {GGanttChart, GGanttLabelRow} from "@infectoone/vue-ganttastic";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {computed, ComputedRef, inject, onMounted, ref, StyleValue, toRefs, toValue, watch} from "vue";
import {getDefaultPileUps, usePileUps} from "@/composable/pileUps";
import {DAYJS_FORMAT} from "@/utils/day";
import {Holiday, Ticket, TicketUser} from "@/api";
import {GLOBAL_GETTER_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {Tippy} from "vue-tippy";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import {useSyncScrollY} from "@/composable/syncWidth";
import {VerticalLine} from "@infectoone/vue-ganttastic/lib_types/types";
import {pileUpsLabelFormat} from "@/utils/filters";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import {PileUps} from "@/composable/pileUps";
import {FacilityType} from "@/const/common";

type PileUpsProps = {
  tickets: ComputedRef<Ticket[]>,
  ticketUsers: TicketUser[],
  chartStart: string,
  chartEnd: string,
  displayType: DisplayType,
  holidays: Holiday[],
  width: string,
  highlightedDates: Date[],
  syncWidth: StyleValue | undefined,
  currentFacilityId: number,
  milestoneVerticalLines: VerticalLine[]
}
const props = defineProps<PileUpsProps>()
const {facilityList, departmentList, userList} = inject(GLOBAL_STATE_KEY)!
const globalState = inject(GLOBAL_STATE_KEY)!
const {getDepartmentName, getFacilityName} = inject(GLOBAL_GETTER_KEY)!
const {selectedDepartment, selectedUser} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!

const displayPrepared = (display: boolean) => {
  return display && globalState.facilityTypes.includes(FacilityType.Prepared)
}

// 結論
// usePileUpsを全ての案件で実行。
// 開始日、終了日は現在の案件。
// DefaultPileUpsByDepartment,DefaultPileUpsByPerson を受け付けるようにする
const isAllMode = props.currentFacilityId === -1
console.log("####### start main getDefaultPileUps")
const {
  globalStartDate,
  defaultPileUps,
} = await getDefaultPileUps(props.currentFacilityId, "day", isAllMode, globalState.facilityTypes)

console.log("####### start main getDefaultPileUps", defaultPileUps)

console.log("####### start main usePileUps", globalStartDate)

const {tickets, ticketUsers, holidays, displayType} = toRefs(props)

const {
  // pileUpFilters,
  // pileUpsByDepartment,
  // pileUpsByPerson,
  // displayPileUps,
    pileUps,
} = usePileUps(
    props.chartStart,
    props.chartEnd,
    isAllMode,
    facilityList.find(v => v.id === props.currentFacilityId)!,
    tickets,
    ticketUsers,
    displayType,
    holidays,
    departmentList,
    userList,
    facilityList,
    defaultPileUps,
    globalStartDate
)

const toggleDisplay = (departmentId: number, type: "" | "assignedUser" | "noOrdersReceivedPileUp" | "unAssignedPileUp") => {
  const target = pileUps.value.find(v => v.departmentId === departmentId)!
  if ( type === "" ) {
    target.display = !target.display
  } else {
    target[type].display = !target[type].display
  }
}

const cPileUps = computed(() => {
  let result: PileUps[] = pileUps.value
  // 部署のフィルタ
  if (selectedDepartment.value !== undefined) {
    result = result.filter(v => v.departmentId === selectedDepartment.value)
  }
  if (selectedUser.value !== undefined) {
    result = result.map(function(v): PileUps {
      // ユーザーに紐づく部署を特定する
      return {
        assignedUser: {display: v.assignedUser.display, labels: v.assignedUser.labels, styles: v.assignedUser.styles, users: v.assignedUser.users.filter(user => user.user.id === selectedUser.value)},
        noOrdersReceivedPileUp: {display: v.noOrdersReceivedPileUp.display, facilities: [], labels: [], styles: []},
        unAssignedPileUp: {display: v.unAssignedPileUp.display, facilities: [], labels: [], styles: []},
        departmentId: v.departmentId,
        display: v.display,
        labels: v.labels,
        styles: v.styles,
      }
    })
  }
  return result

})

const emit = defineEmits(["onMounted"])
onMounted(() => {
  emit("onMounted")
})

const updateDepartmentFilter = (departmentId: number) => {
  selectedUser.value = undefined
  selectedDepartment.value = departmentId
}
const updateUserFilter = (userId: number | undefined) => {
  selectedDepartment.value = undefined
  selectedUser.value = userId
}

const gGanttChartRef = ref<HTMLDivElement>() // ガントチャート本体
useSyncScrollY(gGanttChartRef, gGanttChartRef)


</script>
