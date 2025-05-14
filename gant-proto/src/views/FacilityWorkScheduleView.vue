<template>
  <Suspense>
    <async-facilityWorkSchedule-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="FacilityWorkSchedule" @close-edit-modal="closeModalProxy">
      <async-facilityWorkSchedule-edit :id="id" @close-edit-modal="closeModalProxy"
                                       :facility-id="currentFacilityId"></async-facilityWorkSchedule-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useFacilityWorkScheduleTable} from "@/composable/facilityWorkSchedule";
import AsyncFacilityWorkScheduleTable from "@/components/facilityWorkSchedule/AsyncFacilityWorkScheduleTable.vue";
import AsyncFacilityWorkScheduleEdit from "@/components/facilityWorkSchedule/AsyncFacilityWorkScheduleEdit.vue";
import {inject} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";

const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!

const {list, refresh} = await useFacilityWorkScheduleTable(currentFacilityId)
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>