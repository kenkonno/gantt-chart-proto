<template>
  <nav class="navbar navbar-light bg-light">
    <div>
      <b>全体設定</b>
      <ModalWithLink title="設備一覧" icon="precision_manufacturing">
        <facility-view @update="refreshFacilityList(); refreshGantt(globalState.currentFacilityId)"></facility-view>
      </ModalWithLink>
      <ModalWithLink title="工程一覧" icon="account_tree">
        <process-view></process-view>
      </ModalWithLink>
      <ModalWithLink title="部署一覧" icon="settings_accessibility">
        <department-view></department-view>
      </ModalWithLink>
      <ModalWithLink title="担当者一覧" icon="person">
        <user-view></user-view>
      </ModalWithLink>
    </div>
    <div v-if="facilityList.length > 0" style="width: 100%; text-align: left">
      <b>設備設定</b>
      <select style="display: inline" v-model.number="globalState.currentFacilityId" @input="refreshGantt(Number($event.target.value))">
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
  <gantt-proxy :facilityId="globalState.currentFacilityId" v-if="globalState.currentFacilityId > 0"></gantt-proxy>
  <div v-else>
    ユニットを選択してください。
  </div>
</template>
<style lang="scss" scoped>
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
import FacilityView from "@/views/FacilityView.vue";
import ModalWithLink from "@/components/modal/ModalWithLink.vue";
import ProcessView from "@/views/ProcessView.vue";
import DepartmentView from "@/views/DepartmentView.vue";
import UserView from "@/views/UserView.vue";
import UnitView from "@/views/UnitView.vue";
import OperationSettingView from "@/views/OperationSettingView.vue";
import HolidayView from "@/views/HolidayView.vue";
import {inject, nextTick} from "vue";
import GanttProxy from "@/components/GanttProxy.vue";
import {
  GLOBAL_ACTION_KEY,
  GLOBAL_MUTATION_KEY,
  GLOBAL_STATE_KEY,
} from "@/composable/globalState";

const globalState = inject(GLOBAL_STATE_KEY)!
const facilityList = globalState.facilityList
const {refreshFacilityList, refreshHolidayMap, refreshUnitMap, refreshOperationSettingMap} = inject(GLOBAL_ACTION_KEY)!
const {updateCurrentFacilityId} = inject(GLOBAL_MUTATION_KEY)!

// たぶんwatchしてガントチャートの切り替えにしたほうがいい気がする。
const refreshGantt = async (facilityId: number) => {
  updateCurrentFacilityId(0)
  // facility紐づくデータを初期化する
  await refreshHolidayMap(facilityId)
  await refreshUnitMap(facilityId)
  await refreshOperationSettingMap(facilityId)
  nextTick(() => {
    updateCurrentFacilityId(facilityId)
  })
}

</script>