<template>
  <div class="flex items-center justify-center min-h-screen px-4 py-12 bg-gradient-to-br from-green-50 to-emerald-100 sm:px-6 lg:px-8">
    <div class="w-full max-w-md space-y-8">
      <!-- 头部 -->
      <div class="text-center">
        <h2 class="mt-6 text-2xl font-extrabold text-gray-900 sm:text-3xl">
          创建新账户
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          已有账户？
          <router-link to="/login" class="font-medium transition-colors text-emerald-600 hover:text-emerald-500">
            立即登录
          </router-link>
        </p>
      </div>

      <!-- 注册表单 -->
      <div class="px-6 py-8 bg-white rounded-lg shadow-xl">
        <n-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          size="large"
          @submit.prevent="handleSubmit"
        >
          <!-- 姓名输入 -->
          <div class="grid grid-cols-2 gap-4">
            <n-form-item path="firstName" label="名">
              <n-input
                v-model:value="formData.firstName"
                placeholder="请输入您的名"
                clearable
              />
            </n-form-item>

            <n-form-item path="lastName" label="姓">
              <n-input
                v-model:value="formData.lastName"
                placeholder="请输入您的姓"
                clearable
              />
            </n-form-item>
          </div>

          <n-form-item path="email" label="邮箱地址">
            <n-input
              v-model:value="formData.email"
              placeholder="请输入您的邮箱地址"
              clearable
            />
          </n-form-item>

          <n-form-item path="password" label="密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              show-password-on="click"
              clearable
            />
          </n-form-item>

          <n-form-item path="confirmPassword" label="确认密码">
            <n-input
              v-model:value="formData.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              show-password-on="click"
              clearable
            />
          </n-form-item>


          <!-- 注册按钮 -->
          <n-button
            type="primary"
            size="large"
            block
            :loading="isLoading"
            :disabled="!isFormValid"
            attr-type="submit"
            class="mb-4"
            color="#10b981"
          >
            {{ isLoading ? '注册中...' : '创建账户' }}
          </n-button>
        </n-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { 
  NForm, 
  NFormItem, 
  NInput, 
  NButton, 
  NCheckbox,
  useMessage,
  type FormInst,
  type FormRules
} from 'naive-ui'
import { useAuth } from '@/composables/useAuth'
import type { RegisterFormData } from '../types'

const message = useMessage()
const { register, isLoading, clearError } = useAuth()

const formRef = ref<FormInst | null>(null)
const formData = ref<RegisterFormData>({
  email: '',
  password: '',
  confirmPassword: '',
  firstName: '',
  lastName: '',
  agreeToTerms: false
})

const rules: FormRules = {
  firstName: [
    { required: true, message: '请输入您的名', trigger: ['input', 'blur'] }
  ],
  lastName: [
    { required: true, message: '请输入您的姓', trigger: ['input', 'blur'] }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: ['input', 'blur'] },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['input', 'blur'] }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: ['input', 'blur'] },
    { min: 6, message: '密码长度不能少于6位', trigger: ['input', 'blur'] }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: ['input', 'blur'] },
    {
      validator: (_rule, value) => {
        if (value !== formData.value.password) {
          return new Error('两次输入的密码不一致')
        }
        return true
      },
      trigger: ['input', 'blur']
    }
  ],
  agreeToTerms: [
    {
      validator: (_rule, value) => {
        if (!value) {
          return new Error('请同意服务条款和隐私政策')
        }
        return true
      },
      trigger: ['change']
    }
  ]
}

const isFormValid = computed(() => {
  return formData.value.firstName && 
         formData.value.lastName &&
         formData.value.email && 
         formData.value.password && 
         formData.value.password.length >= 6 &&
         formData.value.confirmPassword === formData.value.password &&
         formData.value.agreeToTerms
})

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    await register({
      email: formData.value.email,
      password: formData.value.password,
      firstName: formData.value.firstName,
      lastName: formData.value.lastName
    })
    
    message.success('注册成功！')
    
  } catch (error: any) {
    console.error('Register failed:', error)
    message.error(error.message || '注册失败，请稍后重试')
  }
}

onMounted(() => {
  clearError()
})
</script>

<style scoped>
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

@media (max-width: 640px) {
  .max-w-md {
    max-width: 100%;
    margin: 0 1rem;
  }
  
  .grid-cols-2 {
    grid-template-columns: 1fr;
  }
}
</style>
