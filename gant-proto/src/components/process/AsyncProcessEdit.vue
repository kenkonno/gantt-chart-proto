<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="process.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Name</label>
      <input class="form-control" type="text" name="name" id="name" v-model="process.name" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">CreatedAt</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="process.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">UpdatedAt</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="process.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postProcess(process, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postProcessById(process, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteProcessById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useProcess, postProcessById, postProcess, deleteProcessById} from "@/composable/process";

interface AsyncProcessEdit {
  id: number | undefined
}

const props = defineProps<AsyncProcessEdit>()
defineEmits(['closeEditModal'])

const {process} = await useProcess(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


