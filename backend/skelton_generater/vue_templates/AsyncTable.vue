<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
@AsyncTableHeader@
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
@AsyncTableBody@
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {@Upper@} from "@/api";

defineEmits(['openEditModal'])

interface Async@Upper@Table {
  list: @Upper@[]
}

const props = defineProps<Async@Upper@Table>()

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>


