<template>
  <div class="container">
    <div class="mb-2">
      <label class="form-label" for="id">Id</label>
      <input class="form-control" type="text" name="id" id="id" v-model="ticket.id" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">GanttGroupId</label>
      <input class="form-control" type="text" name="ganttGroupId" id="ganttGroupId" v-model="ticket.gantt_group_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">ProcessId</label>
      <input class="form-control" type="text" name="processId" id="processId" v-model="ticket.process_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">DepartmentId</label>
      <input class="form-control" type="text" name="departmentId" id="departmentId" v-model="ticket.department_id" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">LimitDate</label>
      <input class="form-control" type="text" name="limitDate" id="limitDate" v-model="ticket.limit_date" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Estimate</label>
      <input class="form-control" type="text" name="estimate" id="estimate" v-model="ticket.estimate" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">DaysAfter</label>
      <input class="form-control" type="text" name="daysAfter" id="daysAfter" v-model="ticket.days_after" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">StartDate</label>
      <input class="form-control" type="text" name="startDate" id="startDate" v-model="ticket.start_date" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">EndDate</label>
      <input class="form-control" type="text" name="endDate" id="endDate" v-model="ticket.end_date" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">ProgressPercent</label>
      <input class="form-control" type="text" name="progressPercent" id="progressPercent" v-model="ticket.progress_percent" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">Order</label>
      <input class="form-control" type="text" name="order" id="order" v-model="ticket.order" :disabled="false">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">作成日</label>
      <input class="form-control" type="text" name="createdAt" id="createdAt" v-model="ticket.created_at" :disabled="true">
    </div>

    <div class="mb-2">
      <label class="form-label" for="id">更新日</label>
      <input class="form-control" type="text" name="updatedAt" id="updatedAt" v-model="ticket.updated_at" :disabled="true">
    </div>

    <template v-if="id == null">
      <button type="submit" class="btn btn-primary" @click="postTicket(ticket, $emit)">更新</button>
    </template>
    <template v-else>
      <button type="submit" class="btn btn-primary" @click="postTicketById(ticket, $emit)">更新</button>
      <button type="submit" class="btn btn-warning" @click="deleteTicketById(id, $emit)">削除</button>
    </template>
  </div>
</template>

<script setup lang="ts">
import {useTicket, postTicketById, postTicket, deleteTicketById} from "@/composable/ticket";

interface AsyncTicketEdit {
  id: number | undefined
}

const props = defineProps<AsyncTicketEdit>()
defineEmits(['closeEditModal'])

const {ticket} = await useTicket(props.id)

</script>

<style scoped lang="scss">
label {
  float: left;
}
</style>


