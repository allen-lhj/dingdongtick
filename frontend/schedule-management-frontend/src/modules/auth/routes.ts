import type { RouteRecordRaw } from 'vue-router'

/**
 * 认证模块路由配置
 */
export const authRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('./views/LoginView.vue'),
    meta: {
      title: '登录',
      guestOnly: true, // 只允许未登录用户访问
      layout: 'auth', // 使用认证布局
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('./views/RegisterView.vue'),
    meta: {
      title: '注册',
      guestOnly: true,
      layout: 'auth',
    }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('./views/ForgotPasswordView.vue'),
    meta: {
      title: '忘记密码',
      guestOnly: true,
      layout: 'auth',
    }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('./views/ResetPasswordView.vue'),
    meta: {
      title: '重置密码',
      guestOnly: true,
      layout: 'auth',
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('./views/ProfileView.vue'),
    meta: {
      title: '个人资料',
      requiresAuth: true, // 需要登录才能访问
    }
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: () => import('./views/ChangePasswordView.vue'),
    meta: {
      title: '修改密码',
      requiresAuth: true,
    }
  }
]

export default authRoutes