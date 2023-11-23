<template>
  <Suspense>
    <async-facility-table
        @open-edit-modal="openEditModal"
        @move-up="updateFacilityOrder($event, -1)"
        @move-down="updateFacilityOrder($event, 1)"
        :list="facilityList"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="設備" @close-edit-modal="closeModalProxy">
      <async-facility-edit :id="id" :order="facilityList.length + 1" :original-id="originalId"
                           @close-edit-modal="closeModalProxy"
                           @update="$emit('update')">

      </async-facility-edit>
    </default-modal>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
</template>

<script setup lang="ts">
import AsyncFacilityTable from "@/components/facility/AsyncFacilityTable.vue";
import AsyncFacilityEdit from "@/components/facility/AsyncFacilityEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {inject} from "vue";

const {facilityList} = inject(GLOBAL_STATE_KEY)!
const {updateFacilityOrder} = inject(GLOBAL_ACTION_KEY)!

const {modalIsOpen, id, originalId, openEditModal, closeEditModal} = useModalWithId()
defineEmits(["update"])
const closeModalProxy = async () => {
  closeEditModal()
}

</script>