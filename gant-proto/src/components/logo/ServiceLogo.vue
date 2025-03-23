<template>
  <div class="company-logo">
    <img
        v-if="useImageLogo"
        src="@/assets/TasMaP_logo_color.png"
        :alt="altText"
        class="logo-image"
        :style="{ width: width, height: height }"
    />
    <div v-else class="text-logo">{{ companyName }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue';

const props = defineProps({
  // 画像ロゴを使用するかテキストロゴを使用するか
  useImageLogo: {
    type: Boolean,
    default: true
  },
  // 画像のalt属性
  altText: {
    type: String,
    default: '会社ロゴ'
  },
  // テキストロゴの場合の会社名
  companyName: {
    type: String,
    default: '+MaP'
  },
  // ロゴの幅
  width: {
    type: String,
    default: '300px'
  },
  // ロゴの高さ
  height: {
    type: String,
    default: 'auto'
  }
});
</script>

<style lang="scss" scoped>
.company-logo {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px 0;
}

.logo-image {
  max-width: 100%;
  object-fit: contain;
}

.text-logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #333;
}
</style>