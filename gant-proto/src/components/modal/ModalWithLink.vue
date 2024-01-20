<template>
  <p @click="open=true" v-if="!disabled" style="border-bottom: 1px solid;">
    <span class="material-symbols-outlined" v-if="icon != null">{{ icon }}</span>
    <span>{{ title }}</span>
  </p>
  <p v-else style="text-decoration: inherit; cursor:inherit">
    <span class="material-symbols-outlined" v-if="icon != null">{{ icon }}</span>
    <span>{{ title }}</span>
  </p>
  <DefaultModal :title="title" v-if="open" @close-edit-modal="open=false">
    <slot/>
  </DefaultModal>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import DefaultModal from "@/components/modal/DefaultModal.vue";

interface ModalWithLink {
  title: string,
  disabled?: boolean
  icon?: string
}

withDefaults(defineProps<ModalWithLink>(), {disabled: false})
defineEmits(['closeEditModal'])

const open = ref(false)
</script>

<style scoped>
p {
  margin-right: 5px;
  cursor: pointer;
  display: inline;
  padding: 5px 0;
}
span {
  vertical-align: middle;
}
.material-symbols-outlined {
  font-size: 1.3rem;
}

</style>