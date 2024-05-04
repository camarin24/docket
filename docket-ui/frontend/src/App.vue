<script setup lang="ts">
import type { GlobalTheme, GlobalThemeOverrides } from 'naive-ui';
import { NConfigProvider, NGlobalStyle, darkTheme } from 'naive-ui';
import { ref } from 'vue';
import HomeComponent from './components/HomeComponent.vue';

const theme = ref<GlobalTheme | null>()

const themeOverrides: GlobalThemeOverrides = {
  common: {
    // baseColor: '#fef282',
    // primaryColor: '#fef282'
    fontWeightStrong: '600',
    fontFamily:"'Inter', sans-serif"
  },
}
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <div class="p-4">
      <HomeComponent />
    </div>
    <n-global-style />
  </n-config-provider>
</template>
