<template>
  <nav>
    <div>
      <b>全体設定</b>
      <Suspense>
        <modal-with-link title="設備一覧">
          <facility-view></facility-view>
        </modal-with-link>
      </Suspense>
      <Suspense>
        <modal-with-link title="工程一覧">
          <process-view></process-view>
        </modal-with-link>
      </Suspense>
      <Suspense>
        <modal-with-link title="部署一覧">
          <department-view></department-view>
        </modal-with-link>
      </Suspense>
      <Suspense>
        <modal-with-link title="担当者一覧">
          <user-view></user-view>
        </modal-with-link>
      </Suspense>
      <p style="display: inline">全体スケジュールビュー(未実装)</p>
    </div>
    <div>
      <select style="display: inline">
        <option>AAAA</option>
      </select>
      <Suspense>
        <modal-with-link title="ユニット一覧">
          <unit-view></unit-view>
        </modal-with-link>
      </Suspense>
      <Suspense>
        <modal-with-link title="稼働設定">
          <operation-setting-view></operation-setting-view>
        </modal-with-link>
      </Suspense>
      <Suspense>
        <modal-with-link title="休日設定">
          <holiday-view></holiday-view>
        </modal-with-link>
      </Suspense>
      <p style="display: inline">権限設定</p>
    </div>
  </nav>
  <Suspense>
    <router-view/>
  </Suspense>
</template>

<style lang="scss" scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

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
</script>