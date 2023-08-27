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

const {list, refresh} = await useOperationSettingTable(1)
const {list: unitList } = await useUnitTable()
const {list: processList } = await useProcessTable()
const {list: userList } = await useUserTable()

</script>