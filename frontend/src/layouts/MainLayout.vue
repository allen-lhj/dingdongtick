<template>
  <n-config-provider :theme="theme">
    <n-layout class="main-layout" has-sider>
      <!-- 1. 固定侧边栏 -->
      <n-layout-sider
        bordered
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        :collapsed="sidebarCollapsed"
        show-trigger
        @collapse="sidebarCollapsed = true"
        @expand="sidebarCollapsed = false"
        class="main-sidebar"
      >
        <MainSidebar :collapsed="sidebarCollapsed" />
      </n-layout-sider>

      <!-- 主内容区域 -->
      <n-layout>
        <n-layout has-sider>
          <!-- 2. 项目列表区域 -->
          <n-layout-sider
            bordered
            :width="projectListWidth"
            :collapsed-width="0"
            :collapsed="projectListCollapsed"
            class="project-list-sider"
          >
            <ProjectList :collapsed="projectListCollapsed" @project-select="handleProjectSelect" />
          </n-layout-sider>

          <!-- 右侧内容区域 -->
          <n-layout>
            <n-split direction="horizontal" :default-size="0.6" :min="0.3" :max="0.8">
              <!-- 3. 任务列表区域 -->
              <template #1>
                <div class="task-list-container">
                  <TaskList
                    :project-id="selectedProjectId"
                    :tasks="tasks"
                    :project-list-collapsed="projectListCollapsed"
                    @task-select="handleTaskSelect"
                    @toggle-project-list="projectListCollapsed = !projectListCollapsed"
                  />
                </div>
              </template>

              <!-- 4. 任务详情区域 -->
              <template #2>
                <div class="task-detail-container">
                  <TaskDetail
                    :task="selectedTask"
                    :visible="!!selectedTask"
                    @close="handleTaskDetailClose"
                  />
                </div>
              </template>
            </n-split>
          </n-layout>
        </n-layout>
      </n-layout>
    </n-layout>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useThemeVars } from 'naive-ui'
import { useTheme } from '@/composables/useTheme'
import MainSidebar from '@/components/layout/MainSidebar.vue'
import ProjectList from '@/components/project/ProjectList.vue'
import TaskList from '@/components/task/TaskList.vue'
import TaskDetail from '@/components/task/TaskDetail.vue'

// 主题配置
const { theme, initTheme } = useTheme()
const themeVars = useThemeVars()

// 布局状态
const sidebarCollapsed = ref(false)
const projectListCollapsed = ref(false)
const projectListWidth = ref(280)

// 数据状态
const selectedProjectId = ref<string | null>(null)
const selectedTask = ref<any>(null)
const tasks = ref([])

// 计算属性
const mainContentWidth = computed(() => {
  let width = '100%'
  if (!sidebarCollapsed.value) width = 'calc(100% - 240px)'
  if (!projectListCollapsed.value) width = `calc(${width} - ${projectListWidth.value}px)`
  return width
})

// 事件处理
const handleProjectSelect = (projectId: string) => {
  selectedProjectId.value = projectId
  selectedTask.value = null
  // 加载项目任务
  loadProjectTasks(projectId)
}

const handleTaskSelect = (task: any) => {
  selectedTask.value = task
}

const handleTaskDetailClose = () => {
  selectedTask.value = null
}

const loadProjectTasks = async (projectId: string) => {
  // TODO: 实现任务加载逻辑
  console.log('Loading tasks for project:', projectId)
}

// 响应式处理
const handleResize = () => {
  const width = window.innerWidth
  if (width < 768) {
    sidebarCollapsed.value = true
    projectListCollapsed.value = true
  } else if (width < 1024) {
    projectListCollapsed.value = true
  }
}

// 主题清理函数
let cleanupTheme: (() => void) | null = null

onMounted(() => {
  window.addEventListener('resize', handleResize)
  handleResize()

  // 初始化主题
  cleanupTheme = initTheme()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)

  // 清理主题监听器
  if (cleanupTheme) {
    cleanupTheme()
  }
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  position: fixed;
  top: 0;
  left: 0;
}

.main-sidebar {
  background: var(--n-color);
  border-right: 1px solid var(--n-border-color);
  height: 100vh;
}

.project-list-sider {
  background: var(--n-color);
  border-right: 1px solid var(--n-border-color);
  transition: all 0.3s ease;
  height: 100vh;
}

.task-list-container,
.task-detail-container {
  height: 100vh;
  overflow: hidden;
  background: var(--n-color);
  width: 100%;
}

.task-list-container {
  border-right: 1px solid var(--n-border-color);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .project-list-sider {
    position: absolute;
    z-index: 1000;
    height: 100vh;
  }
}
</style>
