<template>
  <Suspense>
    <async-ticketUser-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="担当者" @close-edit-modal="closeModalProxy">
      <async-ticketUser-edit :id="id" @close-edit-modal="closeModalProxy"></async-ticketUser-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncTicketUserTable from "@/components/ticketUser/AsyncTicketUserTable.vue";
import AsyncTicketUserEdit from "@/components/ticketUser/AsyncTicketUserEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useTicketUserTable} from "@/composable/ticketUser";
const {list, refresh} = await useTicketUserTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>