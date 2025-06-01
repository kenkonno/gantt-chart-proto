<template>
  <Suspense>
    <async-holiday-table
        @open-edit-modal="openEditModal($event)"
        :list="holidayList"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="祝日" @close-edit-modal="closeModalProxy" :size="'half'">
      <async-holiday-edit :id="id" @close-edit-modal="closeModalProxy"></async-holiday-edit>
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
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {inject} from "vue";

const {holidayList} = inject(GLOBAL_STATE_KEY)!
const {refreshHolidayList} = inject(GLOBAL_ACTION_KEY)!
await refreshHolidayList()

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])

const closeModalProxy = async () => {
  await refreshHolidayList()
  closeEditModal()
  emit("update")
}

</script>
