<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="department.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">名称<input-required /></label>
      <input class="form-control" type="text" name="name" id="name" v-model="department.name" :disabled="false">
    </div>

    <div class="mb-2">
      <p class="form-label">色</p>
      <div><color-picker :disableAlpha="true" v-model:pure-color="department.color"/></div>
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="department.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="department.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="validate(department) && postDepartment(department, order, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="validate(department) && postDepartmentById(department, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteDepartmentById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {
  useDepartment,
  postDepartmentById,
  postDepartment,
  deleteDepartmentById,
  validate
} from "@/composable/department";
import InputRequired from "@/components/form/InputRequired.vue";
import {ColorPicker} from "vue3-colorpicker";

interface AsyncDepartmentEdit {
  id: number | undefined,
  order?: number
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


