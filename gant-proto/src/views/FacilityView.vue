<template>
  <Suspense>
    <async-project-list-name-sort-table v-if="available( 'ProjectListNameSort')"
                                        @open-edit-modal="openEditModal"
                                        @move-up="innerUpdateFacilityOrder($event, -1)"
                                        @move-down="innerUpdateFacilityOrder($event, 1)"
                                        :list="sortedFacilityList"
                                        :is-view-only="isViewOnly"
                                        :is-simulate="isSimulate"
    />
    <async-facility-table v-else
                          @open-edit-modal="openEditModal"
                          @move-up="innerUpdateFacilityOrder($event, -1)"
                          @move-down="innerUpdateFacilityOrder($event, 1)"
                          :list="sortedFacilityList"
                          :is-view-only="isViewOnly"
                          :is-simulate="isSimulate"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  <Suspense v-if="modalIsOpen">
    <default-modal title="案件" @close-edit-modal="closeModalProxy" :size="'half'">
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
import {computed, inject} from "vue";
import {useIsSimulate} from "@/composable/isSimulate";
import AsyncProjectListNameSortTable from "@/components/facility/featureOption/AsyncProjectListNameSortTable.vue";
import {available} from "@/composable/featureOption";

const {facilityList} = inject(GLOBAL_STATE_KEY)!
const {updateFacilityOrder} = inject(GLOBAL_ACTION_KEY)!

const {modalIsOpen, id, originalId, openEditModal, closeEditModal} = useModalWithId()
defineEmits(["update"])
const closeModalProxy = async () => {
  closeEditModal()
}
const sortedFacilityList = computed(() => {
  // TODO: 0がレスポンスから消えている
  return [...facilityList].sort((a, b) => (b.order ? b.order : 0) < (a.order ? a.order : 0) ? -1: 1);
});

// 降順で渡しているのでindexと方向を逆転させる
const innerUpdateFacilityOrder = (index: number, direction: number) => {
  // index は 0番目の時に最大、 最大の時に0
  // 方向は正と負を逆転させる。
  updateFacilityOrder((facilityList.length - 1) - index, direction * -1)
}

const { isViewOnly, isSimulate } = await useIsSimulate()


</script>