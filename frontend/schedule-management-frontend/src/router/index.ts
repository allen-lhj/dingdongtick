import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store'
import { authRoutes } from '@/modules/auth/routes'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        title: '首页',
        requiresAuth: true,
      }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
      meta: {
        title: '关于',
      }
    },
    // 认证相关路由
    ...authRoutes,
    // 404页面
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundView.vue'),
      meta: {
        title: '页面未找到',
      }
    }
  ],
})

// 全局路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - ${import.meta.env.VITE_APP_TITLE || '日程管理系统'}`
  }

  // 如果是首次访问，初始化认证状态
  if (authStore.status === 'guest' && !authStore.isLoading) {
    await authStore.initializeAuth()
  }

  const requiresAuth = to.meta?.requiresAuth ?? false
  const guestOnly = to.meta?.guestOnly ?? false

  if (requiresAuth && !authStore.isAuthenticated) {
    // 需要认证但未登录，重定向到登录页
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else if (guestOnly && authStore.isAuthenticated) {
    // 只允许访客访问但已登录，重定向到首页
    next('/')
  } else {
    // 允许访问
    next()
  }
})

export default router
