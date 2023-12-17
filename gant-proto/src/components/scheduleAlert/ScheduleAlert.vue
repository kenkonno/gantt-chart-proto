<template>
  <a class="icon" @click="tableIsOpen = !tableIsOpen" style="border-bottom: 1px solid;">
    <span class="material-symbols-outlined">warning</span>
    <span>遅延通知</span>
    <span class="text-bg-danger">{{ scheduleAlert.length }}</span>
  </a>
  <div class="schedule-alert-table" v-if="tableIsOpen">
    <div class="clearfix">
      <button type="button" class="btn-close float-end" data-bs-dismiss="modal" aria-label="Close"
              @click="tableIsOpen = false"></button>
    </div>
    <div class="d-flex flex-column">
      <label>
        設備
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
        <u class="link" @click="refreshGantt(item[0].facility_id)">{{ item[0].facility_name }}</u>
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
</template>

<script setup lang="ts">

import {inject} from "vue";
import {GLOBAL_ACTION_KEY, GLOBAL_MUTATION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {GLOBAL_SCHEDULE_ALERT_KEY} from "@/composable/scheduleAlert";

const {scheduleAlert} = inject(GLOBAL_STATE_KEY)!
const {refreshGantt} = inject(GLOBAL_MUTATION_KEY)!
const {getScheduleAlert} = inject(GLOBAL_ACTION_KEY)!
const {
  facilityList,
  filterFacility,
  filterProgressNumber,
  filterDelayDays,
  cScheduleAlert,
  tableIsOpen
} = inject(GLOBAL_SCHEDULE_ALERT_KEY)!

defineEmits(["selectUnit"])


getScheduleAlert()

</script>

<style scoped lang="scss">
.link {
  cursor: pointer;
  text-decoration: underline;
}

.schedule-alert-container {
  z-index: 9999;
}

.table {
  z-index: 9999;
}

.schedule-alert-table {
  position: fixed;
  top: 0;
  left: 70%;
  width: 30%;
  z-index: 201;
  height: 100%;
  overflow: scroll;
  background: white;
  border: 1px solid black;
}

p {
  margin-right: 5px;
  cursor: pointer;
  display: inline;
  padding: 5px 0;
}

span {
  vertical-align: middle;
}

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

.text-overflow {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}
</style>


