<template>
  <nav class="navbar navbar-light bg-light">
    <div v-if="facilityList.length > 0" style="width: 100%; text-align: left">
      <b>設備設定</b>
      <select style="display: inline" v-model.number="globalState.currentFacilityId"
              @input="refreshGantt(Number($event.target.value))">
        <option v-for="item in facilityList" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <ModalWithLink title="ユニット一覧" icon="switch_access" :disabled="globalState.currentFacilityId===-1">
        <unit-view></unit-view>
      </ModalWithLink>
      <ModalWithLink title="稼働設定" icon="timer" :disabled="globalState.currentFacilityId===-1">
        <operation-setting-view></operation-setting-view>
      </ModalWithLink>
      <ModalWithLink title="休日設定" icon="holiday_village" :disabled="globalState.currentFacilityId===-1">
        <holiday-view></holiday-view>
      </ModalWithLink>
      <ModalWithLink title="権限設定" icon="folder_supervised" :disabled="globalState.currentFacilityId===-1">
      </ModalWithLink>
    </div>
    <div v-else>設備の設定がありません。設備一覧から追加してください。</div>
  </nav>

  <div style="display:none">{{gantFacility != undefined}} vuejshack</div>
  <div v-show="globalState.currentFacilityId > 0">
    <gantt-facility-menu
        :gantt-facility-header="GanttHeader"
        :display-type="displayType"
        @set-schedule-by-from-to="gantFacility.setScheduleByFromToProxy()"
        @set-schedule-by-person-day="gantFacility.setScheduleByPersonDayProxy()"
        @updateDisplayType="updateDisplayType"
    ></gantt-facility-menu>
  </div>
  <Suspense>
    <gantt-facility
        ref="gantFacility"
        v-if="globalState.currentFacilityId > 0 && globalState.ganttFacilityRefresh"
        :gantt-facility-header="GanttHeader"
        :display-type="displayType"
    >
    </gantt-facility>
    <div v-else>
      ユニットを選択してください。
    </div>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>
<style lang="scss" scoped>
.navbar {
  padding: 0;
  height: 30px;
  font-size: 0.8rem;

  > div {
    > a, div {
      display: block;
      margin-left: 5px;
      color: inherit;
      padding: 0;
      text-decoration: inherit;
      border-bottom: 1px solid black;

      > .material-symbols-outlined {
        vertical-align: middle;
        font-size: 1rem;
      }

      .text {
        vertical-align: middle;
      }
    }

  }
}

nav {
  padding: 0 0 5px 5px;

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
</style>
<script setup lang="ts">
import {inject, ref} from "vue";
import GanttFacility from "@/components/ganttFacility/GanttFacility.vue";
import {
  GLOBAL_MUTATION_KEY,
  GLOBAL_STATE_KEY,
} from "@/composable/globalState";
import GanttFacilityMenu from "@/components/ganttFacility/GanttFacilityMenu.vue";
import {DisplayType, useGanttFacilityMenu} from "@/composable/ganttFacilityMenu";
import OperationSettingView from "@/views/OperationSettingView.vue";
import ModalWithLink from "@/components/modal/ModalWithLink.vue";
import UnitView from "@/views/UnitView.vue";
import HolidayView from "@/views/HolidayView.vue";
import {computed} from "vue";
import {FacilityStatus} from "@/const/common";

const globalState = inject(GLOBAL_STATE_KEY)!
const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!

const {GanttHeader, displayType} = useGanttFacilityMenu()
const gantFacility = ref(null)
const updateDisplayType = (v: DisplayType) => {
  displayType.value = v
}

const facilityList = computed(() => {
  return globalState.facilityList.filter(v => v.status === FacilityStatus.Enabled)
})


// たぶんwatchしてガントチャートの切り替えにしたほうがいい気がする。
</script>