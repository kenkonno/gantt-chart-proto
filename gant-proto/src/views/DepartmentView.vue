<template>
  <Suspense>
    <async-department-table
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
    <default-modal title="部署" @close-edit-modal="closeModalProxy">
      <async-department-edit :id="id" :order="list.length + 1"
                             @close-edit-modal="closeModalProxy"></async-department-edit>
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
import {useFacilityTable} from "@/composable/facility";

const {list, refresh, updateOrder} = await useDepartmentTable()
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
  emit("update")
}

</script>