<template>
  <component :is="layoutComponent">
    <slot />
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AuthLayout from '@/layouts/AuthLayout.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const route = useRoute()

// 根据路由元信息选择布局
const layoutComponent = computed(() => {
  const layout = route.meta?.layout as string
  
  switch (layout) {
    case 'auth':
      return AuthLayout
    case 'default':
    default:
      return DefaultLayout
  }
})
</script>
