<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>日付</th>
        <th>説明</th>
        <th>作成日</th>
        <th>更新日</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td @click="$emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ $filters.dateFormatYMD(item.date) }}</td>
        <td>{{ item.description }}</td>
        <td>{{ $filters.dateFormatYMD(item.created_at) }}</td>
        <td>{{ $filters.unixTimeFormat(item.updated_at) }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {Milestone} from "@/api";

defineEmits(['openEditModal'])

interface AsyncMilestoneTable {
  list: Milestone[]
}

const props = defineProps<AsyncMilestoneTable>()

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




