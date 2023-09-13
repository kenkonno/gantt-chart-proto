<template>
  <Suspense>
    <async-operationSetting-table
        :list="list"
        :unitList="unitList"
        :processList="processList"
        :facility-id="facilityId"
        @close-edit-modal="$emit('update')"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>

</template>

<script setup lang="ts">
import {useOperationSettingTable} from "@/composable/operationSetting";
import {useUnitTable} from "@/composable/unit";
import {useProcessTable} from "@/composable/process";

interface OperationSettingView {
  facilityId: number
}
const props = defineProps<OperationSettingView>()
const emit = defineEmits(["update"])

const {list, refresh} = await useOperationSettingTable(props.facilityId)
const {list: unitList } = await useUnitTable(props.facilityId)
const {list: processList } = await useProcessTable()

</script>