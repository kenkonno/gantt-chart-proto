<template>
  <div>
    <label>
      案件状況：
      <label v-for="(name, code) in FacilityTypeMap" :key="code">
        <input
            type="checkbox"
            name="facilityType"
            :value="code"
            :checked="modelValue.includes(code)"
            @change="onCheckboxChange(code, $event)"
        />
        {{ name }}
      </label>
    </label>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
import {FacilityTypeMap} from "@/const/common";

// propsの定義
const props = defineProps<{
  modelValue: string[];
}>();

// emitの定義
const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void;
  (e: 'change', value: string[]): void;
}>();

// チェックボックスの変更イベントハンドラー
const onCheckboxChange = (code: string, event: Event) => {
  const isChecked = (event.target as HTMLInputElement).checked;
  let newValue: string[];

  if (isChecked) {
    // 選択された値が配列に存在しない場合は追加
    newValue = [...props.modelValue, code];
  } else {
    // 選択された値を配列から削除
    newValue = props.modelValue.filter(item => item !== code);
  }

  // v-modelの値を更新
  emit('update:modelValue', newValue);

  // 追加のイベントを発行（親コンポーネントのchangeFacilityType関数用）
  emit('change', newValue);
};
</script>

<style scoped>
/* 必要なスタイルがあればここに追加 */
</style>