<template>
  <Suspense>
    <async-unit-table
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
    <default-modal title="Unit" @close-edit-modal="closeModalProxy">
      <async-unit-edit :id="id" :order="list.length + 1" :facility-id="facilityId" @close-edit-modal="closeModalProxy"></async-unit-edit>
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
import {useUnitTable} from "@/composable/unit";
interface UnitView {
  facilityId: number
}
const props = defineProps<UnitView>()

const {list, refresh, updateOrder} = await useUnitTable(props.facilityId)
const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()
const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  await refresh()
  closeEditModal()
  emit("update")
}

</script>