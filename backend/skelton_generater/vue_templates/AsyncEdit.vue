<template>
  <div class="container">
@AsyncEditMapping@
    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="post@Upper@(@Lower@, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="post@Upper@ById(@Lower@, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="delete@Upper@ById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {use@Upper@, post@Upper@ById, post@Upper@, delete@Upper@ById} from "@/composable/@Lower@";

interface Async@Upper@Edit {
  id: number | undefined
}

const props = defineProps<Async@Upper@Edit>()
defineEmits(['closeEditModal'])

const {@Lower@} = await use@Upper@(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>

