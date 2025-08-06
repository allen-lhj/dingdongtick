<template>
  <div class="max-w-4xl px-4 py-8 mx-auto sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow">
      <!-- 头部 -->
      <div class="px-6 py-4 border-b border-gray-200">
        <h1 class="text-2xl font-bold text-gray-900">个人资料</h1>
        <p class="mt-1 text-sm text-gray-600">管理您的账户信息和偏好设置</p>
      </div>

      <!-- 内容区域 -->
      <div class="p-6">
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <!-- 左侧：头像和基本信息 -->
          <div class="lg:col-span-1">
            <div class="text-center">
              <!-- 头像 -->
              <div class="relative inline-block">
                <div class="flex items-center justify-center w-32 h-32 mx-auto bg-indigo-600 rounded-full">
                  <span class="text-4xl font-bold text-white">
                    {{ userInitials }}
                  </span>
                </div>
                <button class="absolute bottom-0 right-0 p-2 bg-white border border-gray-200 rounded-full shadow-lg hover:bg-gray-50">
                  <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
              </div>

              <!-- 用户名 -->
              <h2 class="mt-4 text-xl font-semibold text-gray-900">
                {{ userDisplayName }}
              </h2>
              <p class="text-sm text-gray-600">{{ user?.email }}</p>

              <!-- 账户状态 -->
              <div class="mt-4">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  <svg class="mr-1.5 h-2 w-2 text-green-400" fill="currentColor" viewBox="0 0 8 8">
                    <circle cx="4" cy="4" r="3" />
                  </svg>
                  活跃用户
                </span>
              </div>
            </div>
          </div>

          <!-- 右侧：表单 -->
          <div class="lg:col-span-2">
            <n-form
              ref="formRef"
              :model="formData"
              :rules="rules"
              size="large"
              label-placement="top"
            >
              <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
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
                  placeholder="请输入邮箱地址"
                  :input-props="{ type: 'email' }"
                  disabled
                />
              </n-form-item>

              <n-form-item path="phone" label="手机号码">
                <n-input
                  v-model:value="formData.phone"
                  placeholder="请输入手机号码"
                  clearable
                />
              </n-form-item>

              <n-form-item path="bio" label="个人简介">
                <n-input
                  v-model:value="formData.bio"
                  type="textarea"
                  placeholder="介绍一下自己..."
                  :rows="3"
                  clearable
                />
              </n-form-item>

              <!-- 偏好设置 -->
              <div class="pt-6 mt-6 border-t border-gray-200">
                <h3 class="mb-4 text-lg font-medium text-gray-900">偏好设置</h3>
                
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700">邮件通知</label>
                      <p class="text-sm text-gray-500">接收重要更新和提醒</p>
                    </div>
                    <n-switch v-model:value="formData.preferences.emailNotifications" />
                  </div>

                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700">桌面通知</label>
                      <p class="text-sm text-gray-500">在浏览器中显示通知</p>
                    </div>
                    <n-switch v-model:value="formData.preferences.desktopNotifications" />
                  </div>

                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700">营销邮件</label>
                      <p class="text-sm text-gray-500">接收产品更新和优惠信息</p>
                    </div>
                    <n-switch v-model:value="formData.preferences.marketingEmails" />
                  </div>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="flex justify-end pt-6 mt-6 space-x-3 border-t border-gray-200">
                <n-button @click="handleReset">
                  重置
                </n-button>
                <n-button
                  type="primary"
                  :loading="isLoading"
                  @click="handleSubmit"
                >
                  {{ isLoading ? '保存中...' : '保存更改' }}
                </n-button>
              </div>
            </n-form>
          </div>
        </div>
      </div>
    </div>

    <!-- 危险操作区域 -->
    <div class="mt-8 bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-medium text-red-600">危险操作</h2>
        <p class="mt-1 text-sm text-gray-600">这些操作不可逆，请谨慎操作</p>
      </div>
      
      <div class="p-6">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-sm font-medium text-gray-900">删除账户</h3>
            <p class="text-sm text-gray-500">永久删除您的账户和所有相关数据</p>
          </div>
          <n-button type="error" secondary @click="showDeleteConfirm = true">
            删除账户
          </n-button>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <n-modal v-model:show="showDeleteConfirm">
      <n-card
        style="width: 600px"
        title="确认删除账户"
        :bordered="false"
        size="huge"
        role="dialog"
        aria-modal="true"
      >
        <template #header-extra>
          <n-button quaternary circle @click="showDeleteConfirm = false">
            <template #icon>
              <n-icon><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg></n-icon>
            </template>
          </n-button>
        </template>
        
        <div class="space-y-4">
          <p class="text-gray-700">
            此操作将永久删除您的账户和所有相关数据，包括：
          </p>
          <ul class="space-y-1 text-sm text-gray-600 list-disc list-inside">
            <li>个人资料信息</li>
            <li>所有日程和任务</li>
            <li>项目和协作数据</li>
            <li>设置和偏好</li>
          </ul>
          <p class="font-medium text-red-600">
            此操作不可撤销！
          </p>
        </div>
        
        <template #footer>
          <div class="flex justify-end space-x-3">
            <n-button @click="showDeleteConfirm = false">
              取消
            </n-button>
            <n-button type="error" @click="handleDeleteAccount">
              确认删除
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  NForm,
  NFormItem,
  NInput,
  NButton,
  NSwitch,
  NModal,
  NCard,
  NIcon,
  useMessage,
  type FormInst,
  type FormRules
} from 'naive-ui'
import { useAuth } from '@/composables/useAuth'
import authService from '../service'

