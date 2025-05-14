<template>
  <nav class="navbar navbar-light bg-light" v-if="allowed('MENU')">
    <div v-if="facilityList.length > 0" style="width: 100%; text-align: left">
      <b>案件設定</b>
      <select style="display: inline" v-model.number="globalState.currentFacilityId"
              @input="refreshGantt(Number($event.target.value))">
        <option v-for="item in facilityList" :key="item.id" :value="item.id">{{ item.name }}<template v-if="item.type === FacilityType.Ordered">✅</template></option>
      </select>
      <template v-if="allowed('FACILITY_SETTINGS')">
        <ModalWithLink title="ユニット一覧" icon="switch_access" :disabled="globalState.currentFacilityId===-1">
          <unit-view></unit-view>
        </ModalWithLink>
        <ModalWithLink title="稼働設定" icon="timer" :disabled="globalState.currentFacilityId===-1" v-if="false">
          <operation-setting-view></operation-setting-view>
        </ModalWithLink>
        <ModalWithLink title="稼働日設定" icon="timer" :disabled="globalState.currentFacilityId===-1">
          <facility-work-schedule-view></facility-work-schedule-view>
        </ModalWithLink>
        <ModalWithLink title="マイルストーン" icon="folder_supervised" :disabled="globalState.currentFacilityId===-1">
          <milestone-view></milestone-view>
        </ModalWithLink>
        <ModalWithLink title="共有リンク" icon="public" :disabled="globalState.currentFacilityId===-1">
          <facility-shared-link-view></facility-shared-link-view>
        </ModalWithLink>
      </template>
    </div>
    <div v-else>案件の設定がありません。案件一覧から追加してください。</div>
  </nav>

  <div style="display:none">{{ gantFacility != undefined }} vuejshack</div>
  <div v-show="globalState.currentFacilityId > 0 && globalState.processList.length > 0" class="gantt-facility-menu">
    <gantt-facility-menu
        :gantt-facility-header="GanttHeader"
        :display-type="displayType"
        @set-schedule-by-from-to="confirm(gantFacility.setScheduleByFromToProxy)"
        @set-schedule-by-person-day="confirm(gantFacility.setScheduleByPersonDayProxy)"
        @updateDisplayType="updateDisplayType"
    ></gantt-facility-menu>
  </div>
  <div v-if="globalState.currentFacilityId <= 0">
    案件を選択してください。
  </div>
  <div v-if="globalState.processList.length <= 0">
    工程を登録してください。
  </div>
  <div></div>
  <Suspense v-if="globalState.currentFacilityId > 0 && globalState.processList.length > 0 && globalState.ganttFacilityRefresh">
    <gantt-facility
        ref="gantFacility"
        :gantt-facility-header="GanttHeader"
        :display-type="displayType"
    >
    </gantt-facility>
    <template #fallback>
      <default-spinner></default-spinner>
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

.gantt-facility-menu {
  height: 30px;
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
import {FacilityStatus, FacilityType} from "@/const/common";
import DefaultSpinner from "@/components/spinner/DefaultSpinner.vue";
import {allowed} from "@/composable/role";
import MilestoneView from "@/views/MilestoneView.vue";
import Swal from "sweetalert2"
import {useRoute} from 'vue-router'
import FacilitySharedLinkView from "@/views/FacilitySharedLinkView.vue";
import FacilityWorkScheduleView from "@/views/FacilityWorkScheduleView.vue";

const globalState = inject(GLOBAL_STATE_KEY)!
const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!

const {GanttHeader, displayType} = useGanttFacilityMenu()
const gantFacility = ref(null)

// パラメーターが存在していれば表示する設備を変更する
const route = useRoute()
const defaultFacilityId = parseInt(route.query.facilityId, 10)
if (!isNaN(defaultFacilityId)) {
  refreshGantt(defaultFacilityId, false)
}


const updateDisplayType = (v: DisplayType) => {
  displayType.value = v
}

const facilityList = computed(() => {
  let result = globalState.facilityList.filter(v => v.status === FacilityStatus.Enabled)
  if (globalState.facilityTypes.length > 0) {
    result = result.filter(v => globalState.facilityTypes.includes(v.type))
  }

  return result.sort((a, b) => (b.order ? b.order : 0) < (a.order ? a.order : 0) ? -1 : 1);
})

const confirm = async (func : ()=> any) => {
  let result = await Swal.fire({
    title: 'リスケ機能実行の確認',
    text: "スケジュールを一括で変更しますがよろしいでしょうか？この変更は元に戻すことはできません。",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#3085d6',
    cancelButtonColor: '#d33',
    confirmButtonText: '実行',
    cancelButtonText: 'キャンセル'
  });
  if (result.isConfirmed) {
    func()
  }
}


// たぶんwatchしてガントチャートの切り替えにしたほうがいい気がする。
</script>