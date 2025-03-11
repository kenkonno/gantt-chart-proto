<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="holiday.id" :disabled="true">
    </div>

    <div class="mb-2" v-if="false">
      <label class="form-label" for="id">名称<input-required /></label>
      <input class="form-control" type="text" name="name" id="name" v-model="holiday.name" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">日付<input-required /></label>
      <input class="form-control" type="date" name="date" id="date" v-model="holiday.date" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="holiday.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="holiday.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="validate(holiday) && postHoliday(holiday, facilityId, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="validate(holiday) && postHolidayById(holiday, facilityId, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteHolidayById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useHoliday, postHolidayById, postHoliday, deleteHolidayById, validate} from "@/composable/holiday";
import InputRequired from "@/components/form/InputRequired.vue";

interface AsyncHolidayEdit {
  id: number | undefined,
  facilityId: number
}

const props = defineProps<AsyncHolidayEdit>()
defineEmits(['closeEditModal'])

const {holiday} = await useHoliday(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


