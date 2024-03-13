<template>
  <Multiselect
      class="user-multiselect-wrapper"
      v-model="value"
      mode="tags"
      placeholder="担当者を追加"
      :close-on-select="true"
      :search="true"
      :options="options"
      @input="$emit('update:modelValue', $event)"
  >
    <template v-slot:tag="{ option, handleTagRemove, disabled }">
      <div
          :class="{
          'is-disabled': disabled
        }"
          class="multiselect-tag-remove multiselect-tag"
          style="background-color: inherit"
          :style="getStyle(option.value)"
      >
        <div class="icon-wrapper">
          <SingleRune :name="option.label" :id="option.value" class="hover-test"></SingleRune>
          <div class="close" @mousedown.prevent="allowed('UPDATE_TICKET') && handleTagRemove(option, $event)"></div>
        </div>
      </div>
    </template>
  </Multiselect>
</template>

<script setup lang="ts">
import Multiselect from '@vueform/multiselect'
import {TicketUser, User} from "@/api";
import {computed} from "vue";
import SingleRune from "@/components/form/SingleRune.vue";
import {allowed} from "@/composable/role";

interface UserMultiselect {
  ticketUser: TicketUser[]
  userList: User[]
}

const props = defineProps<UserMultiselect>()
defineEmits(['update:modelValue', 'change'])

const options = computed(() => props.userList.map(v => {
  return {value: v.id, label: v.name}
}))
const value = computed(() => props.ticketUser.map(v => v.user_id))

const getStyle = (userId: number) => {
  // indexを取得する
  const index = value.value.findIndex(v => v === userId)
  // データが更新されるタイミングが１ターン後なのでこれを入れている
  if ( index == -1 ) {
    return {display: "none"}
  }
  return {left: index*(50/value.value.length) + "%"}

}
</script>

<style>
.icon-wrapper {
  z-index: 100;
}
.icon-wrapper:hover .close{
  display: block;
}
.close {
  display: none;
  z-index: 101;
  position: absolute;
  width: 16px;
  height: 16px;
  border: 1px solid #333; /* 枠の調整 */
  border-radius: 50%;  /* 丸みの度合い */
  background-color: white;
  right: 0;
  top: 0;
}

.close::before, .close::after { /* 共通設定 */
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  width: 2px; /* 棒の幅（太さ） */
  height: 12px; /* 棒の高さ */
  background: #333;
}

.close::before {
  transform: translate(-50%,-50%) rotate(45deg);
}

.close::after {
  transform: translate(-50%,-50%) rotate(-45deg);
}

.user-multiselect-wrapper .multiselect-clear-icon {
  display: none;
}

.user-multiselect-wrapper .multiselect-dropdown {
  z-index: 200;
  position: absolute;
}

.user-multiselect-wrapper.multiselect,
.user-multiselect-wrapper.multiselect-wrapper {
  height: 26px !important;
  min-height: 26px !important;
  font-size: 10pt;
}

.user-multiselect-wrapper .multiselect-option {
  padding: 2px;
  font-size: 10pt;
}

.user-multiselect-wrapper .multiselect-tags {
  margin: 0;
}
.user-multiselect-wrapper .multiselect-tag {
  position: absolute;
}
</style>