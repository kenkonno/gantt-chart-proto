<template>
  <Suspense>
    <async-department-table
        @open-edit-modal="openEditModal($event)"
        @move-up="updateDepartmentOrder($event, -1)"
        @move-down="updateDepartmentOrder($event, 1)"
        :list="departmentList"
        :is-view-only="isViewOnly"
        :is-simulate="isSimulate"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="部署" @close-edit-modal="closeModalProxy" :size="'half'">
      <async-department-edit :id="id" :order="departmentList.length + 1"
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
import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {Api} from "@/api/axios";
import AsyncProcessTable from "@/components/process/AsyncProcessTable.vue";
import {useIsSimulate} from "@/composable/isSimulate";

const {departmentList} = inject(GLOBAL_STATE_KEY)!
const {refreshDepartmentList, updateDepartmentOrder} = inject(GLOBAL_ACTION_KEY)!
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  await refreshDepartmentList()
  closeEditModal()
  emit("update")
}

const { isViewOnly, isSimulate } = await useIsSimulate()

</script>