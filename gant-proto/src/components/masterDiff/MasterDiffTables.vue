<template>
  <H5>マスター差分</H5>
  <small>※オリジナルのみに存在。差分有り（オリジナル・シミュレーションの準に上下に表示）。シミュレーションのみに存在の順で表示しています。</small>
  <hr>
  <div
      v-if="masterDiffList.beforeProcesses.length > 0 ||masterDiffList.diffProcesses?.length > 0 || masterDiffList.afterProcesses.length > 0">
    <h5>工程</h5>
    <AsyncProcessTable :list="masterDiffList.beforeProcesses" :is-view-only="true"/>
    <AsyncProcessTable :list="masterDiffList.diffProcesses" :is-view-only="true" :is-no-header="true"/>
    <AsyncProcessTable :list="masterDiffList.afterProcesses" :is-view-only="true" :is-no-header="true"/>
  </div>

  <div
      v-if="masterDiffList.beforeDepartments.length > 0 ||masterDiffList.diffDepartments?.length > 0 || masterDiffList.afterDepartments.length > 0">
    <h5>部署</h5>
    <AsyncDepartmentTable ::list="masterDiffList.beforeDepartments" :is-view-only="true"/>
    <AsyncDepartmentTable ::list="masterDiffList.diffDepartments" :is-view-only="true" :is-no-header="true"/>
    <AsyncDepartmentTable ::list="masterDiffList.afterDepartments" :is-view-only="true" :is-no-header="true"/>
  </div>

  <div
      v-if="masterDiffList.beforeUsers.length > 0 ||masterDiffList.diffUsers?.length > 0 || masterDiffList.afterUsers.length > 0">
    <h5>担当者</h5>
    <AsyncUserTable :list="masterDiffList.beforeUsers" :is-view-only="true"/>
    <AsyncUserTable :list="masterDiffList.diffUsers" :is-view-only="true" :is-no-header="true"/>
    <AsyncUserTable :list="masterDiffList.afterUsers" :is-view-only="true" :is-no-header="true"/>
  </div>

  <div
      v-if="masterDiffList.beforeUnits.length > 0 ||masterDiffList.diffUnits?.length > 0 || masterDiffList.afterUnits.length > 0">
    <h5>ユニット</h5>
    <AsyncUnitTable :list="masterDiffList.beforeUnits" :is-view-only="true"/>
    <AsyncUnitTable :list="masterDiffList.diffUnits" :is-view-only="true" :is-no-header="true"/>
    <AsyncUnitTable :list="masterDiffList.afterUnits" :is-view-only="true" :is-no-header="true"/>
  </div>

  <div
      v-if="masterDiffList.beforeHolidays.length > 0 ||masterDiffList.diffHolidays?.length > 0 || masterDiffList.afterHolidays.length > 0">
    <h5>祝日</h5>
    <AsyncHolidayTable :list="masterDiffList.beforeHolidays" :is-view-only="true"/>
    <AsyncHolidayTable :list="masterDiffList.diffHolidays" :is-view-only="true" :is-no-header="true"/>
    <AsyncHolidayTable :list="masterDiffList.afterHolidays" :is-view-only="true" :is-no-header="true"/>
  </div>

  <div
      v-if="masterDiffList.beforeMilestones.length > 0 ||masterDiffList.diffMilestones?.length > 0 || masterDiffList.afterMilestones.length > 0">
    <h5>マイルストーン</h5>
    <AsyncMilestoneTable :list="masterDiffList.beforeMilestones" :is-view-only="true"/>
    <AsyncMilestoneTable :list="masterDiffList.diffMilestones" :is-view-only="true" :is-no-header="true"/>
    <AsyncMilestoneTable :list="masterDiffList.afterMilestones" :is-view-only="true" :is-no-header="true"/>
  </div>
</template>

<script setup lang="ts">
import AsyncProcessTable from "@/components/process/AsyncProcessTable.vue";
import AsyncUnitTable from "@/components/unit/AsyncUnitTable.vue";
import AsyncUserTable from "@/components/user/AsyncUserTable.vue";
import AsyncHolidayTable from "@/components/holiday/AsyncHolidayTable.vue";
import AsyncMilestoneTable from "@/components/milestone/AsyncMilestoneTable.vue";
import AsyncDepartmentTable from "@/components/department/AsyncDepartmentTable.vue";
import {Api} from "@/api/axios";

const {data: masterDiffList} = await Api.getSimulationMasterDiff()


</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


