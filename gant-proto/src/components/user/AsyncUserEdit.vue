<template>
    <h1>TEST</h1>
    <!-- component with nested async dependencies -->
    <input type="button" value="新規追加">

    <label>ID<input type="text" :value="user.id"> </label>
    <label>Password<input type="text" :value="user.password"> </label>
    <label>Email<input type="text" :value="user.email"> </label>
    <label>CreatedAt<input type="text" :value="user.created_at"> </label>
    <label>UpdatedAt<input type="text" :value="user.updated_at"> </label>
    <!-- loading state via #fallback slot -->
</template>

<script setup lang="ts">
import {useUser} from "@/composable/user";

interface UserEdit {
  id: string | null
}

defineProps<UserEdit>()

const {user} = await useUser()

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
