<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">名称<input-required /></label>
      <input class="form-control" type="text" name="name" id="name" v-model="newUnitName" :disabled="false">
    </div>

    <div class="copy-options mb-3">
      <h5>工程のコピー設定</h5>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyProcessId" v-model="copyOptions.copyProcessId">
        <label class="form-check-label" for="copyProcessId">工程</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyDepartmentId" v-model="copyOptions.copyDepartmentId">
        <label class="form-check-label" for="copyDepartmentId">部署</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyTicketUser" v-model="copyOptions.copyTicketUser">
        <label class="form-check-label" for="copyTicketUser">担当者</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyNumberOfWorker" v-model="copyOptions.copyNumberOfWorker">
        <label class="form-check-label" for="copyNumberOfWorker">人数</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyEstimate" v-model="copyOptions.copyEstimate">
        <label class="form-check-label" for="copyEstimate">工数</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyDaysAfter" v-model="copyOptions.copyDaysAfter">
        <label class="form-check-label" for="copyDaysAfter">日後</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyStartDate" v-model="copyOptions.copyStartDate">
        <label class="form-check-label" for="copyStartDate">開始日</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyEndDate" v-model="copyOptions.copyEndDate">
        <label class="form-check-label" for="copyEndDate">終了日</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyProgressPercent" v-model="copyOptions.copyProgressPercent">
        <label class="form-check-label" for="copyProgressPercent">進捗</label>
      </div>
      <div class="form-check mb-2">
        <input class="form-check-input" type="checkbox" id="copyMemo" v-model="copyOptions.copyMemo">
        <label class="form-check-label" for="copyMemo">メモ</label>
      </div>
    </div>

    <button type="submit" class="btn btn-primary" @click="validateAndDuplicate">複製</button>
    <button type="button" class="btn btn-secondary ms-2" @click="$emit('closeEditModal')">キャンセル</button>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useUnit } from "@/composable/unit";
import InputRequired from "@/components/form/InputRequired.vue";
import {Api} from "@/api/axios";
import {toast} from "vue3-toastify";

interface AsyncUnitDuplicateEdit {
  id: number | undefined,
  facilityId: number,
  order?: number
}

const props = defineProps<AsyncUnitDuplicateEdit>()
const emit = defineEmits(['closeEditModal', 'duplicateSuccess'])

const { unit } = await useUnit(props.id)
const newUnitName = ref(unit.value ? `${unit.value.name} (コピー)` : '')

// コピーオプションの状態を管理
const copyOptions = reactive({
  copyProcessId: true,
  copyDepartmentId: true,
  copyEstimate: true,
  copyNumberOfWorker: true,
  copyDaysAfter: true,
  copyStartDate: true,
  copyEndDate: true,
  copyProgressPercent: true,
  copyTicketUser: true,
  copyMemo: false,
})

const validateAndDuplicate = async () => {
  if (!unit.value || !validateName()) {
    return
  }

  try {
    const response = await Api.postUnitsDuplicate({
      unitId: unit.value.id!,
      unitName: newUnitName.value,
      ...copyOptions
    })

    if (response.status === 200) {
      toast("成功しました。")
      emit('duplicateSuccess')
      emit('closeEditModal')
    }
  } catch (error) {
    console.error('ユニットの複製に失敗しました', error)
  }
}

const validateName = () => {
  if (!newUnitName.value || newUnitName.value.trim() === '') {
    alert('名称は必須です')
    return false
  }
  return true
}
</script>

<style scoped lang="scss">
.form-check label {
  float: none;
  margin-left: 5px;
}

.copy-options {
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
  padding: 15px;
  background-color: #f8f9fa;
}
</style>