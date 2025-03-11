<template>
  <div class="gantt-wrapper" id="gantt-all-view" :class="{withFilter:hasFilter(), byProcess:byProcess()}">
    <div class="gantt-facility-wrapper d-flex overflow-x-scroll" ref="ganttWrapperElement"
         :class="{'hide-scroll': allowed('VIEW_PILEUPS')&& globalState.showPileUp, 'full-max-height': !allowed('VIEW_PILEUPS')  || !globalState.showPileUp}"
    >
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
          color-scheme="creamy"
          :highlighted-dates="holidaysAsDate(displayType)"
          sticky
          :display-today-line="true"
          @today-line-position-x="initScroll($event, ganttWrapperElement)"
          @click-bar="onClickBar($event.bar, $event.e, $event.datetime)"
          ref="gGanttChartRef"
      >
        <template #side-menu>
          <table class="side-menu" ref="ganttSideMenuElement">
            <thead class="side-menu-header">
            <tr>
              <th class="side-menu-cell"></th><!-- css hack min-height -->
              <th :colspan="ganttAllHeader.length">
                「全体」のスケジュール：{{ startDate }}～{{ endDate }}
              </th>
            </tr>
            <tr>
              <th class="side-menu-cell"></th><!-- css hack min-height -->
              <th v-for="item in ganttAllHeader" :key="item" class="side-menu-cell" :class="{'d-none': !item.visible}">
                {{ item.name }}
              </th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="item in ganttAllRow" :key="item.facility.id">
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="ganttAllHeader[0].visible" @click="refreshGantt(item.facility.id, true)"
                        class="pointer facility-name">
                <u class="d-inline-block position-relative">{{ item.facility.name }}<green-check v-if="item.facility.type === FacilityType.Ordered"/></u>
              </gantt-td>
              <gantt-td :visible="ganttAllHeader[1].visible">
                <div class="user-wrapper">
                  <SingleRune v-for="user in item.users" :key="user.id" :name="getRuneName(user.lastName, user.firstName)" :id="user.id"></SingleRune>
                </div>
              </gantt-td>
              <gantt-td :visible="ganttAllHeader[2].visible">{{ item.startDate }}</gantt-td>
              <gantt-td :visible="ganttAllHeader[3].visible">{{ item.endDate }}</gantt-td>
              <gantt-td :visible="ganttAllHeader[4].visible" class="estimate">{{ item.estimate }}</gantt-td>
              <gantt-td :visible="ganttAllHeader[5].visible" class="progress-percent">{{ $filters.progressFormat(item.progress_percent) }}</gantt-td>
            </tr>
            </tbody>
          </table>
        </template>
        <g-gantt-row v-for="item in ganttAllRow" :key="item.facility.id" :bars="item.bars"/>
      </g-gantt-chart>
    </div>
    <!-- 山積み部分 -->
    <hr v-if="allowed('VIEW_PILEUPS') && globalState.showPileUp" />
    <div class="gantt-facility-pile-ups-wrapper d-flex overflow-x-scroll" ref="childGanttWrapperElement">
      <PileUps
          :chart-start="chartStart"
          :chart-end="chartEnd"
          :display-type="displayType"
          :holidays="holidays"
          :tickets="tickets"
          :ticket-users="ticketUsers"
          :width="getGanttChartWidth(displayType)"
          :highlightedDates="holidaysAsDate(displayType)"
          :syncWidth="syncWidth"
          :current-facility-id="-1"
          :milestone-vertical-lines="[]"
          @on-mounted="forceScroll"
          :defaultPileUps="defaultPileUps"
          :global-start-date="globalStartDate"
          :default-valid-user-index-map="defaultValidUserIndexMap"
          v-if="allowed('VIEW_PILEUPS') && globalState.showPileUp"
      >
      </PileUps>
    </div>
  </div>
</template>
<style>
@import '@/assets/gantt-override.scss';
#gantt-all-view.byProcess .g-gantt-row-bars-container .g-gantt-bar {
  opacity: 0.8;
}

#gantt-all-view.withFilter .g-gantt-row-bars-container .g-gantt-bar {
  height: 33% !important;
}

#gantt-all-view.withFilter .g-gantt-row-bars-container .g-gantt-bar:nth-child(2n) {
  top: 33% !important;
  height: 33% !important;
}

