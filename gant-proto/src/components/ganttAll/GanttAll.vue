<template>
  <div class="gantt-wrapper" id="gantt-all-view" :class="{withFilter:hasFilter()}">
    <div class="gantt-facility-wrapper d-flex overflow-x-scroll hide-scroll" ref="ganttWrapperElement">
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
          :highlighted-dates="holidaysAsDate"
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
                        class="pointer"><u>{{ item.facility.name }}</u></gantt-td>
              <gantt-td :visible="ganttAllHeader[1].visible">
                <div class="user-wrapper">
                  <SingleRune v-for="user in item.users" :key="user.id" :name="user.name" :id="user.id"></SingleRune>
                </div>
              </gantt-td>
              <gantt-td :visible="ganttAllHeader[2].visible">{{ item.startDate }}</gantt-td>
              <gantt-td :visible="ganttAllHeader[3].visible">{{ item.endDate }}</gantt-td>
              <gantt-td :visible="ganttAllHeader[4].visible">{{ item.estimate }}</gantt-td>
              <gantt-td :visible="ganttAllHeader[5].visible">{{ item.progress_percent }}</gantt-td>
            </tr>
            </tbody>
          </table>
        </template>
        <g-gantt-row v-for="item in ganttAllRow" :key="item.facility.id" :bars="item.bars"/>
      </g-gantt-chart>
    </div>
    <!-- 山積み部分 -->
    <hr>
    <div class="gantt-facility-pile-ups-wrapper d-flex overflow-x-scroll" ref="childGanttWrapperElement">
      <PileUps
          v-if="allowed('VIEW_PILEUPS')"
          :chart-start="chartStart"
          :chart-end="chartEnd"
          :display-type="displayType"
          :holidays="holidays"
          :tickets="[]"
          :ticket-users="[]"
          :width="getGanttChartWidth(displayType)"
          :highlightedDates="holidaysAsDate"
          :syncWidth="syncWidth"
          :current-facility-id="-1"
          :milestone-vertical-lines="[]"
          @on-mounted="forceScroll"
      >
      </PileUps>
    </div>
  </div>
</template>
<style>
@import '@/assets/gantt-override.scss';

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
import {GLOBAL_MUTATION_KEY} from "@/composable/globalState";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import {allowed} from "@/composable/role";

const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!

type GanttAllProps = {
  ganttAllHeader: Header[],
  displayType: DisplayType,
}
const props = defineProps<GanttAllProps>()

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
} = await useGanttAll()

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

</script>