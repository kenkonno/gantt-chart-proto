<template>
  <Multiselect
      v-model="value"
      mode="tags"
      placeholder="担当者を追加"
      :close-on-select="false"
      :search="true"
      :options="options"
      @input="$emit('update:modelValue', $event)"
  >
    <template v-slot:tag="{ option, handleTagRemove, disabled }">
      <div
          :class="{
          'is-disabled': disabled
        }"
          v-if="!disabled"
          class="multiselect-tag-remove multiselect-tag"
          @mousedown.prevent="handleTagRemove(option, $event)"
          style="background-color: inherit"
      >
        <SingleRune :name="option.label" :id="option.value"></SingleRune>
      </div>
    </template>
  </Multiselect>
</template>

<script setup lang="ts">
import Multiselect from '@vueform/multiselect'
import {TicketUser, User} from "@/api";
import {computed} from "vue";
import SingleRune from "@/components/form/SingleRune.vue";

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

.multiselect-tags {
  margin: 0;
}
.multiselect-tag {
  position: absolute;
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