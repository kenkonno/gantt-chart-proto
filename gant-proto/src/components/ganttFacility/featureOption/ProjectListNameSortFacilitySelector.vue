<template>
  <select style="display: inline" :value="modelValue"
          @input="$emit('update:modelValue', Number($event.target.value))">
    <option v-for="item in sortedFacilityList" :key="item.id" :value="item.id">{{ item.name }}
      <template v-if="item.type === FacilityType.Ordered">✅</template>
    </option>
  </select>
</template>

<script setup lang="ts">
import {computed} from "vue";
import {FacilityType} from "@/const/common";

interface Props {
  modelValue: number;
  facilityList: any[];
}

const props = defineProps<Props>();
const emit = defineEmits(['update:modelValue']);

const sortedFacilityList = computed(() => {
  return [...props.facilityList].sort((a, b) => {
    if (a.name < b.name) {
      return -1;
    } else {
      return 1
    }
  })
})

</script>

<style scoped>
select {
  margin: 0 5px;
}
</style>