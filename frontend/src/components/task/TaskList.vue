<template>
  <div class="task-list-container">
    <!-- 头部操作栏 -->
    <div class="task-header">
      <div class="header-left">
        <!-- 项目列表展开/收起按钮 -->
        <n-tooltip>
          <template #trigger>
            <n-button quaternary circle size="small" @click="$emit('toggle-project-list')">
              <template #icon>
                <n-icon>
                  <MenuIcon />
                </n-icon>
              </template>
            </n-button>
          </template>
          {{ projectListCollapsed ? '展开项目列表' : '收起项目列表' }}
        </n-tooltip>

        <div class="header-title">
          <n-icon size="18">
            <CheckboxIcon />
          </n-icon>
          <span>任务列表</span>
          <n-tag v-if="tasks.length > 0" size="small" type="info">
            {{ completedCount }}/{{ tasks.length }}
          </n-tag>
        </div>
      </div>

      <div class="header-actions">
        <n-tooltip>
          <template #trigger>
            <n-button quaternary circle size="small" @click="showCreateTask = true">
              <template #icon>
                <n-icon>
                  <AddIcon />
                </n-icon>
              </template>
            </n-button>
          </template>
          新建任务
        </n-tooltip>

        <n-dropdown :options="viewOptions" @select="handleViewChange" trigger="click">
          <n-button quaternary circle size="small">
            <template #icon>
              <n-icon>
                <FilterIcon />
              </n-icon>
            </template>
          </n-button>
        </n-dropdown>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filter-section">
      <n-input v-model:value="searchQuery" placeholder="搜索任务..." clearable size="small">
        <template #prefix>
          <n-icon>
            <SearchIcon />
          </n-icon>
        </template>
      </n-input>

      <div class="filter-tabs">
        <n-button-group size="small">
          <n-button
            :type="currentFilter === 'all' ? 'primary' : 'default'"
            @click="currentFilter = 'all'"
          >
            全部
          </n-button>
          <n-button
            :type="currentFilter === 'pending' ? 'primary' : 'default'"
            @click="currentFilter = 'pending'"
          >
            待完成
          </n-button>
          <n-button
            :type="currentFilter === 'completed' ? 'primary' : 'default'"
            @click="currentFilter = 'completed'"
          >
            已完成
          </n-button>
        </n-button-group>
      </div>
    </div>

    <!-- 任务列表内容 -->
    <div class="task-list-content">
      <n-scrollbar>
        <div v-if="filteredTasks.length === 0" class="empty-state">
          <n-empty description="暂无任务">
            <template #icon>
              <n-icon size="48">
                <CheckboxIcon />
              </n-icon>
            </template>
            <template #extra>
              <n-button @click="showCreateTask = true" type="primary"> 创建第一个任务 </n-button>
            </template>
          </n-empty>
        </div>

        <div v-else class="task-items">
          <!-- 按日期分组显示 -->
          <div v-for="group in groupedTasks" :key="group.date" class="task-group">
            <div class="group-header">
              <span class="group-date">{{ group.dateLabel }}</span>
              <span class="group-count">{{ group.tasks.length }} 个任务</span>
            </div>

            <TaskItem
              v-for="task in group.tasks"
              :key="task.id"
              :task="task"
              :active="selectedTaskId === task.id"
              @select="handleTaskSelect"
              @toggle-complete="handleToggleComplete"
              @update="handleTaskUpdate"
            />
          </div>
        </div>
      </n-scrollbar>
    </div>

    <!-- 创建任务对话框 -->
    <CreateTaskModal
      v-model:show="showCreateTask"
      :project-id="projectId"
      @created="handleTaskCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  CheckboxOutline as CheckboxIcon,
  AddOutline as AddIcon,
  FunnelOutline as FilterIcon,
  SearchOutline as SearchIcon,
  MenuOutline as MenuIcon,
} from '@vicons/ionicons5'
import TaskItem from './TaskItem.vue'
import CreateTaskModal from './CreateTaskModal.vue'

interface Task {
  id: string
  title: string
  description?: string
  completed: boolean
  priority: 'low' | 'medium' | 'high'
  dueDate?: string
  createdAt: string
  updatedAt: string
  projectId: string
}

interface Props {
  projectId: string | null
  tasks: Task[]
  projectListCollapsed: boolean
}

