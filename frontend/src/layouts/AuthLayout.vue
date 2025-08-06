<template>
  <div class="flex flex-col justify-center min-h-screen py-12 bg-gray-50 sm:px-6 lg:px-8">
    <!-- 主要内容区域 -->
    <div class="flex items-center justify-center flex-1">
      <div class="w-full max-w-md space-y-8">
        <slot />
      </div>
    </div>


    <!-- 全局错误提示 -->
    <div 
      v-if="globalError" 
      class="fixed z-50 w-full max-w-sm top-4 right-4"
    >
      <div class="p-4 border border-red-200 rounded-md bg-red-50">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">
              操作失败
            </h3>
            <div class="mt-2 text-sm text-red-700">
              {{ globalError }}
            </div>
          </div>
          <div class="pl-3 ml-auto">
            <div class="-mx-1.5 -my-1.5">
              <button
                @click="clearGlobalError"
                class="inline-flex bg-red-50 rounded-md p-1.5 text-red-500 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-red-50 focus:ring-red-600"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载遮罩 -->
    <div 
      v-if="isGlobalLoading" 
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    >
      <div class="flex items-center p-6 space-x-3 bg-white rounded-lg">
        <div class="w-6 h-6 border-b-2 border-indigo-600 rounded-full animate-spin"></div>
        <span class="font-medium text-gray-900">处理中...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

// 路由信息
const route = useRoute()

// 认证状态
const { error, isLoading, clearError } = useAuth()

// 全局错误和加载状态
const globalError = computed(() => error.value)
const isGlobalLoading = computed(() => isLoading.value)

// 清除全局错误
const clearGlobalError = () => {
  clearError()
}
</script>

<style scoped>
/* 自定义样式 */
.router-link-active {
  color: #4f46e5;
}

/* 动画效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式调整 */
@media (max-width: 640px) {
  .max-w-md {
    max-width: 100%;
    margin: 0 1rem;
  }
}

/* 选择框样式 */
select {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 0.5rem center;
  background-repeat: no-repeat;
  background-size: 1.5em 1.5em;
  padding-right: 2.5rem;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f5f9;
}

::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
