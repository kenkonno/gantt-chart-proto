<template>
  <div class="container">
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>ユニット名</th>
        <th v-for="item in processList" :key="item.id">
          {{ item.name }}
        </th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td>{{ item.id }}</td>
        <td>{{ unitMap[item.unit_id] }}</td>
        <td v-for="process in processList" :key="process.id">
          <input type="number" v-model="getWorkHour(item.unit_id!, process.id!).work_hour">
        </td>
      </tr>
      </tbody>
    </table>
    <button type="submit" class="btn btn-primary" @click="postOperationSettingById(facilityId, list, $emit)">更新
    </button>
  </div>
</template>

<script setup lang="ts">
import {OperationSetting, Process, Unit} from "@/api";
import {postOperationSettingById} from "@/composable/operationSetting";

defineEmits(['openEditModal', 'closeEditModal'])

interface AsyncOperationSettingTable {
  list: OperationSetting[]
  unitList: Unit[]
  processList: Process[]
  facilityId: number
}

const props = defineProps<AsyncOperationSettingTable>()
const unitMap: { [x: number]: string; } = {}
props.unitList.forEach(v => {
  unitMap[v.id!] = v.name
})
const getWorkHour = (unitId: number, processId: number) => {
  return props.list.find(v => v.unit_id === unitId)!.workHours.find(v => v.process_id === processId)
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




