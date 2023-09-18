<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="ganttGroup.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">FacilityId</label>
      <input class="form-control" type="text" name="facilityId" id="facilityId" v-model="ganttGroup.facility_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">UnitId</label>
      <input class="form-control" type="text" name="unitId" id="unitId" v-model="ganttGroup.unit_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Order</label>
      <input class="form-control" type="text" name="order" id="order" v-model="ganttGroup.order" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="ganttGroup.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="ganttGroup.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postGanttGroup(ganttGroup, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postGanttGroupById(ganttGroup, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteGanttGroupById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useGanttGroup, postGanttGroupById, postGanttGroup, deleteGanttGroupById} from "@/composable/ganttGroup";

interface AsyncGanttGroupEdit {
  id: number | undefined
}

const props = defineProps<AsyncGanttGroupEdit>()
defineEmits(['closeEditModal'])

const {ganttGroup} = await useGanttGroup(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


