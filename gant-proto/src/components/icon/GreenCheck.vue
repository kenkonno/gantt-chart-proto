<template>
  <span class="check"></span>
</template>

<script setup lang="ts">
</script>
<style scoped lang="scss">
.check {
  content: "";
  display: block;
  position: absolute;
  top: -0.1rem;
  right: -0.5rem;
  width: 10px;
  height: 5px;
  border-left: 2px solid #25AF01;
  border-bottom: 2px solid #25AF01;
  transform: rotate(-45deg);
}
</style>


