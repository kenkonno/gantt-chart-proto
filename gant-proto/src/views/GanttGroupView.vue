<template>
  <Suspense>
    <async-ganttGroup-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="GanttGroup" @close-edit-modal="closeModalProxy">
      <async-ganttGroup-edit :id="id" @close-edit-modal="closeModalProxy"></async-ganttGroup-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncGanttGroupTable from "@/components/ganttGroup/AsyncGanttGroupTable.vue";
import AsyncGanttGroupEdit from "@/components/ganttGroup/AsyncGanttGroupEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useGanttGroupTable} from "@/composable/ganttGroup";
const {list, refresh} = await useGanttGroupTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>