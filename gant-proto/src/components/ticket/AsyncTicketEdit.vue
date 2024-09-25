<template>
  <div class="async-ticket-edit-container">
    <div class="d-flex flex-wrap conditions">
      <div class="d-flex">
        <p>
          <b class="form-label">Id:</b>
          <span>{{ ticket.id }}</span>
        </p>
        <p>
          <b class="form-label">ユニット:</b>
          <span>{{ getUnitName(props.unitId) }}</span>
        </p>
        <p>
          <b class="form-label">工程:</b>
          <select v-model="ticket.process_id"
                  :disabled="!allowed('UPDATE_TICKET')">
            <option v-for="item in processList" :key="item.id" :value="item.id">{{ item.name }}</option>
          </select>
        </p>
        <p>
          <b class="form-label">部署:</b>
          <select v-model="ticket.department_id"
                  :disabled="!allowed('UPDATE_TICKET')">
            <option v-for="item in cDepartmentList" :key="item.id" :value="item.id">{{ item.name }}</option>
          </select>
        </p>
        <p v-if="false">
          <b class="form-label">期日:</b>
          <span>{{ ticket.limit_date }}</span>
        </p>
        <p class="d-flex text-nowrap" style="min-width: 12rem;">
          <b class="form-label">担当者:</b>
          <UserMultiselect :userList="getUserListByDepartmentId(ticket.department_id)"
                           :ticketUser="ticketUsers"
                           :disabled="!allowed('UPDATE_TICKET')"
                           @update:modelValue="updateTicketUsers($event)"></UserMultiselect>
        </p>
        <p>
          <b class="form-label">人数:</b>
          <FormNumber class="small-numeric"
                      v-model="ticket.number_of_worker"
                      :disabled="ticketUsers?.length > 0 || !allowed('UPDATE_TICKET')"/>
        </p>
      </div>
      <div class="d-flex">
        <p>
          <b class="form-label">工数:</b>
          <FormNumber class="small-numeric"
                      v-model="ticket.estimate"
                      :disabled="!allowed('UPDATE_TICKET')"/>

        </p>
        <p>
          <b class="form-label">日後:</b>
          <FormNumber class="small-numeric"
                      v-model="ticket.days_after"
                      :disabled="!allowed('UPDATE_TICKET')"/>
        </p>
        <p>
          <b class="form-label">開始日:</b>
          <input type="date"
                 v-model="ticket.start_date"
                 :disabled="!allowed('UPDATE_TICKET')"/>
        </p>
        <p>
          <b class="form-label">終了日:</b>
          <input type="date"
                 v-model="ticket.end_date"
                 :disabled="!allowed('UPDATE_TICKET')"/>
        </p>
        <p>
          <b class="form-label">進捗:</b>
          <FormNumber class="middle-numeric"
                      v-model="ticket.progress_percent"
                      :disabled="!allowed('UPDATE_PROGRESS')"
                      :min=0 />
        </p>
      </div>
      <div class="d-flex">
        <p>
          <b class="form-label">作成日:</b>
          <span>{{ $filters.dateFormatYMD(ticket.created_at) }}</span>
        </p>
        <p>
          <b class="form-label">更新日:</b>
          <span>{{ $filters.unixTimeFormat(ticket.updated_at) }}</span>
        </p>
      </div>
    </div>
    <div class="editor mt-2">
      <tiptap-editor v-model="ticket.memo"/>
    </div>
    <div class="buttons mt-2" v-if="allowed('UPDATE_TICKET')">
      <button type="submit" class="btn btn-primary" @click="updateTicketMemo()">更新</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useTicket} from "@/composable/ticket";
import {computed, inject, ref} from "vue";
import {GLOBAL_GETTER_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {allowed} from "@/composable/role";
import TiptapEditor from "@/components/tiptap/TiptapEditor.vue";
import {Api} from "@/api/axios";
import {Department, TicketUser} from "@/api";
import UserMultiselect from "@/components/form/UserMultiselect.vue";
import FormNumber from "@/components/form/FormNumber.vue";

const {getUnitName, getDepartmentName, getProcessName} = inject(GLOBAL_GETTER_KEY)!
const {processList, departmentList, facilityTypes, userList} = inject(GLOBAL_STATE_KEY)!
const cDepartmentList = computed(() => {
  const result: Department[] = [{created_at: "", id: undefined, name: "", order: 0, updated_at: 0}]
  result.push(...departmentList)
  return result
})

interface AsyncTicketEdit {
  id: number | undefined
  unitId: number
}

const props = defineProps<AsyncTicketEdit>()
const emit = defineEmits(['closeEditModal'])

// チケットは memo を通常取り扱わないので２個APIを呼び出している。
const {ticket} = await useTicket(props.id)
const {data} = await Api.getTicketMemoId(ticket.value.id!)
ticket.value.memo = data.memo

// TODO: この日付関連はかなり良くないのでだが、負債ということでどこかで治す。
ticket.value.start_date = ticket.value.start_date?.substring(0, 10)
ticket.value.end_date = ticket.value.end_date?.substring(0, 10)
ticket.value.limit_date = ticket.value.limit_date?.substring(0, 10)

const updateTicketMemo = async () => {
  emit('closeEditModal', ticket.value, ticketUsers.value.map( v => v.user_id))
}

// TODO: 重複コード
const getUserListByDepartmentId = (departmentId?: number) => {
  if (departmentId == null) {
    return userList
  }
  return userList.filter(v => v.department_id === departmentId)
}

const ticketId = ticket.value.id ? ticket.value.id : 0
const {data: TicketUserData} = await Api.getTicketUsers([ticketId])
const ticketUsers = ref(TicketUserData.list)


const updateTicketUsers = (value: number[]) => {
  ticketUsers.value.length = 0
  const users = value.map((v,i) => {
    return {
      id: undefined,
      ticket_id: ticket.value.id,
      user_id: v,
      order: i,
      created_at: "",
      updated_at: 0,
    } as TicketUser
  })
  ticketUsers.value.push(...users)
}

</script>

<style scoped lang="scss">
label {
  float: left;
}

.async-ticket-edit-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  margin: 10px;
  overflow: scroll;
  width: 100%;

  &::-webkit-scrollbar {
    display: none;
  }

  -ms-overflow-style: none;
  scrollbar-width: none;

  .conditions {
    flex-basis: auto;
    border: 1px solid #aaaaaa;
    border-radius: 10px;
    min-height: 80px;

    p {
      margin: 5px;
      display: flex;
      text-wrap: nowrap;
    }
  }

  .editor {
    flex-grow: 1;
    flex-basis: 0;
    overflow: scroll;
  }
}

.buttons {
  flex-basis: 50px;
  min-height: 50px;
}

input, select {
  height: 24px;
}

</style>