// 消息提示
const message = useMessage()

// 认证状态
const { user, userDisplayName, logout } = useAuth()

// 表单引用和数据
const formRef = ref<FormInst | null>(null)
const formData = ref({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  bio: '',
  preferences: {
    emailNotifications: true,
    desktopNotifications: false,
    marketingEmails: false
  }
})

// 状态管理
const isLoading = ref(false)
const showDeleteConfirm = ref(false)

// 用户首字母
const userInitials = computed(() => {
  if (!user.value) return ''
  const firstName = user.value.firstName || ''
  const lastName = user.value.lastName || ''
  return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase()
})

// 表单验证规则
const rules: FormRules = {
  firstName: [
    {
      required: true,
      message: '请输入您的名',
      trigger: ['input', 'blur']
    },
    {
      min: 1,
      max: 50,
      message: '名字长度应在1-50个字符之间',
      trigger: ['input', 'blur']
    }
  ],
  lastName: [
    {
      required: true,
      message: '请输入您的姓',
      trigger: ['input', 'blur']
    },
    {
      min: 1,
      max: 50,
      message: '姓氏长度应在1-50个字符之间',
      trigger: ['input', 'blur']
    }
  ],
  phone: [
    {
      pattern: /^1[3-9]\d{9}$/,
      message: '请输入有效的手机号码',
      trigger: ['input', 'blur']
    }
  ],
  bio: [
    {
      max: 500,
      message: '个人简介不能超过500个字符',
      trigger: ['input', 'blur']
    }
  ]
}

// 初始化表单数据
const initializeFormData = () => {
  if (user.value) {
    formData.value = {
      firstName: user.value.firstName || '',
      lastName: user.value.lastName || '',
      email: user.value.email || '',
      phone: user.value.phone || '',
      bio: user.value.bio || '',
      preferences: {
        emailNotifications: user.value.preferences?.emailNotifications ?? true,
        desktopNotifications: user.value.preferences?.desktopNotifications ?? false,
        marketingEmails: user.value.preferences?.marketingEmails ?? false
      }
    }
  }
}

// 处理提交
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    isLoading.value = true

    await authService.updateProfile({
      firstName: formData.value.firstName,
      lastName: formData.value.lastName,
      phone: formData.value.phone,
      bio: formData.value.bio,
      preferences: formData.value.preferences
    })

    message.success('个人资料更新成功')

  } catch (error: any) {
    console.error('Profile update failed:', error)
    message.error(error.message || '更新失败，请稍后重试')
  } finally {
    isLoading.value = false
  }
}

// 重置表单
const handleReset = () => {
  initializeFormData()
  message.info('已重置为原始数据')
}

// 删除账户
const handleDeleteAccount = async () => {
  try {
    await authService.deleteAccount()
    message.success('账户已删除')
    await logout()
  } catch (error: any) {
    console.error('Account deletion failed:', error)
    message.error(error.message || '删除失败，请稍后重试')
  } finally {
    showDeleteConfirm.value = false
  }
}

// 组件挂载时初始化数据
onMounted(() => {
  initializeFormData()
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
@media (max-width: 768px) {
  .grid-cols-1.lg\\:grid-cols-3 {
    grid-template-columns: 1fr;
  }

  .md\\:grid-cols-2 {
    grid-template-columns: 1fr;
  }
}
</style>
