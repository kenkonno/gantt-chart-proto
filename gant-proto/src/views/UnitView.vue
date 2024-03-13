<template>
  <Suspense>
    <async-unit-table
        @open-edit-modal="openEditModal($event)"
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
</template>

<script setup lang="ts">
import AsyncUnitTable from "@/components/unit/AsyncUnitTable.vue";
import AsyncUnitEdit from "@/components/unit/AsyncUnitEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";

const {unitMap, currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {refreshUnitMap, updateUnitOrder} = inject(GLOBAL_ACTION_KEY)!
await refreshUnitMap(currentFacilityId)

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  await refreshUnitMap(currentFacilityId)
  closeEditModal()
  emit("update")
}

</script>