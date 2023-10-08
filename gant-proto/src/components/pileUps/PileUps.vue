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
  >
    <template #side-menu>
      <table class="side-menu" :style="syncWidth">
        <tbody>
        <gantt-nested-row v-for="item in pileUpFilters" :key="item.departmentId">
          <template #parent>
            <tr>
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="true" class="justify-middle">
                <tippy v-if="pileUpsByDepartment.find(v => v.departmentId === item.departmentId).hasError"
                       content="稼働上限を超えている担当者がいます。">
                  <span class="error-over-work-hour">{{ getDepartmentName(item.departmentId) }}(人)</span>
                </tippy>
                <span v-else>{{ getDepartmentName(item.departmentId) }}(人)</span>
                <span class="material-symbols-outlined pointer" v-if="!item.displayUsers"
                      @click="item.displayUsers = true">add</span>
                <span class="material-symbols-outlined pointer" v-else
                      @click="item.displayUsers = false">remove</span>
              </gantt-td>
            </tr>
          </template>
          <template #child v-if="item.displayUsers">
            <tr v-for="user in pileUpsByPerson.filter(v => v.user.department_id === item.departmentId)"
                :key="user.user.id">
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="true" class="justify-middle">
                <tippy v-if="user.hasError" content="稼働上限を超えている日があります。">
                  <span class="error-over-work-hour">{{ user.user.name }}(h)</span>
                </tippy>
                <span v-else>{{ user.user.name }}(h)</span>
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
import {computed, inject, onMounted, StyleValue, toValue} from "vue";
import GanttNestedRow from "@/components/gantt/GanttNestedRow.vue";
import {getDefaultPileUps, usePielUps} from "@/composable/pileUps";
import {DAYJS_FORMAT} from "@/utils/day";
import {Holiday, Ticket, TicketUser} from "@/api";
import {GLOBAL_GETTER_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {Tippy} from "vue-tippy";

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
}
const props = defineProps<PileUpsProps>()
const {currentFacilityId, departmentList, userList} = inject(GLOBAL_STATE_KEY)!
const {getDepartmentName} = inject(GLOBAL_GETTER_KEY)!

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
const {
  globalStartDate,
  defaultPileUpsByPerson,
  defaultPileUpsByDepartment
} = await getDefaultPileUps(currentFacilityId, "day")

console.log("####### start main usePileUps", globalStartDate)
const {
  pileUpFilters,
  pileUpsByDepartment,
  pileUpsByPerson,
  displayPileUps,
  refreshPileUps,
} = usePielUps(
    props.chartStart,
    props.chartEnd,
    tickets,
    ticketUsers,
    displayType,
    holidays,
    departmentList,
    userList,
    toValue(defaultPileUpsByPerson),
    toValue(defaultPileUpsByDepartment),
    globalStartDate
)
const emit = defineEmits(["onMounted"])
onMounted(() => {
  emit("onMounted")
})

</script>
