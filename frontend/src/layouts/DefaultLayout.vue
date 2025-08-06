<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <nav class="bg-white border-b border-gray-200 shadow-sm">
      <div class="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <!-- 左侧 Logo 和导航 -->
          <div class="flex items-center">
            <div class="flex items-center flex-shrink-0">
              <div class="flex items-center justify-center w-8 h-8 bg-indigo-600 rounded-lg">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3a4 4 0 118 0v4m-4 8a2 2 0 100-4 2 2 0 000 4zm0 0v3m-4-3h8m-4-3V8a1 1 0 011-1h6a1 1 0 011 1v8a1 1 0 01-1 1H9a1 1 0 01-1-1V8z" />
                </svg>
              </div>
              <span class="ml-2 text-xl font-bold text-gray-900">TickTick</span>
            </div>
            
            <!-- 主导航 -->
            <div class="hidden md:ml-6 md:flex md:space-x-8">
              <router-link 
                to="/" 
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 transition-colors border-b-2 border-transparent hover:border-gray-300 hover:text-gray-700"
                active-class="text-gray-900 border-indigo-500"
              >
                首页
              </router-link>
              <router-link 
                to="/calendar" 
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 transition-colors border-b-2 border-transparent hover:border-gray-300 hover:text-gray-700"
                active-class="text-gray-900 border-indigo-500"
              >
                日历
              </router-link>
              <router-link 
                to="/projects" 
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 transition-colors border-b-2 border-transparent hover:border-gray-300 hover:text-gray-700"
                active-class="text-gray-900 border-indigo-500"
              >
                项目
              </router-link>
            </div>
          </div>

          <!-- 右侧用户菜单 -->
          <div class="flex items-center">
            <div v-if="isAuthenticated" class="flex items-center space-x-4">
              <!-- 通知按钮 -->
              <button class="p-1 text-gray-400 rounded-full hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM10.5 3.75a6 6 0 0 1 6 6v2.25l2.25 2.25v2.25H2.25V14.25L4.5 12V9.75a6 6 0 0 1 6-6z" />
                </svg>
              </button>

              <!-- 用户头像和菜单 -->
              <div class="relative">
                <button 
                  @click="showUserMenu = !showUserMenu"
                  class="flex items-center space-x-2 text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  <div class="flex items-center justify-center w-8 h-8 bg-indigo-600 rounded-full">
                    <span class="font-medium text-white">
                      {{ userInitials }}
                    </span>
                  </div>
                  <span class="hidden font-medium text-gray-700 md:block">
                    {{ userDisplayName }}
                  </span>
                  <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </button>

                <!-- 用户菜单下拉 -->
                <div 
                  v-if="showUserMenu"
                  @click.away="showUserMenu = false"
                  class="absolute right-0 z-50 w-48 mt-2 origin-top-right bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
                >
                  <div class="py-1">
                    <router-link 
                      to="/profile" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      @click="showUserMenu = false"
                    >
                      个人资料
                    </router-link>
                    <router-link 
                      to="/settings" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      @click="showUserMenu = false"
                    >
                      设置
                    </router-link>
                    <div class="border-t border-gray-100"></div>
                    <button 
                      @click="handleLogout"
                      class="block w-full px-4 py-2 text-sm text-left text-gray-700 hover:bg-gray-100"
                    >
                      退出登录
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- 未登录状态 -->
            <div v-else class="flex items-center space-x-4">
              <router-link 
                to="/login" 
                class="font-medium text-gray-500 hover:text-gray-700"
              >
                登录
              </router-link>
              <router-link 
                to="/register" 
                class="px-4 py-2 font-medium text-white transition-colors bg-indigo-600 rounded-md hover:bg-indigo-700"
              >
                注册
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要内容区域 -->
    <main class="flex-1">
      <slot />
    </main>

    <!-- 底部 -->
    <footer class="mt-auto bg-white border-t border-gray-200">
      <div class="px-4 py-6 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-500">
            © 2024 TickTick. 保留所有权利。
          </div>
          <div class="flex space-x-6">
            <a href="#" class="text-sm text-gray-500 hover:text-gray-900">关于</a>
            <a href="#" class="text-sm text-gray-500 hover:text-gray-900">帮助</a>
            <a href="#" class="text-sm text-gray-500 hover:text-gray-900">隐私</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuth } from '@/composables/useAuth'

// 认证状态
const { user, isAuthenticated, userDisplayName, logout } = useAuth()

// 用户菜单状态
const showUserMenu = ref(false)

// 用户首字母
const userInitials = computed(() => {
  if (!user.value) return ''
  const firstName = user.value.firstName || ''
  const lastName = user.value.lastName || ''
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase()
})

// 处理登出
const handleLogout = async () => {
  showUserMenu.value = false
  await logout()
}
</script>

<style scoped>
/* 点击外部关闭菜单的指令 */
.click-away {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 40;
}
</style>
