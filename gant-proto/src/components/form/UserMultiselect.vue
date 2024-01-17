<template>
  <Multiselect
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
          @mousedown.prevent="handleTagRemove(option, $event)"
          style="background-color: inherit"
          :style="getStyle(option.value)"
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
.multiselect-clear-icon {
  display: none;
}

.multiselect-dropdown {
  z-index: 200;
  position: absolute;
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
</style>