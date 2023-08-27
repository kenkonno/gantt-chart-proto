<template>
  <Suspense>
    <TopPage></TopPage>
  </Suspense>
</template>

<style lang="scss" scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

</style>
<script setup lang="ts">
import TopPage from "@/views/TopPage.vue";
</script>