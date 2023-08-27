<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="unit.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Name</label>
      <input class="form-control" type="text" name="name" id="name" v-model="unit.name" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">CreatedAt</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="unit.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">UpdatedAt</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="unit.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postUnit(unit, facilityId, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postUnitById(unit, facilityId, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteUnitById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useUnit, postUnitById, postUnit, deleteUnitById} from "@/composable/unit";

interface AsyncUnitEdit {
  id: number | undefined,
  facilityId: number
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


