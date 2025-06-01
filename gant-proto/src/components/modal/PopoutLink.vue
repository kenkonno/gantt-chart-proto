<template>
  <p @click="open" v-if="!disabled" style="border-bottom: 1px solid;">
    <span class="material-symbols-outlined" v-if="icon != null">{{ icon }}</span>
    <span>{{ title }}</span>
  </p>
  <p v-else style="text-decoration: inherit; cursor:inherit">
    <span class="material-symbols-outlined" v-if="icon != null">{{ icon }}</span>
    <span>{{ title }}</span>
  </p>
</template>

<script setup lang="ts">

interface PopoutLink {
  title: string,
  url: string,
  disabled?: boolean
  icon?: string
}

const props = withDefaults(defineProps<PopoutLink>(), {disabled: false})

const open = () => {
  // 方法1: より詳細なウィンドウ設定
  const windowFeatures = [
    'popup=yes',
    'width=1200',
    'height=800',
    'left=100',
    'top=100',
    'scrollbars=yes',
    'resizable=yes',
    'toolbar=no',
    'menubar=no',
    'location=yes',
    'status=yes'
  ].join(',');

  const newWindow = window.open(props.url, '_blank', windowFeatures);

  // 方法2: フォーカスを新しいウィンドウに移す
  if (newWindow) {
    newWindow.focus();
  }
}

</script>

<style scoped>
p {
  margin-right: 5px;
  cursor: pointer;
  display: inline;
  padding: 5px 0;
}

span {
  vertical-align: middle;
}

.material-symbols-outlined {
  font-size: 1.3rem;
}

</style>