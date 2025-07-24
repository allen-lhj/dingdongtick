<template>
  <div class="max-w-2xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
    <div class="bg-white shadow rounded-lg">
      <!-- 头部 -->
      <div class="px-6 py-4 border-b border-gray-200">
        <h1 class="text-2xl font-bold text-gray-900">修改密码</h1>
        <p class="mt-1 text-sm text-gray-600">为了账户安全，请定期更新您的密码</p>
      </div>

      <!-- 表单内容 -->
      <div class="p-6">
        <n-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          size="large"
          label-placement="top"
          @submit.prevent="handleSubmit"
        >
          <n-form-item path="currentPassword" label="当前密码">
            <n-input
              v-model:value="formData.currentPassword"
              type="password"
              placeholder="请输入当前密码"
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

          <n-form-item path="newPassword" label="新密码">
            <n-input
              v-model:value="formData.newPassword"
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
          <div v-if="formData.newPassword" class="mb-6">
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

          <!-- 密码要求提示 -->
          <div class="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-md">
            <h3 class="text-sm font-medium text-blue-800 mb-2">密码要求：</h3>
            <ul class="text-sm text-blue-700 space-y-1">
              <li class="flex items-center">
                <svg class="h-4 w-4 mr-2" :class="formData.newPassword.length >= 6 ? 'text-green-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                至少6个字符
              </li>
              <li class="flex items-center">
                <svg class="h-4 w-4 mr-2" :class="/[a-zA-Z]/.test(formData.newPassword) && /\d/.test(formData.newPassword) ? 'text-green-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                包含字母和数字
              </li>
              <li class="flex items-center">
                <svg class="h-4 w-4 mr-2" :class="formData.newPassword !== formData.currentPassword && formData.newPassword ? 'text-green-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                与当前密码不同
              </li>
            </ul>
          </div>

          <!-- 操作按钮 -->
          <div class="flex justify-end space-x-3">
            <n-button @click="handleReset">
              重置
            </n-button>
            <n-button
              type="primary"
              :loading="isLoading"
              :disabled="!isFormValid"
              attr-type="submit"
            >
              {{ isLoading ? '修改中...' : '修改密码' }}
            </n-button>
          </div>
        </n-form>
      </div>
    </div>

    <!-- 安全提示 -->
    <div class="mt-6 bg-yellow-50 border border-yellow-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-yellow-800">
            安全提示
          </h3>
          <div class="mt-2 text-sm text-yellow-700">
            <ul class="list-disc list-inside space-y-1">
              <li>修改密码后，您需要在所有设备上重新登录</li>
              <li>请使用强密码，包含大小写字母、数字和特殊字符</li>
              <li>不要与他人分享您的密码</li>
              <li>建议定期更换密码以保障账户安全</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
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
const message = useMessage()

// 表单引用和数据
const formRef = ref<FormInst | null>(null)
const formData = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 状态管理
const isLoading = ref(false)

// 表单验证规则
const rules: FormRules = {
  currentPassword: [
    {
      required: true,
      message: '请输入当前密码',
      trigger: ['input', 'blur']
    }
  ],
  newPassword: [
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
    },
    {
      validator: (rule, value) => {
        if (value && value === formData.value.currentPassword) {
          return new Error('新密码不能与当前密码相同')
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
        if (value !== formData.value.newPassword) {
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
  return formData.value.currentPassword && 
         formData.value.newPassword && 
         formData.value.newPassword.length >= 6 &&
         formData.value.confirmPassword === formData.value.newPassword &&
         formData.value.newPassword !== formData.value.currentPassword
})

// 密码强度计算
const passwordStrength = computed(() => {
  const password = formData.value.newPassword
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
    
    isLoading.value = true
    
    await authService.changePassword({
      currentPassword: formData.value.currentPassword,
      newPassword: formData.value.newPassword
    })
    
    message.success('密码修改成功，请重新登录')
    
    // 延迟跳转到登录页
    setTimeout(() => {
      router.push('/login')
    }, 2000)
    
  } catch (error: any) {
    console.error('Password change failed:', error)
    message.error(error.message || '修改失败，请稍后重试')
  } finally {
    isLoading.value = false
  }
}

// 重置表单
const handleReset = () => {
  formData.value = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
  message.info('表单已重置')
}
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
  .max-w-2xl {
    max-width: 100%;
    margin: 0 1rem;
  }
}
</style>
