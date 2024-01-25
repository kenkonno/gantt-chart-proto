<template>
  <nav class="navbar navbar-light bg-light">
    <div class="d-flex w-100">
      <b class="d-flex align-items-center">ビューの選択</b>
      <router-link to="/">
        <span class="material-symbols-outlined">edit</span>
        <span class="text">設備ビュー</span>
      </router-link>
      <router-link to="/all-view">
        <span class="material-symbols-outlined">travel_explore</span>
        <span class="text">全体ビュー</span>
      </router-link>
      <schedule-alert></schedule-alert>
      <div>
      </div>
      <div>
        <label>
          受注状況：
          <label v-for="(name, code) in FacilityTypeMap" :key="code" >
            <input type="checkbox" name="facilityType" :value="code" v-model="globalState.facilityTypes" @change="changeFacilityType"/>
            {{name}}
          </label>
        </label>
      </div>
      <DepartmentUserFilter></DepartmentUserFilter>
      <a href="#" @click.prevent="updateFacility">
        <span class="material-symbols-outlined">refresh</span>
        <span class="text">リロード</span>
      </a>
      <a href="#" @click.prevent="modalIsOpen = true" style="margin-left: auto;">
        <span class="material-symbols-outlined">person</span>
        <span class="text">{{ userInfo.name }}</span>
      </a>
    </div>
  </nav>
  <nav class="navbar navbar-light bg-light" v-if="allowed('ALL_SETTINGS')">
    <div>
      <b>全体設定</b>
      <ModalWithLink title="設備一覧" icon="precision_manufacturing">
        <facility-view @update="updateFacility"></facility-view>
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
  </nav>
  <Suspense v-if="modalIsOpen">
    <default-modal title="担当者" @close-edit-modal="closeModalProxy">
      <async-user-edit :id="userInfo.id" @close-edit-modal="closeModalProxy"></async-user-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense>
    <router-view/>
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
</style>

<script setup lang="ts">
import {
  GLOBAL_ACTION_KEY,
  GLOBAL_GETTER_KEY,
  GLOBAL_MUTATION_KEY,
  GLOBAL_STATE_KEY,
  useGlobalState
} from "@/composable/globalState";
import {provide} from "vue";
import ScheduleAlert from "@/components/scheduleAlert/ScheduleAlert.vue";
import {FacilityTypeMap} from "@/const/common";
import router from "@/router";
import {GLOBAL_SCHEDULE_ALERT_KEY, useScheduleAlert} from "@/composable/scheduleAlert";
import UserView from "@/views/UserView.vue";
import FacilityView from "@/views/FacilityView.vue";
import ProcessView from "@/views/ProcessView.vue";
import ModalWithLink from "@/components/modal/ModalWithLink.vue";
import DepartmentView from "@/views/DepartmentView.vue";
import DepartmentUserFilter from "@/components/departmentUserFilter/DepartmentUserFilter.vue";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY, useDepartmentUserFilter} from "@/composable/departmentUserFilter";
import {allowed} from "@/composable/role";
import {getUserInfo, loggedIn} from "@/composable/auth";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import AsyncUserEdit from "@/components/user/AsyncUserEdit.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {initStateValue} from "@/utils/globalFilterState";

// ローカルストレージの初期化
initStateValue()

// GlobalStateのProvide
const {globalState, actions, mutations, getters} = await useGlobalState()
provide(GLOBAL_STATE_KEY, globalState.value)
provide(GLOBAL_ACTION_KEY, actions)
provide(GLOBAL_MUTATION_KEY, mutations)
provide(GLOBAL_GETTER_KEY, getters)

const globalScheduleAlert = useScheduleAlert(globalState.value.scheduleAlert)
provide(GLOBAL_SCHEDULE_ALERT_KEY, globalScheduleAlert)

const globalDepartmentUserFilter = useDepartmentUserFilter()
provide(GLOBAL_DEPARTMENT_USER_FILTER_KEY, globalDepartmentUserFilter)

const userInfo = await getUserInfo()!

const changeFacilityType = () => {
  // 設備ビューの時はpileUpsだけ
  if (router.currentRoute.value.name == "gantt") {
    mutations.refreshPileUpsRefresh()
  }
  if (router.currentRoute.value.name == "gantt-all-view") {
    mutations.refreshGanttAll()
  }
}

const updateFacility = () => {
  actions.refreshFacilityList();
  if (router.currentRoute.value.name == "gantt") {
    mutations.refreshGantt(globalState.value.currentFacilityId, false)
  }
  if (router.currentRoute.value.name == "gantt-all-view") {
    mutations.refreshGanttAll()
  }
}

// profile関連
const {modalIsOpen, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  closeEditModal()
  const {user} = await loggedIn()
  if (user != undefined) {
    userInfo.name = user.name
  }
}

</script>