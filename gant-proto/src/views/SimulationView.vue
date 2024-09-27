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
        <td>{{ simulationLock.lockedBy }}</td>
        <td>{{ simulationLock.lockedAt }}</td>
      </tr>
      </tbody>
    </table>
  </div>
  <hr>
  <div class="d-flex justify-content-between">
    <div>
      <button type="submit" class="btn btn-primary" @click="postSimulation($emit)">開始</button>
    </div>
    <div>
      <button type="submit" class="btn btn-secondary" @click="putSimulation('pending',$emit)">中断</button>
    </div>
    <div>
      <button type="submit" class="btn btn-info" @click="putSimulation('resume',$emit)">再開</button>
    </div>
    <div>
      <button type="submit" class="btn btn-success" @click="putSimulation('apply',$emit)">反映</button>
    </div>
    <div>
      <button type="submit" class="btn btn-danger" @click="deleteSimulation($emit)">破棄</button>
    </div>
  </div>
</template>

<script setup lang="ts">

import {deleteSimulation, postSimulation, putSimulation, useSimulation} from "@/composable/simulation";

const {simulationLock} = await useSimulation()

</script>
