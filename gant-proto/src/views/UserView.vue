<template>
  <Suspense>
    <async-user-table
        @open-edit-modal="openEditModal($event)"
        :list="userList"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="担当者" @close-edit-modal="closeModalProxy">
      <async-user-edit :id="id" @close-edit-modal="closeModalProxy"></async-user-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncUserTable from "@/components/user/AsyncUserTable.vue";
import AsyncUserEdit from "@/components/user/AsyncUserEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";

const {userList} = inject(GLOBAL_STATE_KEY)!
const {refreshUserList, refreshDepartmentList} = inject(GLOBAL_ACTION_KEY)!
await refreshDepartmentList()

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  await refreshUserList()
  closeEditModal()
  emit("update")
}

</script>