<template>
  <VueFinalModal
      display-directive="show"
      background="interactive"
      content-transition="vfm-fade"
      :hide-overlay="true"
      @update:model-value="val => emit('update:modelValue', val)"
  >
    <!-- in case you use Nuxt, make sure to wrap with `<ClientOnly>` -->
    <!-- in case you don't use Nuxt, you don't need `<ClientOnly>` -->
    <VueDragResize
        :is-active="true"
        :x="500"
        :y="50"
        :w="600"
        :h="800"
        class="bg-primary-100 dark:bg-gray-800"
        :preventActiveBehavior="true"
        @resizing="dragResize"
        @dragging="dragResize"
    >
      <div class="wrapper" @mousedown="onSelectElement($event)" @touchstart="onSelectElement($event)">
        <div class="clearfix">
          <button type="button" class="btn-close float-end" data-bs-dismiss="modal" aria-label="Close"
                  @pointerup="emit('update:modelValue', false)"></button>
        </div>
        <div class="d-flex flex-column">
          <label>
            案件
            <select v-model="filterFacility">
              <option v-for="item in facilityList" :key="item.code" :value="item.code">{{ item.name }}</option>
            </select>
          </label>
          <label>
            遅延
            <input type="number" v-model="filterDelayDays" style="width: 3rem">
            日以上を表示
          </label>
          <label>
            進捗
            <input type="number" v-model="filterProgressNumber" style="width: 3rem">
            %以下を表示
          </label>
        </div>
        <div>
          <div v-for="[key, item] in Array.from(cScheduleAlert)" :key="key">
            <hr>
            <u class="link" @click="refreshGantt(item[0].facility_id, true)">{{ item[0].facility_name }}</u>
            <table style="width: 100%; table-layout: fixed">
              <thead>
              <tr>
                <th width="40%">ユニット名</th>
                <th width="40%">工程名</th>
                <th width="20%">進捗</th>
                <th width="20%">遅延日数</th>
              </tr>
              </thead>
              <tr v-for="i in item" :key="i.ticket_id">
                <td class="text-overflow">{{ i.unit_name }}</td>
                <td>{{ i.process_name }}</td>
                <td>{{ i.progress_percent }}%</td>
                <td>{{ i.delay_days }}</td>
              </tr>
            </table>
          </div>
        </div>
      </div>
    </VueDragResize>
  </VueFinalModal>
</template>

<script setup lang="ts">
import {VueFinalModal} from 'vue-final-modal'
import VueDragResize from 'vue3-drag-resize'
import {ref} from "vue";

const emit = defineEmits<{
  (e: 'update:modelValue', modelValue: boolean): void
}>()

const width = ref(0)
const height = ref(0)
const top = ref(0)
const left = ref(0)

function dragResize(newRect: any) {
  width.value = newRect.width
  height.value = newRect.height
  top.value = newRect.top
  left.value = newRect.left
}


import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_MUTATION_KEY} from "@/composable/globalState";
import {GLOBAL_SCHEDULE_ALERT_KEY} from "@/composable/scheduleAlert";

const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!
const {getScheduleAlert} = inject(GLOBAL_ACTION_KEY)!
const {
  facilityList,
  filterFacility,
  filterProgressNumber,
  filterDelayDays,
  cScheduleAlert,
} = inject(GLOBAL_SCHEDULE_ALERT_KEY)!

const onSelectElement = function (event: any) {
  const tagName = event.target.tagName

  if (tagName === 'INPUT' || tagName === 'SELECT') {
    event.stopPropagation()
  }
}
getScheduleAlert()
</script>
<style>
.link {
  cursor: pointer;
  text-decoration: underline;
}

p {
  margin-right: 5px;
  cursor: pointer;
  display: inline;
  padding: 5px 0;
}

.text-overflow {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.content-container {
  overflow: scroll;
  background: white;
}
</style>