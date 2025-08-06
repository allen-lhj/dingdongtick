<template>
  <div class="task-detail-container" :class="{ visible }">
    <div v-if="!task" class="empty-detail">
      <n-empty description="选择一个任务查看详情">
        <template #icon>
          <n-icon size="48">
            <DocumentIcon />
          </n-icon>
        </template>
      </n-empty>
    </div>

    <div v-else class="task-detail-content">
      <!-- 头部操作栏 -->
      <div class="detail-header">
        <div class="header-left">
          <n-button
            quaternary
            circle
            size="small"
            @click="$emit('close')"
          >
            <template #icon>
              <n-icon>
                <CloseIcon />
              </n-icon>
            </template>
          </n-button>
          <span class="header-title">任务详情</span>
        </div>
        
        <div class="header-actions">
          <n-tooltip>
            <template #trigger>
              <n-button
                quaternary
                circle
                size="small"
                @click="isEditing = !isEditing"
              >
                <template #icon>
                  <n-icon>
                    <EditIcon />
                  </n-icon>
                </template>
              </n-button>
            </template>
            {{ isEditing ? '取消编辑' : '编辑任务' }}
          </n-tooltip>
          
          <n-dropdown
            :options="moreOptions"
            @select="handleMoreAction"
            trigger="click"
          >
            <n-button quaternary circle size="small">
              <template #icon>
                <n-icon>
                  <MoreIcon />
                </n-icon>
              </template>
            </n-button>
          </n-dropdown>
        </div>
      </div>

      <!-- 任务内容 -->
      <div class="detail-content">
        <n-scrollbar>
          <!-- 任务标题和状态 -->
          <div class="task-title-section">
            <div class="title-row">
              <n-checkbox
                :checked="task.completed"
                @update:checked="handleToggleComplete"
                size="large"
              />
              <n-input
                v-if="isEditing"
                v-model:value="editForm.title"
                placeholder="任务标题"
                size="large"
                class="title-input"
              />
              <h2 v-else class="task-title" :class="{ completed: task.completed }">
                {{ task.title }}
              </h2>
            </div>
            
            <div class="task-meta">
              <n-tag
                :type="getPriorityType(task.priority)"
                size="small"
              >
                {{ getPriorityLabel(task.priority) }}
              </n-tag>
              <span class="created-time">
                创建于 {{ formatDate(task.createdAt) }}
              </span>
            </div>
          </div>

          <!-- 任务描述 -->
          <div class="task-section">
            <div class="section-title">
              <n-icon>
                <DocumentTextIcon />
              </n-icon>
              描述
            </div>
            <n-input
              v-if="isEditing"
              v-model:value="editForm.description"
              type="textarea"
              placeholder="添加任务描述..."
              :rows="4"
            />
            <div v-else class="task-description">
              {{ task.description || '暂无描述' }}
            </div>
          </div>

          <!-- 任务属性 -->
          <div class="task-section">
            <div class="section-title">
              <n-icon>
                <SettingsIcon />
              </n-icon>
              属性
            </div>
            
            <div class="task-properties">
              <!-- 截止日期 -->
              <div class="property-item">
                <span class="property-label">截止日期</span>
                <n-date-picker
                  v-if="isEditing"
                  v-model:value="editForm.dueDate"
                  type="datetime"
                  clearable
                  size="small"
                />
                <span v-else class="property-value">
                  {{ task.dueDate ? formatDate(task.dueDate) : '未设置' }}
                </span>
              </div>

              <!-- 优先级 -->
              <div class="property-item">
                <span class="property-label">优先级</span>
                <n-select
                  v-if="isEditing"
                  v-model:value="editForm.priority"
                  :options="priorityOptions"
                  size="small"
                />
                <n-tag
                  v-else
                  :type="getPriorityType(task.priority)"
                  size="small"
                >
                  {{ getPriorityLabel(task.priority) }}
                </n-tag>
              </div>

              <!-- 标签 -->
              <div class="property-item">
                <span class="property-label">标签</span>
                <n-dynamic-tags
                  v-if="isEditing"
                  v-model:value="editForm.tags"
                  size="small"
                />
                <div v-else class="task-tags">
                  <n-tag
                    v-for="tag in task.tags || []"
                    :key="tag"
                    size="small"
                    closable
                  >
                    {{ tag }}
                  </n-tag>
                </div>
              </div>
            </div>
          </div>

          <!-- 子任务 -->
          <div class="task-section">
            <div class="section-title">
              <n-icon>
                <ListIcon />
              </n-icon>
              子任务
              <n-button
                v-if="isEditing"
                quaternary
                size="tiny"
                @click="addSubtask"
              >
                <template #icon>
                  <n-icon>
                    <AddIcon />
                  </n-icon>
                </template>
              </n-button>
            </div>
            
            <div class="subtasks">
              <div
                v-for="(subtask, index) in editForm.subtasks"
                :key="index"
                class="subtask-item"
              >
                <n-checkbox
                  v-model:checked="subtask.completed"
                  size="small"
                />
                <n-input
                  v-if="isEditing"
                  v-model:value="subtask.title"
                  placeholder="子任务标题"
                  size="small"
                  class="subtask-input"
                />
                <span v-else class="subtask-title" :class="{ completed: subtask.completed }">
                  {{ subtask.title }}
                </span>
                <n-button
                  v-if="isEditing"
                  quaternary
                  circle
                  size="tiny"
                  @click="removeSubtask(index)"
                >
                  <template #icon>
                    <n-icon>
                      <CloseIcon />
                    </n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div v-if="isEditing" class="action-buttons">
            <n-button type="primary" @click="handleSave">
              保存更改
            </n-button>
            <n-button @click="handleCancel">
              取消
            </n-button>
          </div>
        </n-scrollbar>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import {
  DocumentOutline as DocumentIcon,
  DocumentTextOutline as DocumentTextIcon,
  CloseOutline as CloseIcon,
  CreateOutline as EditIcon,
  EllipsisHorizontalOutline as MoreIcon,
  SettingsOutline as SettingsIcon,
  ListOutline as ListIcon,
  AddOutline as AddIcon,
  TrashOutline as DeleteIcon,
  DuplicateOutline as DuplicateIcon
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
  subtasks?: Array<{ title: string; completed: boolean }>
}

