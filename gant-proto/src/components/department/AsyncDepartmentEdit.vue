<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="department.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Name</label>
      <input class="form-control" type="text" name="name" id="name" v-model="department.name" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">CreatedAt</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="department.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">UpdatedAt</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="department.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postDepartment(department, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postDepartmentById(department, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteDepartmentById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useDepartment, postDepartmentById, postDepartment, deleteDepartmentById} from "@/composable/department";

interface AsyncDepartmentEdit {
  id: number | undefined
}

const props = defineProps<AsyncDepartmentEdit>()
defineEmits(['closeEditModal'])

const {department} = await useDepartment(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


