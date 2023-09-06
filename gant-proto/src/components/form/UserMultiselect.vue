<template>
  <Multiselect
      v-model="value"
      mode="tags"
      placeholder="担当者を追加"
      :close-on-select="false"
      :search="true"
      :options="options"
      @input="$emit('update', $event)"
  >
    <template v-slot:tag="{ option, handleTagRemove, disabled }">
      <div
          :class="{
          'is-disabled': disabled
        }"
          v-if="!disabled"
          class="multiselect-tag-remove multiselect-tag"
          @mousedown.prevent="handleTagRemove(option, $event)"
          :style="getBgColor(option.value)"
      >
        {{ option.label.substring(0, 1) }}
      </div>
    </template>
  </Multiselect>
</template>

<script setup lang="ts">
import Multiselect from '@vueform/multiselect'
import {TicketUser, User} from "@/api";

interface UserMultiselect {
  ticketUser: TicketUser[]
  userList: User[]
}

const props = defineProps<UserMultiselect>()
defineEmits(['update:modelValue', 'change'])

const options = props.userList.map(v => {
  return {value: v.id, label: v.name}
})
const value = props.ticketUser.map(v => v.user_id)
const update = (value: any) => {
  console.log(value)
}

const getBgColor = (value: number) => {
  const master = [
    'background-color: #D97E4C;',
    'background-color: #D18583;',
    'background-color: #C7C9D1;',
    'background-color: #817D5C;',
    'background-color: #8B42D0;',
    'background-color: #557470;',
    'background-color: #55567E;',
    'background-color: #DB7C49;',
    'background-color: #E58376;',
    'background-color: #7C8E58;',
    'background-color: #7B93B3;',
  ]

  return master[value % 11]
}


</script>

<style>
.multiselect-clear-icon {
  display: none;
}

.multiselect, .multiselect-wrapper {
  height: 26px !important;
  min-height: 26px !important;
  font-size: 10pt;
}
.multiselect-option {
  padding: 2px;
  font-size: 10pt;
}
.multiselect-tags{
  margin: 0;
}
.multiselect-tag {
  width: 24px;
  height: 24px;
  line-height: 24px;
  border-radius: 50%;
  color: #fff;
  text-align: center;
  position: absolute;
  z-index: 10;
  border: solid 1px white;
}

.multiselect-tag:hover {
  opacity: 0.8;
}


.multiselect-tag:nth-child(1) {
  left: 0%
}

.multiselect-tag:nth-child(2) {
  left: 15%
}

.multiselect-tag:nth-child(3) {
  left: 30%
}

.multiselect-tag:nth-child(4) {
  left: 45%
}

.multiselect-tag:nth-child(5) {
  left: 60%
}

.multiselect-tag:nth-child(6) {
  left: 75%
}

.multiselect-tag:nth-child(7) {
  left: 90%
}

.multiselect-tag:nth-child(8) {
  left: 105%
}

.multiselect-tag:nth-child(9) {
  left: 120%
}
</style>