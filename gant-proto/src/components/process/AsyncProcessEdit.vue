<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="process.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">名称<input-required /></label>
      <input class="form-control" type="text" name="name" id="name" v-model="process.name" :disabled="false">
    </div>

    <div class="mb-2">
      <p class="form-label">色</p>
      <div><color-picker :disableAlpha="true" v-model:pure-color="process.color"/></div>
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="process.created_at"
             :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="process.updated_at"
             :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="validate(process) && postProcess(process, order, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="validate(process) && postProcessById(process, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteProcessById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {ColorPicker} from "vue3-colorpicker";
import {useProcess, postProcessById, postProcess, deleteProcessById, validate} from "@/composable/process";
import InputRequired from "@/components/form/InputRequired.vue";

interface AsyncProcessEdit {
  id: number | undefined,
  order?: number
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


