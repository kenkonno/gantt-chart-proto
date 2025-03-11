<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="unit.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">名称<input-required /></label>
      <input class="form-control" type="text" name="name" id="name" v-model="unit.name" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="unit.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="unit.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="validate(unit) && postUnit(unit, facilityId, order,$emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="validate(unit) && postUnitById(unit, facilityId, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteUnitById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useUnit, postUnitById, postUnit, deleteUnitById, validate} from "@/composable/unit";
import InputRequired from "@/components/form/InputRequired.vue";
interface AsyncUnitEdit {
  id: number | undefined,
  facilityId: number,
  order?: number
}

const props = defineProps<AsyncUnitEdit>()
defineEmits(['closeEditModal'])

const {unit} = await useUnit(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


