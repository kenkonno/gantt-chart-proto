<template>
  <Suspense>
    <async-unit-table
        @open-edit-modal="openEditModal($event)"
        @open-copy-edit-modal="openCopyEditModal($event)"
        @move-up="updateUnitOrder($event, -1)"
        @move-down="updateUnitOrder($event, 1)"
        :list="unitMap[currentFacilityId]"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="ユニット" @close-edit-modal="closeModalProxy" :size="'half'">
      <async-unit-edit :id="id" :order="unitMap[currentFacilityId].length + 1" :facility-id="currentFacilityId"
                       @close-edit-modal="closeModalProxy"></async-unit-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="copyModalIsOpen">
    <default-modal title="ユニットのコピー" @close-edit-modal="closeModalProxy" :size="'half'">
      <async-unit-duplicate-edit :id="copyId" :order="unitMap[currentFacilityId].length + 1"
                                 :facility-id="currentFacilityId"
                                 @close-edit-modal="closeModalProxy"></async-unit-duplicate-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncUnitTable from "@/components/unit/AsyncUnitTable.vue";
import AsyncUnitEdit from "@/components/unit/AsyncUnitEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import AsyncUnitDuplicateEdit from "@/components/unit/AsyncUnitDuplicateEdit.vue";

const {unitMap, currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {refreshUnitMap, updateUnitOrder} = inject(GLOBAL_ACTION_KEY)!
const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!
await refreshUnitMap(currentFacilityId)

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const {
  modalIsOpen: copyModalIsOpen,
  id: copyId,
  openEditModal: openCopyEditModal,
  closeEditModal: closeCopyEditModal
} = useModalWithId()
const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  await refreshUnitMap(currentFacilityId)
  // NOTE: Unitの場合はガントチャートの更新が必要。
  refreshGantt(currentFacilityId, false)
  closeEditModal()
  closeCopyEditModal()
  emit("update")
}

</script>