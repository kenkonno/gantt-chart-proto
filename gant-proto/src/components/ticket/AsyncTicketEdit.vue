<template>
  <div class="container">
    <div class="d-flex flex-wrap conditions">
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
        <span>{{ getProcessName(ticket.process_id) }}</span>
      </p>
      <p>
        <b class="form-label">部署:</b>
        <span>{{ getDepartmentName(ticket.department_id) }}</span>
      </p>
      <p v-if="false">
        <b class="form-label">期日:</b>
        <span>{{ ticket.limit_date }}</span>
      </p>
      <p>
        <b class="form-label">工数:</b>
        <span>{{ ticket.estimate }}</span>
      </p>
      <p>
        <b class="form-label">日後:</b>
        <span>{{ ticket.days_after }}</span>
      </p>
      <p>
        <b class="form-label">開始日:</b>
        <span>{{ $filters.dateFormatYMD(ticket.start_date) }}</span>
      </p>
      <p>
        <b class="form-label">終了日:</b>
        <span>{{ $filters.dateFormatYMD(ticket.end_date) }}</span>
      </p>
      <p>
        <b class="form-label">進捗:</b>
        <span>{{ ticket.progress_percent ? ticket.progress_percent : 0 }}%</span>
      </p>
      <p>
        <b class="form-label">作成日:</b>
        <span>{{ $filters.dateFormatYMD(ticket.created_at) }}</span>
      </p>
      <p>
        <b class="form-label">更新日:</b>
        <span>{{ $filters.unixTimeFormat(ticket.updated_at) }}</span>
      </p>
    </div>
    <div class="quill-editor">
      <QuillEditor ref="myQuillEditor" theme="snow" v-model="memo"/>
    </div>
    <div class="buttons" v-if="allowed('UPDATE_TICKET')">
      <button type="submit" class="btn btn-primary" @click="updateTicketMemo()">更新</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useTicket, postTicketMemoById} from "@/composable/ticket";
import {inject, onMounted, ref} from "vue";
import {GLOBAL_GETTER_KEY} from "@/composable/globalState";
import {QuillEditor} from "@vueup/vue-quill";
import {Api} from "@/api/axios";
import {allowed} from "@/composable/role";

const {getUnitName, getDepartmentName, getProcessName} = inject(GLOBAL_GETTER_KEY)!

interface AsyncTicketEdit {
  id: number | undefined
  unitId: number
}

const props = defineProps<AsyncTicketEdit>()
const emit = defineEmits(['closeEditModal'])

const {ticket} = await useTicket(props.id)
const memo = ref<string>("")
const {data} = await Api.getTicketMemoId(ticket.value.id!)
memo.value = data.memo
const myQuillEditor = ref(null)

const updateTicketMemo = async () => {
  try {
    const result = await postTicketMemoById(ticket.value.id!, myQuillEditor.value.getHTML(), ticket.value.updated_at)
    emit('closeEditModal', result)
  } catch(e) {
    console.warn(e)
  }
}

onMounted(() => {
  myQuillEditor.value.setHTML(memo.value)
})


</script>

<style scoped lang="scss">
label {
  float: left;
}

.container {
  display: flex;
  flex-direction: column;
  margin: 10px;
  height: 80%;
  overflow: scroll;
}

.conditions {
  border: 1px solid #aaaaaa;
  border-radius: 10px;
  flex: 1;
  min-height: 80px;

  > p {
    margin: 5px;
  }
}

.buttons {
  flex: 1;
  min-height: 50px;
  margin-top: 50px;
}

</style>
<style>
.quill-editor, .ql-container, .ql-editor {
  min-height: 10rem;
}
</style>

