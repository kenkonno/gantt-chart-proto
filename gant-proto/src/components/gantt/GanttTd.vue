<template>
  <td class="side-menu-cell" :class="{'d-none': !visible}">
    <slot></slot>
  </td>
</template>

<script setup lang="ts">
interface GanttTd {
  visible: boolean
}
defineProps<GanttTd>()

</script>