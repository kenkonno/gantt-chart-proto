<template>
  <ModalsContainer></ModalsContainer>
  <a class="icon" @click="() => open()" style="border-bottom: 1px solid;">
    <span class="material-symbols-outlined">warning</span>
    <span>遅延通知</span>
    <span class="text-bg-danger">{{ scheduleAlert.length }}</span>
  </a>
</template>

<script setup lang="ts">

import {inject, onBeforeUnmount} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {GLOBAL_SCHEDULE_ALERT_KEY} from "@/composable/scheduleAlert";
import {ModalsContainer} from "vue-final-modal";

const {scheduleAlert} = inject(GLOBAL_STATE_KEY)!
const {getScheduleAlert} = inject(GLOBAL_ACTION_KEY)!
const {
  open, destroy
} = inject(GLOBAL_SCHEDULE_ALERT_KEY)!

defineEmits(["selectUnit"])


onBeforeUnmount(() => {
  destroy()
})
getScheduleAlert()

</script>

<style scoped lang="scss">

.material-symbols-outlined {
  font-size: 1.3rem;
}

a {
  display: block;
  margin-left: 5px;
  color: inherit;
  padding: 0;
  text-decoration: inherit;
  border-bottom: 1px solid black;
  cursor: pointer;
}

span {
  vertical-align: middle;
}

</style>


