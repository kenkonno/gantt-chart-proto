<template>
  <nav class="navbar navbar-light bg-warning" v-if="isSimulateUser">
    <div class="d-flex w-100">
      <b class="m-auto">シミュレーション中</b>
    </div>
  </nav>
  <nav class="navbar navbar-light bg-light" v-if="allowed('MENU')">
    <div class="d-flex w-100">
      <b class="d-flex align-items-center">ビューの選択</b>
      <router-link to="/">
        <span class="material-symbols-outlined">edit</span>
        <span class="text">案件ビュー</span>
      </router-link>
      <router-link to="/all-view" v-if="allowed('ALL_VIEW')">
        <span class="material-symbols-outlined">travel_explore</span>
        <span class="text">全体ビュー</span>
      </router-link>
      <schedule-alert v-if="allowed('VIEW_SCHEDULE_ALERT')"></schedule-alert>
      <div>
      </div>
      <FacilityTypeFilter
          v-model="globalState.facilityTypes"
          @change="changeFacilityType"
      />
      <DepartmentUserFilter></DepartmentUserFilter>
      <div v-if="allowed('VIEW_PILEUPS')">
        <label>
          <label>
            <input type="checkbox" name="facilityType" :value="true" v-model="globalState.showPileUp" />
            山積み表示
          </label>
        </label>
      </div>
      <a href="#" @click.prevent="updateFacility">
        <span class="material-symbols-outlined">refresh</span>
        <span class="text">リロード</span>
      </a>
      <a href="#" @click.prevent="modalIsOpen = true" style="margin-left: auto;">
        <span class="material-symbols-outlined">person</span>
        <span class="text">{{ `${userInfo.lastName} ${userInfo.firstName}` }}</span>
      </a>
    </div>
  </nav>
  <nav class="navbar navbar-light bg-light" v-if="allowed('ALL_SETTINGS')">
    <div>
      <b>全体設定</b>
      <ModalWithLink title="案件一覧" icon="precision_manufacturing">
        <facility-view @update="updateFacility"></facility-view>
      </ModalWithLink>
      <ModalWithLink title="工程一覧" icon="account_tree" v-if="allowed('UPDATE_MASTER')">
        <process-view @update="updateFacility"></process-view>
      </ModalWithLink>
      <ModalWithLink title="部署一覧" icon="settings_accessibility" v-if="allowed('UPDATE_MASTER')">
        <department-view @update="updateFacility"></department-view>
      </ModalWithLink>
      <ModalWithLink title="担当者一覧" icon="person" v-if="allowed('UPDATE_MASTER')">
        <user-view @update="updateFacility"></user-view>
      </ModalWithLink>
      <ModalWithLink title="シミュレーション" icon="timeline">
        <simulation-view @update="updateSimulation"></simulation-view>
      </ModalWithLink>
    </div>
  </nav>
  <Suspense v-if="modalIsOpen">
    <default-modal title="担当者" @close-edit-modal="closeModalProxy">
      <async-user-edit :id="userInfo.id" @close-edit-modal="closeModalProxy" :mode="'profile'"></async-user-edit>
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
.navbar-container {
  display: flex;
  border-bottom: 1px solid #dee2e6;
}

.company-logo-area {
  width: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 5px;
  background-color: #f8f9fa;
}

.company-logo {
  max-width: 100%;
  max-height: 50px;
  object-fit: contain;
}

.navbars-wrapper {
  flex-grow: 1;
}

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
import {provide, ref} from "vue";
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
import SimulationView from "@/views/SimulationView.vue";
import ServiceLogo from "@/components/logo/ServiceLogo.vue";
import FacilityTypeFilter from "@/components/form/FacilityTypeFilter.vue";

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

let {userInfo: tempUserInfo, isSimulateUser: tempIsSimulateUser} = getUserInfo()!
const userInfo = ref(tempUserInfo)
const isSimulateUser = ref(tempIsSimulateUser)

const changeFacilityType = () => {
  // 案件ビューの時はpileUpsだけ
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
    userInfo.value.lastName = user.lastName
    userInfo.value.firstName = user.firstName
  }
}

// simulation関連
const updateSimulation = async () => {
  await loggedIn()
  let {userInfo: tempUserInfo, isSimulateUser: tempIsSimulateUser} = getUserInfo()!
  userInfo.value = tempUserInfo
  isSimulateUser.value = tempIsSimulateUser
  console.log('AAAAAAAAAAAAA', isSimulateUser.value, tempIsSimulateUser)
}

</script>