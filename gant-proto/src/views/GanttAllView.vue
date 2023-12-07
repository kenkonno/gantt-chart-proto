<template>
  <gantt-all-menu
      :display-type="displayType"
      :gantt-all-header="GanttHeader"
      @update-display-type="updateDisplayType"
  >
  </gantt-all-menu>
  <Suspense>
    <gantt-all
        v-if="globalState.ganttAllRefresh"
        :gantt-all-header="GanttHeader"
        :display-type="displayType">
    </gantt-all>
    <template #fallback>
      Loading...
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
const globalState = inject(GLOBAL_STATE_KEY)!

const {
  GanttHeader,
  displayType,
} = useGanttAllMenu()
const updateDisplayType = (v: DisplayType) => {
  console.log("################ displayType", v)
  displayType.value = v
}

</script>