interface Props {
  task: Task | null
  visible: boolean
}

interface Emits {
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 状态管理
const isEditing = ref(false)
const editForm = reactive({
  title: '',
  description: '',
  priority: 'medium' as const,
  dueDate: null as number | null,
  tags: [] as string[],
  subtasks: [] as Array<{ title: string; completed: boolean }>
})

// 选项配置
const priorityOptions = [
  { label: '低优先级', value: 'low' },
  { label: '中优先级', value: 'medium' },
  { label: '高优先级', value: 'high' }
]

const moreOptions = [
  { label: '复制任务', key: 'duplicate', icon: DuplicateIcon },
  { type: 'divider', key: 'divider' },
  { label: '删除任务', key: 'delete', icon: DeleteIcon }
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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 事件处理
const handleToggleComplete = (completed: boolean) => {
  // TODO: 实现任务完成状态切换
  console.log('Toggle complete:', completed)
}

const handleMoreAction = (key: string) => {
  switch (key) {
    case 'duplicate':
      // TODO: 复制任务
      break
    case 'delete':
      // TODO: 删除任务
      break
  }
}

const addSubtask = () => {
  editForm.subtasks.push({ title: '', completed: false })
}

const removeSubtask = (index: number) => {
  editForm.subtasks.splice(index, 1)
}

const handleSave = () => {
  // TODO: 保存任务更改
  console.log('Save task:', editForm)
  isEditing.value = false
}

const handleCancel = () => {
  // 重置表单
  if (props.task) {
    Object.assign(editForm, {
      title: props.task.title,
      description: props.task.description || '',
      priority: props.task.priority,
      dueDate: props.task.dueDate ? new Date(props.task.dueDate).getTime() : null,
      tags: [...(props.task.tags || [])],
      subtasks: [...(props.task.subtasks || [])]
    })
  }
  isEditing.value = false
}

// 监听任务变化
watch(() => props.task, (newTask) => {
  if (newTask) {
    Object.assign(editForm, {
      title: newTask.title,
      description: newTask.description || '',
      priority: newTask.priority,
      dueDate: newTask.dueDate ? new Date(newTask.dueDate).getTime() : null,
      tags: [...(newTask.tags || [])],
      subtasks: [...(newTask.subtasks || [])]
    })
  }
  isEditing.value = false
}, { immediate: true })
</script>

<style scoped>
.task-detail-container {
  height: 100%;
  background: var(--n-color);
  transition: all 0.3s ease;
}

.empty-detail {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.task-detail-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.detail-header {
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--n-border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-title {
  font-weight: 600;
  color: var(--n-text-color);
}

.header-actions {
  display: flex;
  gap: 4px;
}

.detail-content {
  flex: 1;
  overflow: hidden;
}

.task-title-section {
  padding: 20px;
  border-bottom: 1px solid var(--n-border-color);
}

.title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.task-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--n-text-color);
  flex: 1;
}

.task-title.completed {
  text-decoration: line-through;
  color: var(--n-text-color-2);
}

.title-input {
  flex: 1;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.created-time {
  font-size: 12px;
  color: var(--n-text-color-2);
}

.task-section {
  padding: 20px;
  border-bottom: 1px solid var(--n-border-color);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--n-text-color);
  margin-bottom: 12px;
}

.task-description {
  color: var(--n-text-color-2);
  line-height: 1.6;
  white-space: pre-wrap;
}

.task-properties {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.property-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.property-label {
  font-weight: 500;
  color: var(--n-text-color);
  min-width: 80px;
}

.property-value {
  color: var(--n-text-color-2);
}

.task-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.subtasks {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.subtask-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.subtask-title {
  flex: 1;
  color: var(--n-text-color);
}

.subtask-title.completed {
  text-decoration: line-through;
  color: var(--n-text-color-2);
}

.subtask-input {
  flex: 1;
}

.action-buttons {
  padding: 20px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
</style>
