<template>
  <nav class="navbar navbar-light bg-light">
    <div class="d-flex">
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
        <label>
          受注状況
          <label v-for="(name, code) in FacilityTypeMap" :key="code" >
            {{name}}
            <input type="checkbox" name="facilityType" :value="code" v-model="globalState.facilityTypes" @change="changeFacilityType"/>
          </label>
        </label>
      </div>
    </div>
  </nav>
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
import {FacilityStatusMap, FacilityTypeMap} from "@/const/common";
import router from "@/router";

const {globalState, actions, mutations, getters} = await useGlobalState()
provide(GLOBAL_STATE_KEY, globalState.value)
provide(GLOBAL_ACTION_KEY, actions)
provide(GLOBAL_MUTATION_KEY, mutations)
provide(GLOBAL_GETTER_KEY, getters)
const changeFacilityType = () => {
  // 設備ビューの時はpileUpsだけ
  if (router.currentRoute.value.name == "gantt") {
    console.log(router.currentRoute.value.name)
    console.log(globalState.value.pileUpsRefresh)
    mutations.refreshPileUpsRefresh()
    console.log(globalState.value.pileUpsRefresh)
  }
  if (router.currentRoute.value.name == "gantt-all-view") {
    mutations.refreshGanttAll()
  }
}
</script>