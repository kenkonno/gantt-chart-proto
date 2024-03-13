<template>
  <Suspense>
    <async-process-table
        @open-edit-modal="openEditModal($event)"
        @move-up="updateProcessOrder($event, -1)"
        @move-down="updateProcessOrder($event, 1)"
        :list="processList"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="工程" @close-edit-modal="closeModalProxy" :size="'half'">
      <async-process-edit :id="id" :order="processList.length + 1" @close-edit-modal="closeModalProxy"></async-process-edit>
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
import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])

const {processList} = inject(GLOBAL_STATE_KEY)!
const {refreshProcessList,updateProcessOrder} = inject(GLOBAL_ACTION_KEY)!


const closeModalProxy = async () => {
  await refreshProcessList()
  closeEditModal()
  emit("update")
}

</script>