<template>
  <div class="wrapper d-flex">
    <div class="toggle-icon" @click="show = !show">
      <slot name="icon"/>
    </div>
    <Transition>
      <div class="accordion-body" v-if="show">
        <slot name="body"/>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
const show = ref(false)
</script>

<style scoped>
.toggle-icon {
  cursor: pointer;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease, width 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>