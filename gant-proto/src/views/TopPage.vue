<template>
  <nav>
    <div>
      <b>全体設定</b>
      <modal-with-link title="設備一覧" :disabled="false">
        <facility-view @update="facilityRefresh"></facility-view>
      </modal-with-link>
      <modal-with-link title="工程一覧" :disabled="false">
        <process-view></process-view>
      </modal-with-link>
      <modal-with-link title="部署一覧" :disabled="false">
        <department-view></department-view>
      </modal-with-link>
      <modal-with-link title="担当者一覧" :disabled="false">
        <user-view></user-view>
      </modal-with-link>
      <p style="display: inline">全体スケジュールビュー(未実装)</p>
    </div>
    <div v-if="facilityList.length > 0">
      <select style="display: inline" v-model="currentFacilityId">
        <option v-for="item in facilityList" :key="item.id" :value="item.id">{{ item.name }}</option>
      </select>
      <modal-with-link title="ユニット一覧" :disabled="currentFacilityId===-1">
        <unit-view :facility-id="currentFacilityId"></unit-view>
      </modal-with-link>
      <modal-with-link title="稼働設定" :disabled="currentFacilityId===-1">
        <operation-setting-view :facility-id="currentFacilityId"></operation-setting-view>
      </modal-with-link>
      <modal-with-link title="休日設定" :disabled="currentFacilityId===-1">
        <holiday-view :facilityId="currentFacilityId"></holiday-view>
      </modal-with-link>
      <p style="display: inline">権限設定</p>
    </div>
    <div v-else>設備の設定がありません。設備一覧から追加してください。</div>
  </nav>
  <Suspense>
    <router-view/>
  </Suspense>
</template>

<style lang="scss" scoped>
nav {
  padding: 30px;

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
import {ref} from "vue";
import {useFacilityTable} from "@/composable/facility";
import {useGanttGroup, useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";

const {list: facilityList, refresh: facilityRefresh} = await useFacilityTable()
const currentFacilityId = ref<number>(-1)
</script>