<template>
  <p @click="open=true">{{title}}</p>
  <DefaultModal v-if="open" @close-edit-modal="open=false" title="">
    <slot />
  </DefaultModal>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import DefaultModal from "@/components/modal/DefaultModal.vue";

interface ModalWithLink {
  title: string
}

const props = defineProps<ModalWithLink>()
defineEmits(['closeEditModal'])

const open = ref(false)
</script>

<style scoped>
p {
  text-decoration: underline;
  cursor: pointer;
  display: inline;
  padding: 5px;
}

</style>