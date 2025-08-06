<template>
  <div class="theme-demo">
    <n-card title="主题切换测试">
      <div class="demo-content">
        <h2>当前主题状态</h2>
        <p>主题模式: {{ isDark ? '深色模式' : '浅色模式' }}</p>
        
        <div class="theme-controls">
          <n-button @click="toggleTheme" type="primary">
            <template #icon>
              <n-icon>
                <component :is="isDark ? SunIcon : MoonIcon" />
              </n-icon>
            </template>
            切换到{{ isDark ? '浅色' : '深色' }}模式
          </n-button>
          
          <n-button @click="setTheme(false)" :type="!isDark ? 'primary' : 'default'">
            设置浅色模式
          </n-button>
          
          <n-button @click="setTheme(true)" :type="isDark ? 'primary' : 'default'">
            设置深色模式
          </n-button>
        </div>
        
        <div class="demo-components">
          <h3>组件展示</h3>
          
          <n-space vertical>
            <n-alert title="信息提示" type="info">
              这是一个信息提示框，用于测试主题效果
            </n-alert>
            
            <n-alert title="成功提示" type="success">
              这是一个成功提示框，用于测试主题效果
            </n-alert>
            
            <n-alert title="警告提示" type="warning">
              这是一个警告提示框，用于测试主题效果
            </n-alert>
            
            <n-alert title="错误提示" type="error">
              这是一个错误提示框，用于测试主题效果
            </n-alert>
          </n-space>
          
          <n-divider />
          
          <n-space>
            <n-button>默认按钮</n-button>
            <n-button type="primary">主要按钮</n-button>
            <n-button type="info">信息按钮</n-button>
            <n-button type="success">成功按钮</n-button>
            <n-button type="warning">警告按钮</n-button>
            <n-button type="error">错误按钮</n-button>
          </n-space>
          
          <n-divider />
          
          <n-input placeholder="输入框测试" />
          
          <n-divider />
          
          <n-select
            :options="[
              { label: '选项1', value: '1' },
              { label: '选项2', value: '2' },
              { label: '选项3', value: '3' }
            ]"
            placeholder="选择框测试"
          />
        </div>
        
        <div class="storage-info">
          <h3>存储信息</h3>
          <p>localStorage 中的主题设置: {{ storedTheme }}</p>
          <p>系统偏好: {{ systemPreference }}</p>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme'
import {
  SunnyOutline as SunIcon,
  MoonOutline as MoonIcon
} from '@vicons/ionicons5'

// 使用主题管理
const { isDark, toggleTheme, setTheme } = useTheme()

// 存储信息
const storedTheme = ref('')
const systemPreference = ref('')

// 获取存储和系统信息
const updateInfo = () => {
  storedTheme.value = localStorage.getItem('theme-preference') || '未设置'
  systemPreference.value = window.matchMedia('(prefers-color-scheme: dark)').matches ? '深色' : '浅色'
}

onMounted(() => {
  updateInfo()
  
  // 监听存储变化
  const interval = setInterval(updateInfo, 1000)
  
  // 清理
  return () => clearInterval(interval)
})
</script>

<style scoped>
.theme-demo {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.demo-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.theme-controls {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.demo-components {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.storage-info {
  background: var(--n-card-color);
  padding: 16px;
  border-radius: 6px;
  border: 1px solid var(--n-border-color);
}

.storage-info h3 {
  margin-top: 0;
  color: var(--n-text-color);
}

.storage-info p {
  margin: 8px 0;
  color: var(--n-text-color-2);
  font-family: monospace;
}
</style>
