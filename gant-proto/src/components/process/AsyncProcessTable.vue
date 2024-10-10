<template>
  <div class="container">
    <div v-if="isSimulate" class="mb-2 bg-warning text-center">シミュレーション中のため変更できません。</div>
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)" v-if="!isViewOnly">新規追加</button>
    <table class="table">
      <thead v-if="!isNoHeader">
      <tr>
        <th>Id</th>
        <th>名称</th>
        <th>色</th>
        <th>作成日</th>
        <th>更新日</th>
        <th v-if="!isViewOnly">並び替え</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(item, index) in list" :key="item.id">
        <td @click="!isViewOnly && $emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ item.name }}</td>
        <td>
          <div :style="{backgroundColor: item.color}"><span style="visibility: hidden">{{ item.color }}</span></div>
        </td>
        <td>{{ $filters.dateFormat(item.created_at) }}</td>
        <td>{{ $filters.unixTimeFormat(item.updated_at) }}</td>
        <td v-if="!isViewOnly">
          <a href="#" @click="$emit('moveUp', index)"><span class="material-symbols-outlined">arrow_upward</span></a>
          <a href="#" @click="$emit('moveDown', index)"><span
              class="material-symbols-outlined">arrow_downward</span></a>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {Process} from "@/api";

defineEmits(['openEditModal', 'moveUp', 'moveDown'])

interface AsyncProcessTable {
  list: Process[]
  isViewOnly?: boolean
  isNoHeader?: boolean
  isSimulate?: boolean
}

defineProps<AsyncProcessTable>()

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




