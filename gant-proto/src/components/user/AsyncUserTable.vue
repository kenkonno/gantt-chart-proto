<template>
  <div class="container">
    <div v-if="isSimulate" class="mb-2 bg-warning text-center">シミュレーション中のため変更できません。</div>
    <div class="d-flex justify-content-between" v-if="allowed('CSV_UPLOAD')">
      <button type="submit" class="btn btn-primary" @click="$emit('openEditModal',undefined)" v-if="!isViewOnly">
        新規追加
      </button>
      <UploadFile description="CSVアップロード：" @file-upload="postUploadUsersCsvFile" :sample="csvSample"/>
    </div>
    <table class="table">
      <thead v-if="!isNoHeader">
      <tr>
        <th>Id</th>
        <th>部署</th>
        <th>氏名</th>
        <th v-if="false">稼働上限</th>
        <th>Role</th>
        <th v-if="false">Password</th>
        <th>Email</th>
        <th>在籍期間(開始)</th>
        <th>在籍期間(終了)</th>
        <th>作成日</th>
        <th>更新日</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="item in list" :key="item.id">
        <td @click="!isViewOnly && $emit('openEditModal', item.id)">{{ item.id }}</td>
        <td>{{ getDepartmentName(item.department_id) }}</td>
        <td>{{ `${item.lastName} ${item.firstName}` }}</td>
        <td v-if="false">{{ item.limit_of_operation }}</td>
        <td>{{ RoleTypeMap[item.role] }}</td>
        <td v-if="false">{{ item.password }}</td>
        <td class="text-break">{{ item.email }}</td>
        <td>{{ $filters.dateFormatYMD(item.employment_start_date) }}</td>
        <td>{{ $filters.dateFormatYMD(item.employment_end_date) }}</td>
        <td>{{ $filters.dateFormat(item.created_at) }}</td>
        <td>{{ $filters.unixTimeFormat(item.updated_at) }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import {User} from "@/api";
import {inject} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {RoleTypeMap} from "../../const/common";
import UploadFile from "@/components/form/UploadFile.vue";
import {postUploadUsersCsvFile} from "@/composable/user";
import {allowed} from "@/composable/role";

defineEmits(['openEditModal'])
const {departmentList} = inject(GLOBAL_STATE_KEY)!
const getDepartmentName = (id: number) => {
  return departmentList.find(v => v.id === id)?.name
}


interface AsyncUserTable {
  list: User[]
  isViewOnly?: boolean
  isNoHeader?: boolean
  isSimulate?: boolean
}

defineProps<AsyncUserTable>()

const csvSample = `部署ID,性,名,Role,Email,Password,在籍期間(開始),在籍期間(終了)
1,タスマップ,太郎,マネージャー,abcd@tasmap.com,password,2025-06-14,
`.split("\n").map(v => v.split(","))

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
tr > td:nth-child(1) {
  text-decoration: underline;
  cursor: pointer;
}
</style>




