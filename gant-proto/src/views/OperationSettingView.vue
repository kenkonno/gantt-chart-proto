<template>
  <Suspense>
    <async-operationSetting-table
        :list="list"
        :unitList="unitList"
        :processList="processList"
        :userList="userList"
    />
    <template #fallback>
      Loading...
    </template>
  </Suspense>

</template>

<script setup lang="ts">
import AsyncOperationSettingTable from "@/components/operationSetting/AsyncOperationSettingTable.vue";
import {useOperationSettingTable} from "@/composable/operationSetting";
import {useUnitTable} from "@/composable/unit";
import {useProcessTable} from "@/composable/process";
import {useUserTable} from "@/composable/user";

interface OperationSettingView {
  facilityId: number
}
const props = defineProps<OperationSettingView>()


const {list, refresh} = await useOperationSettingTable(props.facilityId)
const {list: unitList } = await useUnitTable(props.facilityId)
const {list: processList } = await useProcessTable()
const {list: userList } = await useUserTable()

</script>