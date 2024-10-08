<template>
  <div class="container" style="min-height: 100px">
    <table class="table">
      <thead>
      <tr>
        <th>ステータス</th>
        <th>作成者</th>
        <th>作成日</th>
      </tr>
      </thead>
      <tbody>
      <tr v-if="simulationLock.simulationName != ''">
        <td>{{ simulationLock.status }}</td>
        <td>{{ getUserName(simulationLock.lockedBy) }}</td>
        <td>{{ simulationLock.lockedAt }}</td>
      </tr>
      </tbody>
    </table>
  </div>
  <hr>
  <MasterDiffTables v-if="simulationLock.status == 'in_progress'"></MasterDiffTables>
  <div class="d-flex justify-content-between">
    <div>
      <button type="submit" class="btn btn-primary"
              :disabled="startDisabled"
              @click="postSimulation($emit)">開始
      </button>
    </div>
    <div>
      <button type="submit" class="btn btn-secondary"
              :disabled="pendingDisabled"
              @click="putSimulation('pending',$emit); refresh()">保留
      </button>
    </div>
    <div>
      <button type="submit" class="btn btn-info"
              :disabled="resumeDisabled"
              @click="putSimulation('resume',$emit); refresh()">再開
      </button>
    </div>
    <div>
      <button type="submit" class="btn btn-success"
              :disabled="applyDisabled"
              @click="putSimulation('apply',$emit); refresh()">反映
      </button>
    </div>
    <div>
      <button type="submit" class="btn btn-danger"
              :disabled="deleteDisabled"
              @click="deleteSimulation($emit); refresh()">破棄
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">

import {deleteSimulation, postSimulation, putSimulation, useSimulation} from "@/composable/simulation";
import {computed, inject} from "vue";
import MasterDiffTables from "@/components/masterDiff/MasterDiffTables.vue";
import {getUserInfo} from "@/composable/auth";
import {allowed} from "@/composable/role";
import {GLOBAL_GETTER_KEY} from "@/composable/globalState";

const {simulationLock, refresh} = await useSimulation()
const {getUserName} = inject(GLOBAL_GETTER_KEY)!

const userInfo = getUserInfo()

const isSimulateUser = () => {
  return userInfo.userInfo?.id == simulationLock.value.lockedBy
}

const startDisabled = computed(() => {
  // 開始はロックがないときだけ
  return !(simulationLock.value.status === '')
})
const pendingDisabled = computed(() => {
  // 保留は開始中だけ
  if (!allowed('FORCE_SIMULATE_USER') && !isSimulateUser()) return true
  return !(simulationLock.value.status === 'in_progress')
})
const resumeDisabled = computed(() => {
  // 再開は保留中だけ
  if (!allowed('FORCE_SIMULATE_USER') && !isSimulateUser()) return true
  return !(simulationLock.value.status === 'in_pending')
})
const applyDisabled = computed(() => {
  if (!allowed('FORCE_SIMULATE_USER') && !isSimulateUser()) return true
  return !(simulationLock.value.status === 'in_progress')
})
const deleteDisabled = computed(() => {
  // 破棄は何かあるときにだけ押せる。
  if (!allowed('FORCE_SIMULATE_USER') && !isSimulateUser()) return true
  return !(simulationLock.value.status !== '')
})

</script>
