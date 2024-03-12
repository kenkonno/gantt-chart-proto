<template>
  <div class="wrapper d-flex department-user-filter-wrapper">
    <div>部署：</div>
    <multiselect
        class="select-box"
        :value="cSelectedDepartment"
        @input="inputDepartment"
        :searchable="true"
        :options="departmentOptions"
        :allow-empty="true"
        mode="tags"
        :max="1"
        @change="refresh"
    />
    <div>担当者：</div>
    <multiselect
        class="select-box"
        :value="cSelectedUser"
        @input="inputUser"
        :searchable="true"
        :options="userOptions"
        :allow-empty="true"
        :taggable="true"
        mode="tags"
        :max="1"
        @change="refresh"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, inject} from "vue";
import {GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import Multiselect from "@vueform/multiselect";
import router from "@/router";

const globalState = inject(GLOBAL_STATE_KEY)!
const {selectedUser, selectedDepartment} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!
const {refreshGantt, refreshGanttAll} = inject(GLOBAL_MUTATION_KEY)!

const departmentOptions = computed(() => {
  return globalState.departmentList.map(v => {
    return {value: v.id, label: v.name}
  })
})
const userOptions = computed(() => {
  return globalState.userList.filter(v => {
    if (selectedDepartment.value == undefined) {
      return true
    } else {
      return v.department_id == selectedDepartment.value
    }
  }).map(v => {
    return {value: v.id, label: v.name}
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

const cSelectedDepartment = computed(() => {
  return [selectedDepartment.value]
})
const cSelectedUser = computed(() => {
  return [selectedUser.value]
})


const inputDepartment = (v: number[]) => {
  selectedDepartment.value = v[0]
}
const inputUser = (v: number[]) => {
  selectedUser.value = v[0]
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