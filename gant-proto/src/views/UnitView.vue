<template>
  <Suspense>
    <async-unit-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="Unit" @close-edit-modal="closeModalProxy">
      <async-unit-edit :id="id" @close-edit-modal="closeModalProxy"></async-unit-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncUnitTable from "@/components/unit/AsyncUnitTable.vue";
import AsyncUnitEdit from "@/components/unit/AsyncUnitEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useUnitTable} from "@/composable/unit";
const {list, refresh} = await useUnitTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>