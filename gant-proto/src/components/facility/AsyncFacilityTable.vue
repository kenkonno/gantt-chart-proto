<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined, undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
        <th>Id</th>
        <th>名称</th>
        <th>開始日</th>
        <th>終了日</th>
        <th>出荷期日</th>
        <th>ステータス</th>
        <th>受注状況</th>
        <th>作成日</th>
        <th>更新日</th>
        <th>コピー</th>
        <th>並び替え</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(item, index) in list" :key="item.id">
        <td @click="$emit('openEditModal', item.id, undefined)">{{ item.id }}</td>
        <td>{{ item.name }}</td>
        <td>{{ $filters.dateFormatYMD(item.term_from) }}</td>
        <td>{{ $filters.dateFormatYMD(item.term_to) }}</td>
        <td>{{ $filters.dateFormatYMD(item.shipment_due_date) }}</td>
        <td>{{ FacilityStatusMap[item.status]}}</td>
        <td>{{ FacilityTypeMap[item.type]}}</td>
        <td>{{ $filters.dateFormat(item.created_at) }}</td>
        <td>{{ $filters.unixTimeFormat(item.updated_at) }}</td>
        <td>
          <a href="#" @click="$emit('openEditModal',item.id, item.id)"><span class="material-symbols-outlined">note</span></a>
        </td>
        <td>
          <a href="#" @click="$emit('moveUp', index)"><span class="material-symbols-outlined">arrow_upward</span></a>
          <a href="#" @click="$emit('moveDown', index)"><span class="material-symbols-outlined">arrow_downward</span></a>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {Facility} from "@/api";
import {FacilityStatusMap, FacilityTypeMap} from "@/const/common";

defineEmits(['openEditModal', 'moveUp', 'moveDown'])

interface AsyncFacilityTable {
  list: Facility[]
}

defineProps<AsyncFacilityTable>()

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




