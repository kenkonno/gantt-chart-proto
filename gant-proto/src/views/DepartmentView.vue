<template>
  <Suspense>
    <async-department-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="Department" @close-edit-modal="closeModalProxy">
      <async-department-edit :id="id" @close-edit-modal="closeModalProxy"></async-department-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncDepartmentTable from "@/components/department/AsyncDepartmentTable.vue";
import AsyncDepartmentEdit from "@/components/department/AsyncDepartmentEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useDepartmentTable} from "@/composable/department";
const {list, refresh} = await useDepartmentTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}

</script>