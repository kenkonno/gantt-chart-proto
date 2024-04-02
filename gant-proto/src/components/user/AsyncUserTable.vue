<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>部署</th>
        <th>氏名</th>
        <th v-if="false">稼働上限</th>
        <th>Role</th>
        <th v-if="false">Password</th>
        <th>Email</th>
        <th>作成日</th>
        <th>更新日</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td @click="$emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ getDepartmentName(item.department_id) }}</td>
        <td>{{ `${item.lastName} ${item.firstName}` }}</td>
        <td v-if="false">{{ item.limit_of_operation }}</td>
        <td>{{ RoleTypeMap[item.role] }}</td>
        <td v-if="false">{{ item.password }}</td>
        <td>{{ item.email }}</td>
        <td>{{ $filters.dateFormat(item.created_at) }}</td>
        <td>{{ $filters.unixTimeFormat(item.updated_at) }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {User} from "@/api";
import {inject} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {RoleTypeMap} from "../../const/common";

defineEmits(['openEditModal'])
const {departmentList} = inject(GLOBAL_STATE_KEY)!
const getDepartmentName = (id: number) => {
  return departmentList.find(v => v.id === id)?.name
}

interface AsyncUserTable {
  list: User[]
}

defineProps<AsyncUserTable>()

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




