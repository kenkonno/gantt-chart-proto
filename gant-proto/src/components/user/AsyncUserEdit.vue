<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="user.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Password</label>
      <input class="form-control" type="text" name="password" id="password" v-model="user.password" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Email</label>
      <input class="form-control" type="text" name="email" id="email" v-model="user.email" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">CreatedAt</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="user.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">UpdatedAt</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="user.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postUser(user, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postUserById(user, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteUserById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useUser, postUserById, postUser, deleteUserById} from "@/composable/user";

interface AsyncUserEdit {
  id: number | undefined
}

const props = defineProps<AsyncUserEdit>()
defineEmits(['closeEditModal'])

const {user} = await useUser(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


