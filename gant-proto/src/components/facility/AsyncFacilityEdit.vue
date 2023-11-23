<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="facility.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">名称</label>
      <input class="form-control" type="text" name="name" id="name" v-model="facility.name" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">開始日</label>
      <input class="form-control" type="date" name="termFrom" id="termFrom" v-model="facility.term_from"
             :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">終了日</label>
      <input class="form-control" type="date" name="termTo" id="termTo" v-model="facility.term_to" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="facility.created_at"
             :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="facility.updated_at"
             :disabled="true">
    </div>

    <template v-if="id != null && originalId != null">
      <button type="submit" class="btn btn-primary" @click="copyFacility(facility, order, originalId, emit)">コピー</button>
    </template>
    <template v-else-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postFacility(facility, order, emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postFacilityById(facility, emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteFacilityById(id, emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useFacility, postFacilityById, postFacility, copyFacility,deleteFacilityById} from "@/composable/facility";

interface AsyncFacilityEdit {
  id: number | undefined,
  originalId: number | undefined,
  order?: number
}

const props = defineProps<AsyncFacilityEdit>()
const emit = defineEmits(['closeEditModal', 'update'])

const {facility} = await useFacility(props.id)

if(props.originalId != null) {
  facility.value.name = facility.value.name + "のコピー"
  facility.value.created_at = ""
  facility.value.updated_at = 0
}

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


