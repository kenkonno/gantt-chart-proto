<template>
  <Suspense>
    <async-user-table @open-edit-modal="openEditModal($event)"></async-user-table>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="User" @close-edit-modal="closeEditModal">
      <async-user-edit :id="id" @close-edit-modal="closeEditModal"></async-user-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncUserTable from "@/components/user/AsyncUserTable.vue";
import AsyncUserEdit from "@/components/user/AsyncUserEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()

</script>