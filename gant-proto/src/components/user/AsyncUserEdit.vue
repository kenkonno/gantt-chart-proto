<template>
  <div class="container">
    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="user.id" :disabled="true">
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="department_id">部署</label>
      <select class="form-control" v-model="user.department_id" name="department_id">
        <option v-for="item in departmentList" :value="item.id" :key="item.id">{{ item.name }}</option>
      </select>
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="id">姓</label>
      <input class="form-control" type="text" name="name" id="name" v-model="user.lastName" :disabled="false" autocomplete="off">
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="id">名</label>
      <input class="form-control" type="text" name="name" id="name" v-model="user.firstName" :disabled="false" autocomplete="off">
    </div>

    <div class="mb-2" v-if="false">
      <label class="form-label" for="limit_of_operation">稼働上限</label>
      <input class="form-control" type="number" name="limit_of_operation" step="0.1" id="limit_of_operation"
             v-model="user.limit_of_operation" :disabled="false" autocomplete="off">
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="role">Role</label>
      <select class="form-control" v-model="user.role" name="role" :disabled="!allowed('CHANGE_ROLE')">
        <option v-for="(name,key) in roleList" :value="key" :key="key">{{ name }}</option>
      </select>
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Password</label>
      <input class="form-control" type="password" name="password" id="password" v-model="user.password"
             :disabled="false" autocomplete="off">
      <small>パスワードは小文字、大文字、数字、特殊文字「 [!@#\$%^&*()] 」を含み8文字以上にしてください。</small>
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="id">Email</label>
      <input class="form-control" type="text" name="email" id="email" v-model="user.email" :disabled="false" autocomplete="off">
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="user.created_at"
             :disabled="true">
    </div>

    <div class="mb-2" v-if="mode != 'reset-password'">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="user.updated_at"
             :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="validate(user, true) && postUser(user, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" :disabled="!allowed('UPDATE_USER')" @click="validate(user) && postUserById(user, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" :disabled="!allowed('UPDATE_USER')" @click="deleteUserById(id, $emit)" v-if="mode != 'reset-password'">削除</button>
      <template v-if="mode == 'profile'">
        <button type="submit" class="btn btn-info float-end" @click="logout()">ログアウト</button>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useUser, postUserById, postUser, deleteUserById, validate} from "@/composable/user";
import {computed, inject} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {RoleTypeMap} from "@/const/common";
import {allowed} from "@/composable/role";
import {Api} from "@/api/axios";
import router from "@/router";

interface AsyncUserEdit {
  id: number | undefined
  mode: string | undefined
}

const {departmentList} = inject(GLOBAL_STATE_KEY)!

const props = defineProps<AsyncUserEdit>()
defineEmits(['closeEditModal'])

const {user} = await useUser(props.id)

const roleList = computed(() => {
  return RoleTypeMap
})

const logout = async () => {
  await Api.postLogout()
  await router.push("/login")
}

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


