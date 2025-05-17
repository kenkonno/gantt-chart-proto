<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="facilityWorkSchedule.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">日付
        <input-required/>
      </label>
      <input class="form-control" type="date" name="date" id="date" v-model="facilityWorkSchedule.date"
             :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">稼動種別
        <input-required/>
      </label>
      <select class="form-control" name="type" id="type" v-model="facilityWorkSchedule.type" :disabled="false">
        <option v-for="(name, code) in FacilityWorkScheduleTypeMap" :value="code" :key="code">{{ name }}</option>
      </select>

    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="facilityWorkSchedule.created_at"
             :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="facilityWorkSchedule.updated_at"
             :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary"
              @click="postFacilityWorkSchedule(facilityWorkSchedule, facilityId, $emit)">更新
      </button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary"
              @click="postFacilityWorkScheduleById(facilityWorkSchedule, facilityId, $emit)">更新
      </button>
      <button type="submit" class="btn btn-warning" @click="deleteFacilityWorkScheduleById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {
  useFacilityWorkSchedule,
  postFacilityWorkScheduleById,
  postFacilityWorkSchedule,
  deleteFacilityWorkScheduleById
} from "@/composable/facilityWorkSchedule";
import {FacilityWorkScheduleTypeMap} from "@/const/common";
import InputRequired from "@/components/form/InputRequired.vue";

interface AsyncFacilityWorkScheduleEdit {
  id: number | undefined,
  facilityId: number,
}

const props = defineProps<AsyncFacilityWorkScheduleEdit>()
defineEmits(['closeEditModal'])

const {facilityWorkSchedule} = await useFacilityWorkSchedule(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


