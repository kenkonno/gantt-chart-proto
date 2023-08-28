<template>
  <Suspense>
    <async-ticket-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="Ticket" @close-edit-modal="closeModalProxy">
      <async-ticket-edit :id="id" @close-edit-modal="closeModalProxy"></async-ticket-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncTicketTable from "@/components/ticket/AsyncTicketTable.vue";
import AsyncTicketEdit from "@/components/ticket/AsyncTicketEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useTicketTable} from "@/composable/ticket";
const {list, refresh} = await useTicketTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>