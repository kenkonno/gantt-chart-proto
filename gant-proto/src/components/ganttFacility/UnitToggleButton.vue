<template>
    <span class="material-symbols-outlined align-middle pointer" @click="$emit('toggle')">
      {{ isOpen ? "remove" : "add" }}
    </span>
</template>
<script setup lang="ts">
type UnitToggleButtonProps = {
  isOpen: boolean
}

defineProps<UnitToggleButtonProps>()
defineEmits(['toggle'])
</script>
<style lang="scss" scoped>
.pointer {
  cursor: pointer;
}
</style>