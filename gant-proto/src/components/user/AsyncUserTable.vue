<template>
  <div class="container">
    <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)">新規追加</button>
    <table class="table">
      <thead>
      <tr>
        <th>id</th>
        <th>password</th>
        <th>email</th>
        <th>created_at</th>
        <th>updated_at</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in users" :key="item.id">
        <td @click="$emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ item.password }}</td>
        <td>{{ item.email }}</td>
        <td>{{ item.created_at }}</td>
        <td>{{ item.updated_at }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {useUserTable} from "@/composable/user";

const users = await useUserTable()
defineEmits(['openEditModal'])

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
h3 {
  margin: 40px 0 0;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
