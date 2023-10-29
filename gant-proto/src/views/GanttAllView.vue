<template>
  <div class="action-menu d-flex">
    <div class="wrapper d-flex">
      <div class="justify-middle">
        <div>メニュー</div>
      </div>
    </div>
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
          color-scheme="creamy"
          :highlighted-dates="holidaysAsDate"
          sticky
          :display-today-line="true"
          @today-line-position-x="initScroll($event, ganttWrapperElement)"
      >
        <template #side-menu>
          <table class="side-menu" ref="ganttSideMenuElement">
            <thead class="side-menu-header">
            <tr>
              <th class="side-menu-cell"></th><!-- css hack min-height -->
              <th :colspan="GanttHeader.length">
                「全体」のスケジュール：{{ startDate }}～{{ endDate }}
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
            <tr v-for="item in ganttAllRow" :key="item.facility.id">
              <td class="side-menu-cell"></td><!-- css hack min-height -->
              <gantt-td :visible="GanttHeader[0].visible">{{ item.facility.name }}</gantt-td>
              <gantt-td :visible="GanttHeader[1].visible">
                <div class="user-wrapper">
                  <SingleRune v-for="user in item.users" :key="user.id" :name="user.name" :id="user.id"></SingleRune>
                </div>
              </gantt-td>
              <gantt-td :visible="GanttHeader[2].visible">{{ item.startDate }}</gantt-td>
              <gantt-td :visible="GanttHeader[3].visible">{{ item.endDate }}</gantt-td>
              <gantt-td :visible="GanttHeader[4].visible">{{ item.estimate }}</gantt-td>
              <gantt-td :visible="GanttHeader[5].visible">{{ item.progress_percent }}</gantt-td>
            </tr>
            </tbody>
          </table>
        </template>
        <g-gantt-row v-for="item in ganttAllRow" :key="item.facility.id" :bars="item.bars"/>
      </g-gantt-chart>
    </div>
    <!-- 山積み部分 -->
    <hr>
    <div class="d-flex overflow-x-scroll" ref="childGanttWrapperElement">
      <PileUps :chart-start="chartStart"
               :chart-end="chartEnd"
               :display-type="displayType"
               :holidays="holidays"
               :tickets="[]"
               :ticket-users="[]"
               :width="getGanttChartWidth"
               :highlightedDates="holidaysAsDate"
               :syncWidth="syncWidth"
               @on-mounted="forceScroll"
      >
      </PileUps>
    </div>
  </div>
</template>
<style>
@import '@/assets/gantt-override.scss';
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
import {nextTick} from "vue";
import {GGanttChart, GGanttRow} from "@infectoone/vue-ganttastic";
import {DAYJS_FORMAT} from "@/utils/day";
import GanttTd from "@/components/gantt/GanttTd.vue";
import PileUps from "@/components/pileUps/PileUps.vue";
import {useGanttAll} from "@/composable/ganttAll";
import {useSyncWidthAndScroll} from "@/composable/syncWidth";
import {ref, watch} from "vue";
import AccordionHorizontal from "@/components/accordionHorizontal/AccordionHorizontal.vue";
import SingleRune from "@/components/form/SingleRune.vue";
import {initScroll} from "@/utils/initScroll";

const {
  GanttHeader,
  startDate,
  endDate,
  ganttAllRow,
  holidaysAsDate,
  displayType,
  getGanttChartWidth,
  tickets,
  ticketUsers,
  holidays,
  chartStart,
  chartEnd
} = await useGanttAll()

// スクロールや大きさの同期系
const ganttSideMenuElement = ref<HTMLDivElement>()
const ganttWrapperElement = ref<HTMLDivElement>()
const childGanttWrapperElement = ref<HTMLDivElement>()

const {
  syncWidth,
  resizeSyncWidth,
  forceScroll
} = useSyncWidthAndScroll(ganttSideMenuElement, ganttWrapperElement, childGanttWrapperElement)
watch(GanttHeader, () => {
  nextTick(resizeSyncWidth)
}, {deep: true})

</script>