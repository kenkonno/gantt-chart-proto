<template>
  <input type="number"
         :value="modelValue"
         @input="$emit('update:modelValue',parseFloat($event.target.value) || undefined)"
         @change="$emit('change')"
  />
</template>

<script setup lang="ts">
interface FormNumber {
  modelValue: number | undefined
}
const props = defineProps<FormNumber>()
defineEmits(['update:modelValue', 'change'])


</script>