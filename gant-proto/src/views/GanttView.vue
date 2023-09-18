<template>
  <nav class="navbar navbar-light bg-light">
    <div>
      <b>全体設定</b>
      <ModalWithLink title="設備一覧">
        <facility-view @update="refreshFacilityList"></facility-view>
      </ModalWithLink>
      <ModalWithLink title="工程一覧">
        <process-view></process-view>
      </ModalWithLink>
      <ModalWithLink title="部署一覧">
        <department-view></department-view>
      </ModalWithLink>
      <ModalWithLink title="担当者一覧">
        <user-view></user-view>
      </ModalWithLink>
      <p style="display: inline">全体スケジュールビュー</p>
    </div>
    <div v-if="facilityList.length > 0" style="width: 100%; text-align: left">
      <b>設備設定</b>
      <select style="display: inline" v-model.number="globalState.currentFacilityId" @input="refreshGantt(Number($event.target.value))">
        <option v-for="item in facilityList" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <ModalWithLink title="ユニット一覧" :disabled="globalState.currentFacilityId===-1">
        <unit-view></unit-view>
      </ModalWithLink>
      <ModalWithLink title="稼働設定" :disabled="globalState.currentFacilityId===-1">
        <operation-setting-view></operation-setting-view>
      </ModalWithLink>
      <ModalWithLink title="休日設定" :disabled="globalState.currentFacilityId===-1">
        <holiday-view></holiday-view>
      </ModalWithLink>
      <p style="display: inline">権限設定</p>
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
import {nextTick, provide} from "vue";
import GanttProxy from "@/components/GanttProxy.vue";
import {GLOBAL_ACTION_KEY, GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY, useGlobalState} from "@/composable/globalState";

const {globalState, actions, mutations} = await useGlobalState()
provide(GLOBAL_STATE_KEY, globalState.value)
provide(GLOBAL_ACTION_KEY, actions)
provide(GLOBAL_MUTATION_KEY, mutations)

const {facilityList} = globalState.value
const {refreshFacilityList} = actions
const {updateCurrentFacilityId} = mutations

// たぶんwatchしてガントチャートの切り替えにしたほうがいい気がする。
const refreshGantt = (facilityId: number) => {
  updateCurrentFacilityId(0)
  nextTick(() => {
    updateCurrentFacilityId(facilityId)
  })
}

</script>