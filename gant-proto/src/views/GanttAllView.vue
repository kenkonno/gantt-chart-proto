<template>
  <gantt-all-menu
      :display-type="displayType"
      :gantt-all-header="GanttHeader"
      @update-display-type="updateDisplayType"
  >
  </gantt-all-menu>
  <Suspense v-if="globalState.ganttAllRefresh">
    <gantt-all
        :gantt-all-header="GanttHeader"
        :display-type="displayType">
    </gantt-all>
    <template #fallback>
      <DefaultSpinner></DefaultSpinner>
    </template>
  </Suspense>
</template>
<script setup lang="ts">
import {useGanttAllMenu} from "@/composable/ganttAllMenu";
import GanttAllMenu from "@/components/ganttAll/GanttAllMenu.vue";
import GanttAll from "@/components/ganttAll/GanttAll.vue";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import {inject} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import DefaultSpinner from "@/components/spinner/DefaultSpinner.vue";

const globalState = inject(GLOBAL_STATE_KEY)!

const {
  GanttHeader,
  displayType,
} = useGanttAllMenu()
const updateDisplayType = (v: DisplayType) => {
  displayType.value = v
}

</script>