<template>
  <div class="wrapper d-flex">
    <div>部署：</div>
    <multiselect
        class="select-box"
        v-model="selectedDepartment"
        :searchable="true"
        :options="departmentOptions"
        :allow-empty="true"
        @change="refresh"
    />
    <div>担当者：</div>
    <multiselect
        class="select-box"
        v-model="selectedUser"
        :searchable="true"
        :options="userOptions"
        :allow-empty="true"
        @change="refresh"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, inject} from "vue";
import {GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import Multiselect from "@vueform/multiselect";

const globalState = inject(GLOBAL_STATE_KEY)!
const {selectedUser, selectedDepartment} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!
const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!

const departmentOptions = computed(() => {
  return globalState.departmentList.map(v => {
    return {value: v.id, label: v.name}
  })
})
const userOptions = computed(() => {
  return globalState.userList.filter(v => {
    if ( selectedDepartment.value == undefined) {
      return true
    } else {
      return v.department_id == selectedDepartment.value
    }
  }).map(v => {
    return {value: v.id, label: v.name}
  })
})

const refresh = () => {
  console.log(globalState.currentFacilityId)
  refreshGantt(globalState.currentFacilityId, false)
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.wrapper {
  white-space: nowrap;
}
.select-box {
  min-width: 12rem;
  z-index: 301;
}
</style>




