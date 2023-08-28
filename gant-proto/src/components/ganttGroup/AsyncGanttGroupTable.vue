<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>FacilityId</th>
        <th>UnitId</th>
        <th>Order</th>
        <th>CreatedAt</th>
        <th>UpdatedAt</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td @click="$emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ item.facility_id }}</td>
        <td>{{ item.unit_id }}</td>
        <td>{{ item.order }}</td>
        <td>{{ item.created_at }}</td>
        <td>{{ item.updated_at }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {GanttGroup} from "@/api";

defineEmits(['openEditModal'])

interface AsyncGanttGroupTable {
  list: GanttGroup[]
}

const props = defineProps<AsyncGanttGroupTable>()

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