interface Emits {
  (e: 'task-select', task: Task): void
  (e: 'toggle-project-list'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 状态管理
const searchQuery = ref('')
const currentFilter = ref<'all' | 'pending' | 'completed'>('all')
const selectedTaskId = ref<string | null>(null)
const showCreateTask = ref(false)
const currentView = ref<'list' | 'board'>('list')

// 视图选项
const viewOptions = [
  { label: '列表视图', key: 'list' },
  { label: '看板视图', key: 'board' },
  { type: 'divider', key: 'divider' },
  { label: '按优先级排序', key: 'sort-priority' },
  { label: '按日期排序', key: 'sort-date' },
]

// 计算属性
const completedCount = computed(() => {
  return props.tasks.filter((task) => task.completed).length
})

const filteredTasks = computed(() => {
  let filtered = props.tasks

  // 按完成状态筛选
  if (currentFilter.value === 'pending') {
    filtered = filtered.filter((task) => !task.completed)
  } else if (currentFilter.value === 'completed') {
    filtered = filtered.filter((task) => task.completed)
  }

  // 按搜索关键词筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(
      (task) =>
        task.title.toLowerCase().includes(query) || task.description?.toLowerCase().includes(query),
    )
  }

  return filtered
})

const groupedTasks = computed(() => {
  const groups = new Map<string, Task[]>()

  filteredTasks.value.forEach((task) => {
    const date = task.dueDate || task.createdAt
    const dateKey = new Date(date).toDateString()

    if (!groups.has(dateKey)) {
      groups.set(dateKey, [])
    }
    groups.get(dateKey)!.push(task)
  })

  return Array.from(groups.entries())
    .map(([date, tasks]) => ({
      date,
      dateLabel: formatDateLabel(date),
      tasks: tasks.sort((a, b) => {
        // 按优先级和创建时间排序
        const priorityOrder = { high: 3, medium: 2, low: 1 }
        if (a.priority !== b.priority) {
          return priorityOrder[b.priority] - priorityOrder[a.priority]
        }
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      }),
    }))
    .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())
})

// 工具函数
const formatDateLabel = (dateString: string) => {
  const date = new Date(dateString)
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)

  if (date.toDateString() === today.toDateString()) {
    return '今天'
  } else if (date.toDateString() === tomorrow.toDateString()) {
    return '明天'
  } else {
    return date.toLocaleDateString('zh-CN', {
      month: 'long',
      day: 'numeric',
      weekday: 'short',
    })
  }
}

// 事件处理
const handleTaskSelect = (task: Task) => {
  selectedTaskId.value = task.id
  emit('task-select', task)
}

const handleToggleComplete = (taskId: string) => {
  // TODO: 实现任务完成状态切换
  console.log('Toggle complete:', taskId)
}

const handleTaskUpdate = (task: Task) => {
  // TODO: 实现任务更新
  console.log('Update task:', task)
}

const handleTaskCreated = (task: Task) => {
  // TODO: 添加新任务到列表
  console.log('Task created:', task)
}

const handleViewChange = (key: string) => {
  if (key === 'list' || key === 'board') {
    currentView.value = key
  }
  // TODO: 实现其他视图操作
}

// 监听项目变化
watch(
  () => props.projectId,
  () => {
    selectedTaskId.value = null
    searchQuery.value = ''
    currentFilter.value = 'all'
  },
)
</script>

<style scoped>
.task-list-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--n-color);
}

.task-header {
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--n-border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--n-text-color);
}

.header-actions {
  display: flex;
  gap: 4px;
}

.filter-section {
  padding: 12px 16px;
  border-bottom: 1px solid var(--n-border-color);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.filter-tabs {
  display: flex;
  justify-content: center;
}

.task-list-content {
  flex: 1;
  overflow: hidden;
}

.empty-state {
  padding: 40px 20px;
  text-align: center;
}

.task-items {
  padding: 8px 0;
}

.task-group {
  margin-bottom: 16px;
}

.group-header {
  padding: 8px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--n-modal-color);
  border-bottom: 1px solid var(--n-border-color);
  position: sticky;
  top: 0;
  z-index: 1;
}

.group-date {
  font-size: 14px;
  font-weight: 600;
  color: var(--n-text-color);
}

.group-count {
  font-size: 12px;
  color: var(--n-text-color-2);
}
</style>