#gantt-all-view.withFilter .g-gantt-row-bars-container .g-gantt-bar:nth-child(3n) {
  top: 66% !important;
  height: 33% !important;
}
</style>
<style lang="scss" scoped>
@import '@/assets/gantt.scss';

nav {
  padding: 10px;

  > div {
    width: 100%;
    text-align: left;

    > select {
      margin: 0 5px;
    }
  }

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}

.user-wrapper {
  display: flex;
  width: 100%;
  height: 100%;
  padding-left: 10px;

  > div {
    margin: auto 0 auto -10px;
  }
}
.facility-name {
  text-align: left;
}
.estimate {
  text-align: right;
}
.progress-percent {
  text-align: right;
}

</style>
<script setup lang="ts">
import {GanttBarObject, GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {DAYJS_FORMAT} from "@/utils/day";
import GanttTd from "@/components/gantt/GanttTd.vue";
import PileUps from "@/components/pileUps/PileUps.vue";
import {useGanttAll} from "@/composable/ganttAll";
import {useSyncScrollY, useSyncWidthAndScroll} from "@/composable/syncWidth";
import SingleRune from "@/components/form/SingleRune.vue";
import {initScroll} from "@/utils/initScroll";
import {inject, nextTick, ref, watch} from "vue";
import {Header} from "@/composable/ganttAllMenu";
import {GLOBAL_SCHEDULE_ALERT_KEY} from "@/composable/scheduleAlert";
import {GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {AggregationAxis, DisplayType} from "@/composable/ganttFacilityMenu";
import {allowed} from "@/composable/role";
import {FacilityType} from "@/const/common";
import GreenCheck from "@/components/icon/GreenCheck.vue";
import {Ticket, TicketUser} from "@/api";
import {getDefaultPileUps} from "@/composable/pileUps";

const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!
const {facilityTypes} = inject(GLOBAL_STATE_KEY)!
const globalState = inject(GLOBAL_STATE_KEY)!

type GanttAllProps = {
  ganttAllHeader: Header[],
  displayType: DisplayType,
  aggregationAxis: AggregationAxis,
}
const props = defineProps<GanttAllProps>()

// パフォーマンス対応のためにガントチャートの時点で並行してAPIコールするように修正。
const ret = await Promise.all([getDefaultPileUps(-1, "day", true, facilityTypes), useGanttAll(props.aggregationAxis)])
const {
  globalStartDate,
  defaultPileUps,
  defaultValidUserIndexMap,
} =  ret[0]

const {
  startDate,
  endDate,
  ganttAllRow,
  holidaysAsDate,
  getGanttChartWidth,
  holidays,
  chartStart,
  chartEnd,
  hasFilter
} = ret[1]

// NOTE: PileUpsのPropsだが、明示的に変数じゃないとwatchが多重で実行されてしまう。
const tickets: Ticket[] = []
const ticketUsers: TicketUser[] = []

const byProcess = () => {
  return props.aggregationAxis == "process"
}

// スクロールや大きさの同期系
const gGanttChartRef = ref<HTMLDivElement>() // ガントチャート本体
const ganttSideMenuElement = ref<HTMLDivElement>()
const ganttWrapperElement = ref<HTMLDivElement>()
const childGanttWrapperElement = ref<HTMLDivElement>()

const {
  syncWidth,
  resizeSyncWidth,
  forceScroll
} = useSyncWidthAndScroll(ganttSideMenuElement, ganttWrapperElement, childGanttWrapperElement)
watch(props.ganttAllHeader, () => {
  nextTick(resizeSyncWidth)
}, {deep: true})

useSyncScrollY(gGanttChartRef, gGanttChartRef)

// 全体ビューの場合は遅延通知をフィルタして開く
const {open, filterFacility} = inject(GLOBAL_SCHEDULE_ALERT_KEY)!

const onClickBar = async (bar: GanttBarObject, e: MouseEvent, datetime?: string | Date) => {
  console.log("click-bar", bar, e, datetime)
  open()
  filterFacility.value = Number(bar.ganttBarConfig.id)
}

// TODO: 重複コード
const getRuneName = (lastName?: string, firstName?: string) => {
  let result = ""
  if(lastName != undefined) {
    result += lastName.substring(0,1)
  }
  if(firstName != undefined) {
    result += firstName.substring(0,1)
  }
  return result
}


</script>