<template>
  <div
    class="task-item"
    :class="{ active, completed: task.completed }"
    @click="$emit('select', task)"
  >
    <!-- 任务完成状态 -->
    <n-checkbox
      :checked="task.completed"
      @update:checked="$emit('toggle-complete', task.id)"
      @click.stop
      size="medium"
    />

    <!-- 任务内容 -->
    <div class="task-content">
      <div class="task-title" :class="{ completed: task.completed }">
        {{ task.title }}
      </div>
      <div v-if="task.description" class="task-description">
        {{ task.description }}
      </div>

      <!-- 任务元信息 -->
      <div class="task-meta">
        <n-tag :type="getPriorityType(task.priority)" size="small" class="priority-tag">
          {{ getPriorityLabel(task.priority) }}
        </n-tag>

        <span v-if="task.dueDate" class="due-date" :class="{ overdue: isOverdue }">
          <n-icon size="12">
            <TimeIcon />
          </n-icon>
          {{ formatDueDate(task.dueDate) }}
        </span>

        <div v-if="task.tags && task.tags.length > 0" class="task-tags">
          <n-tag v-for="tag in task.tags.slice(0, 2)" :key="tag" size="tiny">
            {{ tag }}
          </n-tag>
          <span v-if="task.tags.length > 2" class="more-tags"> +{{ task.tags.length - 2 }} </span>
        </div>
      </div>
    </div>

    <!-- 任务操作 -->
    <div class="task-actions">
      <n-dropdown :options="actionOptions" @select="handleAction" trigger="click">
        <n-button quaternary circle size="small" @click.stop class="action-btn">
          <template #icon>
            <n-icon>
              <MoreIcon />
            </n-icon>
          </template>
        </n-button>
      </n-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  TimeOutline as TimeIcon,
  EllipsisHorizontalOutline as MoreIcon,
  CreateOutline as EditIcon,
  TrashOutline as DeleteIcon,
  DuplicateOutline as DuplicateIcon,
} from '@vicons/ionicons5'

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
  tags?: string[]
}

interface Props {
  task: Task
  active: boolean
}

interface Emits {
  (e: 'select', task: Task): void
  (e: 'toggle-complete', taskId: string): void
  (e: 'update', task: Task): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 计算属性
const isOverdue = computed(() => {
  if (!props.task.dueDate) return false
  return new Date(props.task.dueDate) < new Date() && !props.task.completed
})

const actionOptions = [
  {
    label: '编辑',
    key: 'edit',
  },
  {
    label: '复制',
    key: 'duplicate',
  },
  {
    type: 'divider',
    key: 'divider',
  },
  {
    label: '删除',
    key: 'delete',
  },
]

// 工具函数
const getPriorityType = (priority: string) => {
  const types = { low: 'default', medium: 'warning', high: 'error' }
  return types[priority as keyof typeof types] || 'default'
}

const getPriorityLabel = (priority: string) => {
  const labels = { low: '低', medium: '中', high: '高' }
  return labels[priority as keyof typeof labels] || '中'
}

const formatDueDate = (dateString: string) => {
  const date = new Date(dateString)
  const today = new Date()
  const diffTime = date.getTime() - today.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return '今天'
  if (diffDays === 1) return '明天'
  if (diffDays === -1) return '昨天'
  if (diffDays < 0) return `逾期 ${Math.abs(diffDays)} 天`
  if (diffDays <= 7) return `${diffDays} 天后`

  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  })
}

// 事件处理
const handleAction = (key: string) => {
  switch (key) {
    case 'edit':
      // TODO: 打开编辑对话框
      break
    case 'duplicate':
      // TODO: 复制任务
      break
    case 'delete':
      // TODO: 删除任务
      break
  }
}
</script>

<style scoped>
.task-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  margin: 2px 0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.task-item:hover {
  background: var(--n-button-color-hover);
}

.task-item.active {
  background: var(--n-button-color-pressed);
  border-color: var(--n-color-primary);
}

.task-item.completed {
  opacity: 0.6;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--n-text-color);
  line-height: 1.4;
  margin-bottom: 4px;
}

.task-title.completed {
  text-decoration: line-through;
  color: var(--n-text-color-2);
}

.task-description {
  font-size: 12px;
  color: var(--n-text-color-2);
  line-height: 1.4;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.priority-tag {
  flex-shrink: 0;
}

.due-date {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--n-text-color-2);
  flex-shrink: 0;
}

.due-date.overdue {
  color: var(--n-color-error);
}

.task-tags {
  display: flex;
  align-items: center;
  gap: 4px;
}

.more-tags {
  font-size: 11px;
  color: var(--n-text-color-3);
}

.task-actions {
  opacity: 0;
  transition: opacity 0.2s ease;
}

.task-item:hover .task-actions {
  opacity: 1;
}

.action-btn {
  flex-shrink: 0;
}
</style>
