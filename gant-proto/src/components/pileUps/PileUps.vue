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
      sticky
      :display-today-line="true"
      ref="gGanttChartRef"
  >
    <template #side-menu>
      <table class="side-menu" :style="syncWidth">
        <tbody>
        <gantt-nested-row v-for="item in cPileUpFilters" :key="item.departmentId">
          <template #parent>
            <tr>
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="true" class="justify-middle">
                <tippy v-if="pileUpsByDepartment.find(v => v.departmentId === item.departmentId).hasError"
                       content="稼働上限を超えている担当者がいます。">
                  <span class="error-over-work-hour" @click="updateDepartmentFilter(item.departmentId)">{{ getDepartmentName(item.departmentId) }}(人)</span>
                </tippy>
                <span v-else  @click="updateDepartmentFilter(item.departmentId)">{{ getDepartmentName(item.departmentId) }}(人)</span>
                <span class="material-symbols-outlined pointer" v-if="!item.displayUsers"
                      @click="item.displayUsers = true">add</span>
                <span class="material-symbols-outlined pointer" v-else
                      @click="item.displayUsers = false">remove</span>
              </gantt-td>
            </tr>
          </template>
          <template #child v-if="item.displayUsers">
            <tr v-for="user in cPileUpsPerson.filter(v => v.user.department_id === item.departmentId)"
                :key="user.user.id">
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="true" class="justify-middle">
                <tippy v-if="user.hasError" content="稼働上限を超えている日があります。">
                  <span class="error-over-work-hour" @click="updateUserFilter(user.user.id)">{{ user.user.name }}(h)</span>
                </tippy>
                <span v-else @click="updateUserFilter(user.user.id)">{{ user.user.name }}(h)</span>
              </gantt-td>
            </tr>
          </template>
        </gantt-nested-row>
        </tbody>
      </table>
    </template>
    <g-gantt-label-row v-for="(item, index) in displayPileUps" :key="index" :labels="item.labels"
                       :styles="item.styles"></g-gantt-label-row>
  </g-gantt-chart>
</template>
<style lang="scss" scoped>
@import '@/assets/gantt.scss';
</style>
<script setup lang="ts">
import {GGanttChart, GGanttLabelRow} from "@infectoone/vue-ganttastic";
import {DisplayType} from "@/composable/ganttFacility";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {computed, inject, onMounted, ref, StyleValue, toValue} from "vue";
import GanttNestedRow from "@/components/gantt/GanttNestedRow.vue";
import {getDefaultPileUps, usePielUps} from "@/composable/pileUps";
import {DAYJS_FORMAT} from "@/utils/day";
import {Holiday, Ticket, TicketUser} from "@/api";
import {GLOBAL_GETTER_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {Tippy} from "vue-tippy";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import {useSyncScrollY} from "@/composable/syncWidth";

type PileUpsProps = {
  tickets: Ticket[],
  ticketUsers: TicketUser[],
  chartStart: string,
  chartEnd: string,
  displayType: DisplayType,
  holidays: Holiday[],
  width: string,
  highlightedDates: Date[],
  syncWidth: StyleValue | undefined,
  currentFacilityId: number,
}
const props = defineProps<PileUpsProps>()
const {departmentList, userList} = inject(GLOBAL_STATE_KEY)!
const {getDepartmentName} = inject(GLOBAL_GETTER_KEY)!
const {selectedDepartment, selectedUser} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!

const tickets = computed(() => {
  return props.tickets
})
const ticketUsers = computed(() => {
  return props.ticketUsers
})
const displayType = computed(() => props.displayType)
const holidays = computed(() => props.holidays)

// 結論
// usePileUpsを全ての設備で実行。
// 開始日、終了日は現在の設備。
// DefaultPileUpsByDepartment,DefaultPileUpsByPerson を受け付けるようにする
console.log("####### start main getDefaultPileUps")
const {
  globalStartDate,
  defaultPileUpsByPerson,
  defaultPileUpsByDepartment,
} = await getDefaultPileUps(props.currentFacilityId, "day", selectedDepartment, selectedUser,)
console.log("####### start main getDefaultPileUps", defaultPileUpsByPerson, defaultPileUpsByDepartment)

console.log("####### start main usePileUps", globalStartDate)
const {
  pileUpFilters,
  pileUpsByDepartment,
  pileUpsByPerson,
  displayPileUps,
} = usePielUps(
    props.chartStart,
    props.chartEnd,
    tickets,
    ticketUsers,
    displayType,
    holidays,
    departmentList,
    userList,
    selectedDepartment,
    selectedUser,
    toValue(defaultPileUpsByPerson),
    toValue(defaultPileUpsByDepartment),
    globalStartDate
)

// TODO: フィルタ処理が分散してわかりにくい。
const cPileUpFilters = computed(() => {
  let targetDepartmentId = selectedDepartment.value
  if (selectedUser.value != undefined) {
    targetDepartmentId = userList.find(v => v.id == selectedUser.value)?.department_id
  } else {
    targetDepartmentId = selectedDepartment.value
  }
  if (targetDepartmentId == undefined) {
    return pileUpFilters.value
  } else {
    return pileUpFilters.value.filter(v => v.departmentId == targetDepartmentId)
  }
})

const cPileUpsPerson = computed(() => {
  if (selectedUser.value == undefined) {
    return pileUpsByPerson.value
  } else {
    return pileUpsByPerson.value.filter(v => v.user.id == selectedUser.value)
  }
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
