<template>
  <Suspense>
    <async-process-table
        @open-edit-modal="openEditModal($event)"
        @move-up="updateOrder($event, -1)"
        @move-down="updateOrder($event, 1)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="Process" @close-edit-modal="closeModalProxy">
      <async-process-edit :id="id" :order="list.length + 1" @close-edit-modal="closeModalProxy"></async-process-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncProcessTable from "@/components/process/AsyncProcessTable.vue";
import AsyncProcessEdit from "@/components/process/AsyncProcessEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useProcessTable} from "@/composable/process";
const {list, refresh, updateOrder} = await useProcessTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])

const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
  emit("update")
}

</script>