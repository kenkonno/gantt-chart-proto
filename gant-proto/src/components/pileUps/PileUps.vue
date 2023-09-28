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
  >
    <template #side-menu>
      <table class="side-menu" :style="syncWidth">
        <tbody>
        <gantt-nested-row v-for="item in pileUpFilters" :key="item.departmentId">
          <template #parent>
            <tr>
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="true" class="justify-middle">
                      <span v-if="pileUpsByDepartment.find(v => v.departmentId === item.departmentId).hasError"
                            class="error-over-work-hour">稼働上限を超えている担当者がいます。</span>
                <span>{{ getDepartmentName(item.departmentId) }}(人)</span>
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
                <span v-if="user.hasError" class="error-over-work-hour">稼働上限を超えている日があります。</span>
                <span>{{ user.user.name }}(h)</span>
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
import {DisplayType, GanttChartGroup, GanttRow} from "@/composable/ganttFacility";
import GanttTd from "@/components/gantt/GanttTd.vue";
import {computed, ComputedRef, inject, StyleValue, watch} from "vue";
import GanttNestedRow from "@/components/gantt/GanttNestedRow.vue";
import {usePielUps} from "@/composable/pileUps";
import {DAYJS_FORMAT} from "@/utils/day";
import {Holiday, Ticket, TicketUser} from "@/api";
import {GLOBAL_GETTER_KEY} from "@/composable/globalState";
import {Api} from "@/api/axios";

type PileUpsProps = {
  tickets: Ticket[],
  ticketUsers: TicketUser[],
  chartStart: string,
  chartEnd: string,
  displayType: DisplayType,
  holidays: Holiday[],
  width: string,
  highlightedDates: Date[],
  syncWidth: CSSStyleDeclaration | undefined,
}
const props = defineProps<PileUpsProps>()
const {getDepartmentName} = inject(GLOBAL_GETTER_KEY)!
// FIXME: pileUpsに渡すときに、ほかの設備の奴らも渡してあげれば全体積み上げにできる。
// TODO: と思ったけど、やってみたら祝日問題とか色々有るので、全部の設備分計算してデフォルトのラベルを作ってあげるのがよさそう。結局全体ビューではやることになるので内部実装とする。祝日の概念がグローバルになれば処理としては簡単。
// TODO: 仮にやるとしたら。start, endはfacilityの全体。ticket,holiday,ticketUsersは各facility毎、でいい感じにindex操作してデフォルトの積み上げを用意しておくこと。
// TODO: API通信は増えるけど、一旦パフォーマンス無視して全部通信させてフロントだけで完結させよう。
// 設備一覧の取得
// 全チケット情報を取得するAPIの追加
// 全担当者情報を取得するAPIの追加
const {data: allTickets} = await Api.getAllTickets()
const {data: allTicketUsers} = await Api.getTicketUsers(allTickets.list.map(v => v.id!))

const tickets = computed(() => {
  const currentTicketIds = props.tickets.map(v => v.id)
  return props.tickets.concat(allTickets.list.filter(allTicket => !currentTicketIds.includes(allTicket.gantt_group_id)))

})
const ticketUsers = computed(() => {
  const currentTicketIds = props.tickets.map(v => v.id)
  return props.ticketUsers.concat(allTicketUsers.list.filter(allTicketUser => !currentTicketIds.includes(allTicketUser.ticket_id)))

})
const displayType = computed(() => props.displayType)
const holidays = computed(() => props.holidays)
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
    holidays
)

</script>
