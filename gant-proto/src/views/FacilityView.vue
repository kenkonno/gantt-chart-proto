<template>
  <Suspense>
    <async-facility-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="Facility" @close-edit-modal="closeModalProxy">
      <async-facility-edit :id="id" @close-edit-modal="closeModalProxy"></async-facility-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncFacilityTable from "@/components/facility/AsyncFacilityTable.vue";
import AsyncFacilityEdit from "@/components/facility/AsyncFacilityEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useFacilityTable} from "@/composable/facility";
const {list, refresh} = await useFacilityTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>