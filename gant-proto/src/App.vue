<template>
  <Suspense>
    <TopPage></TopPage>
  </Suspense>
</template>

<style lang="scss">
@import '@/assets/common.scss';

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  font-size: 10pt;
  padding: 10px;
}

html, body, #app, #gantt-wrapper {
  height: 100%
}

nav {
  height: 60px;
}

#gantt-wrapper {

  margin-top: -60px;
  padding-top: 60px;
}

</style>
<script setup lang="ts">
import TopPage from "@/views/TopPage.vue";
</script>