<template>
  <Suspense>
    <async-milestone-table
        @open-edit-modal="openEditModal($event)"
        :list="list"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="マイルストーン" @close-edit-modal="closeModalProxy">
      <async-milestone-edit :id="id"
                            :order="list.length + 1"
                            @close-edit-modal="closeModalProxy"
      @update="updateGantt"></async-milestone-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncMilestoneTable from "@/components/milestone/AsyncMilestoneTable.vue";
import AsyncMilestoneEdit from "@/components/milestone/AsyncMilestoneEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {useMilestoneTable} from "@/composable/milestone";
import {inject} from "vue";
import {GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";

const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!

const {list, refresh} = await useMilestoneTable(currentFacilityId)
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
}
const updateGantt = () => {
  refreshGantt(currentFacilityId, false)
}

</script>