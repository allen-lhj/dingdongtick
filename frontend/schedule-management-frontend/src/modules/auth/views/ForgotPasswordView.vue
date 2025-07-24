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
          忘记密码
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          输入您的邮箱地址，我们将发送重置密码的链接给您
        </p>
      </div>

      <!-- 忘记密码表单 -->
      <div class="bg-white py-8 px-6 shadow-xl rounded-lg">
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
              type="email"
              placeholder="请输入您的邮箱地址"
              :input-props="{ autocomplete: 'email' }"
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

          <!-- 发送重置链接按钮 -->
          <n-button
            type="primary"
            size="large"
            block
            :loading="isLoading"
            :disabled="!isFormValid"
            attr-type="submit"
            class="mb-4"
          >
            {{ isLoading ? '发送中...' : '发送重置链接' }}
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

        <!-- 成功提示 -->
        <div v-if="emailSent" class="mt-6 p-4 bg-green-50 border border-green-200 rounded-md">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-green-800">
                邮件已发送
              </h3>
              <div class="mt-2 text-sm text-green-700">
                <p>
                  我们已向 <strong>{{ formData.email }}</strong> 发送了密码重置链接。
                  请检查您的邮箱（包括垃圾邮件文件夹）。
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部帮助 -->
      <div class="text-center text-sm text-gray-600">
        <p>
          没有收到邮件？
          <button 
            @click="handleResend"
            :disabled="isLoading || !canResend"
            class="text-indigo-600 hover:text-indigo-500 disabled:text-gray-400 disabled:cursor-not-allowed"
          >
            {{ canResend ? '重新发送' : `${resendCountdown}秒后可重新发送` }}
          </button>
        </p>
        <p class="mt-2">
          需要帮助？
          <a href="#" class="text-indigo-600 hover:text-indigo-500">联系客服</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

// 消息提示
const message = useMessage()

// 表单引用和数据
const formRef = ref<FormInst | null>(null)
const formData = ref({
  email: ''
})

// 状态管理
const isLoading = ref(false)
const emailSent = ref(false)
const canResend = ref(true)
const resendCountdown = ref(0)

// 定时器
let countdownTimer: NodeJS.Timeout | null = null

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
  ]
}

// 计算属性
const isFormValid = computed(() => {
  return formData.value.email && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.value.email)
})

// 处理提交
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    isLoading.value = true
    
    await authService.requestPasswordReset({
      email: formData.value.email
    })
    
    emailSent.value = true
    startResendCountdown()
    
    message.success('密码重置邮件已发送')
    
  } catch (error: any) {
    console.error('Password reset request failed:', error)
    message.error(error.message || '发送失败，请稍后重试')
  } finally {
    isLoading.value = false
  }
}

// 处理重新发送
const handleResend = async () => {
  if (!canResend.value || isLoading.value) return
  
  await handleSubmit()
}

// 开始重新发送倒计时
const startResendCountdown = () => {
  canResend.value = false
  resendCountdown.value = 60
  
  countdownTimer = setInterval(() => {
    resendCountdown.value--
    
    if (resendCountdown.value <= 0) {
      canResend.value = true
      if (countdownTimer) {
        clearInterval(countdownTimer)
        countdownTimer = null
      }
    }
  }, 1000)
}

// 清理定时器
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
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
