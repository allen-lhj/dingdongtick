<template>
  <div class="settings-panel">
    <n-tabs type="line" animated>
      <n-tab-pane name="general" tab="通用设置">
        <div class="settings-section">
          <h3>外观</h3>
          <n-form-item label="主题模式">
            <n-switch v-model:value="settings.darkMode">
              <template #checked>深色</template>
              <template #unchecked>浅色</template>
            </n-switch>
          </n-form-item>

          <n-form-item label="语言">
            <n-select v-model:value="settings.language" :options="languageOptions" />
          </n-form-item>
        </div>

        <div class="settings-section">
          <h3>通知</h3>
          <n-form-item label="桌面通知">
            <n-switch v-model:value="settings.notifications" />
          </n-form-item>

          <n-form-item label="任务提醒">
            <n-switch v-model:value="settings.taskReminders" />
          </n-form-item>
        </div>
      </n-tab-pane>

      <n-tab-pane name="account" tab="账户设置">
        <div class="settings-section">
          <h3>个人信息</h3>
          <n-form-item label="用户名">
            <n-input v-model:value="settings.username" />
          </n-form-item>

          <n-form-item label="邮箱">
            <n-input v-model:value="settings.email" />
          </n-form-item>
        </div>
      </n-tab-pane>
    </n-tabs>

    <div class="settings-actions">
      <n-button type="primary" @click="handleSave">保存设置</n-button>
      <n-button @click="handleReset">重置</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useTheme } from '@/composables/useTheme'

// 主题管理
const { isDark, setTheme } = useTheme()

const settings = reactive({
  darkMode: isDark.value,
  language: 'zh-CN',
  notifications: true,
  taskReminders: true,
  username: '用户名',
  email: 'user@example.com',
})

// 监听主题变化
watch(
  () => settings.darkMode,
  (newValue) => {
    setTheme(newValue)
  },
)

// 监听全局主题状态变化
watch(isDark, (newValue) => {
  settings.darkMode = newValue
})

const languageOptions = [
  { label: '中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
]

const handleSave = () => {
  // TODO: 保存设置
  console.log('保存设置:', settings)
}

const handleReset = () => {
  // TODO: 重置设置
  console.log('重置设置')
}
</script>

<style scoped>
.settings-panel {
  padding: 20px 0;
}

.settings-section {
  margin-bottom: 24px;
}

.settings-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--n-text-color);
}

.settings-actions {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--n-border-color);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
</style>
