<template>
  <div class="wrapper d-flex department-user-filter-wrapper">
    <div>部署：</div>
    <CustomMultiselect
        class="select-box"
        v-model="selectedDepartment"
        :options="departmentOptions"
        @update:modelValue="refresh"
        :max="max"
        placeholder="部署を選択"
        no-options-text="部署を登録してください"
    />
    <div>担当者：</div>
    <CustomMultiselect
        class="select-box"
        v-model="selectedUser"
        :options="userOptions"
        @update:modelValue="refresh"
        :max="max"
        placeholder="担当者を選択"
        no-options-text="担当者を登録してください"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, inject} from "vue";
import {GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import router from "@/router";
import {available} from "@/composable/featureOption";
import {FeatureOption} from "@/const/common";
import CustomMultiselect from "@/components/form/CustomMultiselect.vue";

const max = computed(() => {
  return available(FeatureOption.MultiSelectFilter) ? -1 : 1
})

const globalState = inject(GLOBAL_STATE_KEY)!
const {selectedUser, selectedDepartment} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!
const {refreshGantt, refreshGanttAll} = inject(GLOBAL_MUTATION_KEY)!

const departmentOptions = computed(() => {
  return globalState.departmentList.map(v => {
    return {value: v.id, label: v.name, runeName: v.name.substring(0, 1)}
  })
})

const userOptions = computed(() => {
  return globalState.userList.filter(v => {
    if (selectedDepartment.value.length == 0) {
      return true
    } else {
      return selectedDepartment.value.includes(v.department_id)
    }
  }).map(v => {
    const runeName = (v.lastName ? v.lastName.substring(0, 1) : "") + (v.firstName ? v.firstName.substring(0, 1) : "")
    return {value: v.id, label: `${v.lastName} ${v.firstName}`, runeName: runeName}
  })
})

const refresh = () => {
  if (router.currentRoute.value.name == "gantt") {
    refreshGantt(globalState.currentFacilityId, false)
  }
  if (router.currentRoute.value.name == "gantt-all-view") {
    refreshGanttAll()
  }
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.wrapper {
  white-space: nowrap;
}

</style>

<style>
.department-user-filter-wrapper .select-box {
  min-width: 12rem;
  z-index: 301;
  min-height: 24px;
  height: 24px;
}
.department-user-filter-wrapper .multiselect-tags-search{
  min-height: 22px;
  height: 22px;
  top: 1px;
}
</style>