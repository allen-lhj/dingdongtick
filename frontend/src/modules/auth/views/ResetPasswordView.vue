<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- 头部 -->
      <div class="text-center">
        <div class="mx-auto h-12 w-12 bg-indigo-600 rounded-full flex items-center justify-center">
          <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v-2H7v-2H4a1 1 0 01-1-1v-4c0-2.632 2.122-5.367 4.688-5.394A6.977 6.977 0 0112 5a7.021 7.021 0 013 .629" />
          </svg>
        </div>
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
          重置密码
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          请输入您的新密码
        </p>
      </div>

      <!-- 重置密码表单 -->
      <div class="bg-white py-8 px-6 shadow-xl rounded-lg">
        <n-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          size="large"
          @submit.prevent="handleSubmit"
        >
          <n-form-item path="password" label="新密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              placeholder="请输入新密码（至少6位）"
              :input-props="{ autocomplete: 'new-password' }"
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

          <n-form-item path="confirmPassword" label="确认新密码">
            <n-input
              v-model:value="formData.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              :input-props="{ autocomplete: 'new-password' }"
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

          <!-- 密码强度指示器 -->
          <div v-if="formData.password" class="mb-4">
            <div class="text-sm text-gray-600 mb-2">密码强度：</div>
            <div class="flex space-x-1">
              <div 
                v-for="i in 4" 
                :key="i"
                class="h-2 flex-1 rounded-full"
                :class="getPasswordStrengthColor(i)"
              />
            </div>
            <div class="text-xs text-gray-500 mt-1">
              {{ getPasswordStrengthText() }}
            </div>
          </div>

          <!-- 重置密码按钮 -->
          <n-button
            type="primary"
            size="large"
            block
            :loading="isLoading"
            :disabled="!isFormValid"
            attr-type="submit"
            class="mb-4"
          >
            {{ isLoading ? '重置中...' : '重置密码' }}
          </n-button>

          <!-- 返回登录 -->
          <div class="text-center">
            <router-link 
              to="/login" 
              class="text-sm text-indigo-600 hover:text-indigo-500 transition-colors"
            >
              ← 返回登录
            </router-link>
          </div>
        </n-form>
      </div>

      <!-- 底部帮助 -->
      <div class="text-center text-sm text-gray-600">
        <p>
          重置链接已过期？
          <router-link to="/forgot-password" class="text-indigo-600 hover:text-indigo-500">
            重新申请
          </router-link>
        </p>
      </div>
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
  NIcon,
  useMessage,
  type FormInst,
  type FormRules
} from 'naive-ui'
import authService from '../service'

// 路由和消息
const router = useRouter()
const route = useRoute()
const message = useMessage()

// 表单引用和数据
const formRef = ref<FormInst | null>(null)
const formData = ref({
  password: '',
  confirmPassword: ''
})

// 状态管理
const isLoading = ref(false)
const token = ref('')

// 表单验证规则
const rules: FormRules = {
  password: [
    {
      required: true,
      message: '请输入新密码',
      trigger: ['input', 'blur']
    },
    {
      min: 6,
      message: '密码长度不能少于6位',
      trigger: ['input', 'blur']
    },
    {
      validator: (rule, value) => {
        if (!value) return true
        const hasLetter = /[a-zA-Z]/.test(value)
        const hasNumber = /\d/.test(value)
        if (!hasLetter || !hasNumber) {
          return new Error('密码必须包含字母和数字')
        }
        return true
      },
      trigger: ['input', 'blur']
    }
  ],
  confirmPassword: [
    {
      required: true,
      message: '请确认新密码',
      trigger: ['input', 'blur']
    },
    {
      validator: (rule, value) => {
        if (value !== formData.value.password) {
          return new Error('两次输入的密码不一致')
        }
        return true
      },
      trigger: ['input', 'blur']
    }
  ]
}

// 计算属性
const isFormValid = computed(() => {
  return formData.value.password && 
         formData.value.password.length >= 6 &&
         formData.value.confirmPassword === formData.value.password
})

// 密码强度计算
const passwordStrength = computed(() => {
  const password = formData.value.password
  if (!password) return 0
  
  let strength = 0
  
  // 长度检查
  if (password.length >= 6) strength++
  if (password.length >= 8) strength++
  
  // 字符类型检查
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++
  if (/\d/.test(password)) strength++
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) strength++
  
  return Math.min(strength, 4)
})

// 密码强度颜色
const getPasswordStrengthColor = (index: number) => {
  if (index <= passwordStrength.value) {
    if (passwordStrength.value <= 1) return 'bg-red-400'
    if (passwordStrength.value <= 2) return 'bg-yellow-400'
    if (passwordStrength.value <= 3) return 'bg-blue-400'
    return 'bg-green-400'
  }
  return 'bg-gray-200'
}

// 密码强度文本
const getPasswordStrengthText = () => {
  const texts = ['很弱', '弱', '一般', '强', '很强']
  return texts[passwordStrength.value] || '很弱'
}

// 处理提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    if (!token.value) {
      message.error('重置链接无效或已过期')
      return
    }
    
    isLoading.value = true
    
    await authService.resetPassword({
      token: token.value,
      password: formData.value.password
    })
    
    message.success('密码重置成功，请使用新密码登录')
    router.push('/login')
    
  } catch (error: any) {
    console.error('Password reset failed:', error)
    message.error(error.message || '重置失败，请稍后重试')
  } finally {
    isLoading.value = false
  }
}

// 组件挂载时获取token
onMounted(() => {
  token.value = route.query.token as string || ''
  
  if (!token.value) {
    message.error('重置链接无效或已过期')
    router.push('/forgot-password')
  }
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
