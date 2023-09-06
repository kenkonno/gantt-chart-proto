<template>
  <Suspense>
    <TopPage></TopPage>
  </Suspense>
</template>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  font-size: 10pt;
  padding: 10px;
}

</style>
<script setup lang="ts">
import TopPage from "@/views/TopPage.vue";
</script>