<template>
  <input type="number"
         :value="modelValue"
         @input="$emit('update:modelValue',parseFloat($event.target.value) != null ? parseFloat($event.target.value) : undefined)"
         @change="$emit('change')"
         :disabled="disabled"
  />
</template>

<script setup lang="ts">
interface FormNumber {
  modelValue: number | undefined
  disabled?: boolean
}

const props = withDefaults(defineProps<FormNumber>(), {disabled: false, modelValue: undefined})
defineEmits(['update:modelValue', 'change'])


</script>