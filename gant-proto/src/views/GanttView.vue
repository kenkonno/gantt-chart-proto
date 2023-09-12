<template>
  <nav class="navbar navbar-light bg-light">
    <div>
      <b>全体設定</b>
      <modal-with-link title="設備一覧">
        <facility-view @update="facilityRefresh"></facility-view>
      </modal-with-link>
      <modal-with-link title="工程一覧">
        <process-view @update="refreshGantt(currentFacilityId)"></process-view>
      </modal-with-link>
      <modal-with-link title="部署一覧">
        <department-view @update="refreshGantt(currentFacilityId)"></department-view>
      </modal-with-link>
      <modal-with-link title="担当者一覧">
        <user-view @update="refreshGantt(currentFacilityId)"></user-view>
      </modal-with-link>
      <p style="display: inline">全体スケジュールビュー(未実装)</p>
    </div>
    <div v-if="facilityList.length > 0" style="width: 100%; text-align: left">
      <b>設備設定</b>
      <select style="display: inline" v-model="currentFacilityId" @input="refreshGantt($event.target.value)">
        <option v-for="item in facilityList" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <modal-with-link title="ユニット一覧" :disabled="currentFacilityId===-1">
        <unit-view :facility-id="currentFacilityId" @update="refreshGantt(currentFacilityId)"></unit-view>
      </modal-with-link>
      <modal-with-link title="稼働設定" :disabled="currentFacilityId===-1">
        <operation-setting-view :facility-id="currentFacilityId" @update="refreshGantt(currentFacilityId)"></operation-setting-view>
      </modal-with-link>
      <modal-with-link title="休日設定" :disabled="currentFacilityId===-1">
        <holiday-view :facilityId="currentFacilityId" @update="refreshGantt(currentFacilityId)"></holiday-view>
      </modal-with-link>
      <p style="display: inline">権限設定</p>
    </div>
    <div v-else>設備の設定がありません。設備一覧から追加してください。</div>
  </nav>
  <gantt-proxy :facilityId="currentFacilityId" v-if="currentFacilityId > 0"></gantt-proxy>
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
import DefaultModal from "@/components/modal/DefaultModal.vue";
import ModalWithLink from "@/components/modal/ModalWithLink.vue";
import ProcessView from "@/views/ProcessView.vue";
import DepartmentView from "@/views/DepartmentView.vue";
import UserView from "@/views/UserView.vue";
import UnitView from "@/views/UnitView.vue";
import OperationSettingView from "@/views/OperationSettingView.vue";
import HolidayView from "@/views/HolidayView.vue";
import {nextTick, ref} from "vue";
import {useFacilityTable} from "@/composable/facility";
import {useGanttGroup, useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";

const {list: facilityList, refresh: facilityRefresh} = await useFacilityTable()
const currentFacilityId = ref<number>(-1)

const refreshGantt = (facilityId: number) => {
  currentFacilityId.value = 0
  nextTick(() => {
    currentFacilityId.value = facilityId
  })
}

import GanttProxy from "@/components/GanttProxy.vue";
</script>