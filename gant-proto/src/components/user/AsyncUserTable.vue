<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>名称</th>
        <th v-if="false">Password</th>
        <th>Email</th>
        <th>作成日</th>
        <th>更新日</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td @click="$emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ item.name }}</td>
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

defineEmits(['openEditModal'])

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




