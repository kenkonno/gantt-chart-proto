<template>
  <div
      class="single-rune"
      :style="getBgColor(id)"
  >
    {{ name }}
  </div>
</template>

<script setup lang="ts">

interface SingleRuneProps {
  name: string
  id: number
}

defineProps<SingleRuneProps>()
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

<style scoped>
.single-rune {
  width: 24px;
  height: 24px;
  line-height: 24px;
  border-radius: 50% !important;
  color: #fff;
  text-align: center;
  z-index: 10;
  border: solid 1px white;
  font-size: 0.7rem;
}
.single-rune:hover {
  opacity: 0.8;
}
</style>