<template>
  <Suspense>
    <async-holiday-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="Holiday" @close-edit-modal="closeModalProxy">
      <async-holiday-edit :id="id" :facility-id="facilityId" @close-edit-modal="closeModalProxy"></async-holiday-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncHolidayTable from "@/components/holiday/AsyncHolidayTable.vue";
import AsyncHolidayEdit from "@/components/holiday/AsyncHolidayEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useHolidayTable} from "@/composable/holiday";

interface HolidayView {
  facilityId: number
}
const props = defineProps<HolidayView>()
const {list, refresh} = await useHolidayTable(props.facilityId)
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])

const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
  emit("update")
}

</script>