<template>
  <div class="main-sidebar-content">
    <!-- 用户头像区域 -->
    <div class="user-section" :class="{ collapsed }">
      <n-avatar
        round
        size="medium"
        src="https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg"
        fallback-src="https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg"
      />
      <div v-if="!collapsed" class="user-info">
        <div class="username">用户名</div>
        <div class="user-status">在线</div>
      </div>
    </div>

    <!-- 主导航菜单 -->
    <n-menu
      :collapsed="collapsed"
      :collapsed-width="64"
      :collapsed-icon-size="20"
      :options="menuOptions"
      :value="activeKey"
      @update:value="handleMenuSelect"
      class="main-menu"
    />

    <!-- 底部功能区 -->
    <div class="sidebar-footer" :class="{ collapsed }">
      <n-tooltip placement="right" :disabled="!collapsed">
        <template #trigger>
          <n-button quaternary circle size="medium" @click="toggleTheme" class="theme-toggle">
            <template #icon>
              <n-icon>
                <component :is="isDark ? SunIcon : MoonIcon" />
              </n-icon>
            </template>
          </n-button>
        </template>
        切换主题
      </n-tooltip>

      <n-tooltip placement="right" :disabled="!collapsed">
        <template #trigger>
          <n-button
            quaternary
            circle
            size="medium"
            @click="showSettings = true"
            class="settings-btn"
          >
            <template #icon>
              <n-icon>
                <SettingsIcon />
              </n-icon>
            </template>
          </n-button>
        </template>
        设置
      </n-tooltip>
    </div>

    <!-- 设置抽屉 -->
    <n-drawer v-model:show="showSettings" :width="400" placement="right">
      <n-drawer-content title="设置">
        <SettingsPanel />
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import { useTheme } from '@/composables/useTheme'
import {
  HomeOutline as HomeIcon,
  CalendarOutline as CalendarIcon,
  FolderOutline as ProjectIcon,
  CheckboxOutline as TaskIcon,
  SettingsOutline as SettingsIcon,
  SunnyOutline as SunIcon,
  MoonOutline as MoonIcon,
  StatsChartOutline as StatsIcon,
  PersonOutline as ProfileIcon,
} from '@vicons/ionicons5'
import SettingsPanel from '@/components/settings/SettingsPanel.vue'

interface Props {
  collapsed: boolean
}

const props = defineProps<Props>()
const router = useRouter()

// 主题管理
const { isDark, toggleTheme } = useTheme()

// 状态管理
const activeKey = ref('home')
const showSettings = ref(false)

// 菜单配置
const menuOptions = computed(() => [
  {
    label: '首页',
    key: 'home',
    icon: renderIcon(HomeIcon),
  },
  {
    label: '今天',
    key: 'today',
    icon: renderIcon(CalendarIcon),
  },
  {
    label: '项目',
    key: 'projects',
    icon: renderIcon(ProjectIcon),
  },
  {
    label: '任务',
    key: 'tasks',
    icon: renderIcon(TaskIcon),
  },
  {
    type: 'divider',
    key: 'divider1',
  },
  {
    label: '统计',
    key: 'stats',
    icon: renderIcon(StatsIcon),
  },
  {
    label: '个人资料',
    key: 'profile',
    icon: renderIcon(ProfileIcon),
  },
])

// 工具函数
function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

// 事件处理
const handleMenuSelect = (key: string) => {
  activeKey.value = key
  router.push(`/${key}`)
}

// toggleTheme 已经从 useTheme 导入，无需重新定义
</script>

<style scoped>
.main-sidebar-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--n-color);
}

.user-section {
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid var(--n-border-color);
  transition: all 0.3s ease;
}

.user-section.collapsed {
  justify-content: center;
  padding: 16px 8px;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.username {
  font-weight: 600;
  font-size: 14px;
  color: var(--n-text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-status {
  font-size: 12px;
  color: var(--n-text-color-2);
  margin-top: 2px;
}

.main-menu {
  flex: 1;
  padding: 8px 0;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--n-border-color);
  display: flex;
  gap: 8px;
  justify-content: center;
}

.sidebar-footer.collapsed {
  flex-direction: column;
  align-items: center;
}

.theme-toggle,
.settings-btn {
  transition: all 0.3s ease;
}

.theme-toggle:hover,
.settings-btn:hover {
  background: var(--n-button-color-hover);
}
</style>
