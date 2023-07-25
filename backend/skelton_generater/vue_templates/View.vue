<template>
  <Suspense>
    <async-@Lower@-table @open-edit-modal="openEditModal($event)"></async-@Lower@-table>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="@Upper@" @close-edit-modal="closeEditModal">
      <async-@Lower@-edit :id="id" @close-edit-modal="closeEditModal"></async-@Lower@-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import Async@Upper@Table from "@/components/@Lower@/Async@Upper@Table.vue";
import Async@Upper@Edit from "@/components/@Lower@/Async@Upper@Edit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()

</script>