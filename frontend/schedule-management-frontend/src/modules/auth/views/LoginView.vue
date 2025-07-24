<template>
  <div class="flex items-center justify-center min-h-screen px-4 py-12 bg-gradient-to-br from-blue-50 to-indigo-100 sm:px-6 lg:px-8">
    <div class="w-full max-w-md space-y-8">
      <!-- 头部 -->
      <div class="text-center">

        <h2 class="mt-6 text-2xl font-extrabold text-gray-900 sm:text-3xl">
          登录您的账户
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          还没有账户？
          <router-link to="/register" class="font-medium text-indigo-600 transition-colors hover:text-indigo-500">
            立即注册
          </router-link>
        </p>
      </div>

      <!-- 登录表单 -->
      <div class="px-6 py-8 bg-white rounded-lg shadow-xl">
        <n-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          size="large"
          @submit.prevent="handleSubmit"
        >
          <n-form-item path="email" label="邮箱地址">
            <n-input
              v-model:value="formData.email"
              placeholder="请输入您的邮箱地址"
              :input-props="{ type: 'email', autocomplete: 'email' }"
              clearable
            >
              <template #prefix>
                <n-icon>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                    <polyline points="22,6 12,13 2,6"/>
                  </svg>
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="password" label="密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              placeholder="请输入您的密码"
              :input-props="{ autocomplete: 'current-password' }"
              show-password-on="click"
              clearable
            >
              <template #prefix>
                <n-icon>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <circle cx="12" cy="16" r="1"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                  </svg>
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <!-- 记住我和忘记密码 -->
          <div class="flex items-center justify-between mb-6">
            <n-checkbox v-model:checked="formData.remember">
              记住我
            </n-checkbox>
            <router-link 
              to="/forgot-password" 
              class="text-sm text-indigo-600 transition-colors hover:text-indigo-500"
            >
              忘记密码？
            </router-link>
          </div>

          <!-- 登录按钮 -->
          <n-button
            type="primary"
            size="large"
            block
            :loading="isLoading"
            :disabled="!isFormValid"
            attr-type="submit"
            class="mb-4"
          >
            {{ isLoading ? '登录中...' : '登录' }}
          </n-button>

          <!-- 分割线 -->
          <div class="relative my-6">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-gray-300" />
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 text-gray-500 bg-white">或者</span>
            </div>
          </div>
        </n-form>
      </div>

      <!-- 底部链接 -->
      <!-- <div class="text-sm text-center text-gray-600">
        <p>
          登录即表示您同意我们的
          <a href="#" class="text-indigo-600 hover:text-indigo-500">服务条款</a>
          和
          <a href="#" class="text-indigo-600 hover:text-indigo-500">隐私政策</a>
        </p>
      </div> -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  NForm, 
  NFormItem, 
  NInput, 
  NButton, 
  NCheckbox, 
  NIcon,
  useMessage,
  type FormInst,
  type FormRules
} from 'naive-ui'
import { useAuth } from '@/composables/useAuth'
import type { LoginFormData } from '../types'

// 路由和消息
const router = useRouter()
const route = useRoute()
const message = useMessage()

// 认证相关
const { login, isLoading, error, clearError } = useAuth()

// 表单引用和数据
const formRef = ref<FormInst | null>(null)
const formData = ref<LoginFormData>({
  email: '',
  password: '',
  remember: false
})

// 表单验证规则
const rules: FormRules = {
  email: [
    {
      required: true,
      message: '请输入邮箱地址',
      trigger: ['input', 'blur']
    },
    {
      type: 'email',
      message: '请输入有效的邮箱地址',
      trigger: ['input', 'blur']
    }
  ],
  password: [
    {
      required: true,
      message: '请输入密码',
      trigger: ['input', 'blur']
    },
    {
      min: 6,
      message: '密码长度不能少于6位',
      trigger: ['input', 'blur']
    }
  ]
}

// 计算属性
const isFormValid = computed(() => {
  return formData.value.email && 
         formData.value.password && 
         formData.value.password.length >= 6
})

// 处理登录提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    const redirectTo = route.query.redirect as string
    await login({
      email: formData.value.email,
      password: formData.value.password
    }, redirectTo)
    
  } catch (error: any) {
    console.error('Login validation failed:', error)
  }
}

// 第三方登录处理
const handleGoogleLogin = () => {
  message.info('Google登录功能开发中...')
}

const handleGithubLogin = () => {
  message.info('GitHub登录功能开发中...')
}

// 组件挂载时清除错误
onMounted(() => {
  clearError()
})
</script>

<style scoped>
/* 自定义样式 */
.n-form-item {
  margin-bottom: 1.5rem;
}

.n-input {
  border-radius: 0.5rem;
}

.n-button {
  border-radius: 0.5rem;
  font-weight: 500;
}

/* 响应式调整 */
@media (max-width: 640px) {
  .max-w-md {
    max-width: 100%;
    margin: 0 1rem;
  }
}
</style>
