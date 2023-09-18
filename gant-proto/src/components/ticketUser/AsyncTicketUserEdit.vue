<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="ticketUser.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">TicketId</label>
      <input class="form-control" type="text" name="ticketId" id="ticketId" v-model="ticketUser.ticket_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">UserId</label>
      <input class="form-control" type="text" name="userId" id="userId" v-model="ticketUser.user_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="ticketUser.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="ticketUser.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postTicketUser(ticketUser, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postTicketUserById(ticketUser, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteTicketUserById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useTicketUser, postTicketUserById, postTicketUser, deleteTicketUserById} from "@/composable/ticketUser";

interface AsyncTicketUserEdit {
  id: number | undefined
}

const props = defineProps<AsyncTicketUserEdit>()
defineEmits(['closeEditModal'])

const {ticketUser} = await useTicketUser(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


