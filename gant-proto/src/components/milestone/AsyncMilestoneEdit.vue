<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="milestone.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">日付<input-required /></label>
      <input class="form-control" type="date" name="date" id="date" v-model="milestone.date" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">説明<input-required /></label>
      <input class="form-control" type="text" name="description" id="description" v-model="milestone.description"
             :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="milestone.created_at"
             :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="milestone.updated_at"
             :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="validate(milestone) && postMilestone(milestone, order, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="validate(milestone) && postMilestoneById(milestone, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteMilestoneById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useMilestone, postMilestoneById, postMilestone, deleteMilestoneById, validate} from "@/composable/milestone";

interface AsyncMilestoneEdit {
  id: number | undefined
  order?: number
}

import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {inject} from "vue";
import InputRequired from "@/components/form/InputRequired.vue";

const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!

const props = defineProps<AsyncMilestoneEdit>()
defineEmits(['closeEditModal'])

const {milestone} = await useMilestone(currentFacilityId, props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


