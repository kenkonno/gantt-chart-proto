<template>
  <Suspense>
    <AsyncOperationSettingTable
        :list="operationSettingMap[currentFacilityId]"
        :unitList="unitMap[currentFacilityId]"
        :processList="processList"
        :facility-id="currentFacilityId"
        @close-edit-modal="$emit('update')"
    />

    <template #fallback>
      Loading...
    </template>
  </Suspense>

</template>

<script setup lang="ts">
import AsyncOperationSettingTable from "@/components/operationSetting/AsyncOperationSettingTable.vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {inject} from "vue";

defineEmits(["update"])

const {processList, unitMap, operationSettingMap, currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {refreshUnitMap, refreshOperationSettingMap} = inject(GLOBAL_ACTION_KEY)!
await refreshUnitMap(currentFacilityId)
await refreshOperationSettingMap(currentFacilityId)

</script>