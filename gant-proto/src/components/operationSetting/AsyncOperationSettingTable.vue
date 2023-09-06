<template>
  <div class="container">
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>User</th>
        <th>Unit</th>
        <th v-for="item in processList" :key="item.id">
          {{ item.name }}
        </th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td>{{ item.id }}</td>
        <td>{{ userMap[item.user_id] }}</td>
        <td>{{ unitMap[item.unit_id] }}</td>
        <td v-for="v in item.workHours" :key="v.process_id">
          <input type="number" v-model="v.work_hour">
        </td>
      </tr>
      </tbody>
    </table>
    <button type="submit" class="btn btn-primary" @click="postOperationSettingById(facilityId, list, $emit)">更新</button>
  </div>
</template>

<script setup lang="ts">
import {OperationSetting, Process, Unit, User} from "@/api";
import {postOperationSettingById} from "@/composable/operationSetting";

defineEmits(['openEditModal'])

interface AsyncOperationSettingTable {
  list: OperationSetting[]
  unitList: Unit[]
  processList: Process[]
  userList: User[]
  facilityId: number
}

const props = defineProps<AsyncOperationSettingTable>()
const unitMap: { [x: number]: string; } = {}
props.unitList.forEach(v => {
  unitMap[v.id!] = v.name
})
const userMap: { [x: number]: string; } = {}
props.userList.forEach(v => {
  userMap[v.id!] = v.name
})

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




