<template>
  <div class="project-item" :class="{ active, collapsed }" @click="$emit('select', project)">
    <!-- 项目图标/颜色标识 -->
    <div class="project-indicator">
      <div class="project-color" :style="{ backgroundColor: project.color }"></div>
    </div>

    <!-- 项目信息 -->
    <div class="project-info" v-if="!collapsed">
      <div class="project-name">{{ project.name }}</div>
      <div class="project-meta">
        <span class="task-count">{{ project.completedCount }}/{{ project.taskCount }}</span>
        <div class="progress-bar">
          <div
            class="progress-fill"
            :style="{ width: progressPercentage + '%', backgroundColor: project.color }"
          ></div>
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="project-actions" v-if="!collapsed">
      <n-tooltip>
        <template #trigger>
          <n-button
            quaternary
            circle
            size="tiny"
            @click.stop="$emit('favorite', project.id)"
            class="favorite-btn"
          >
            <template #icon>
              <n-icon :color="project.isFavorite ? '#f0a020' : undefined">
                <component :is="project.isFavorite ? StarFilledIcon : StarIcon" />
              </n-icon>
            </template>
          </n-button>
        </template>
        {{ project.isFavorite ? '取消收藏' : '收藏项目' }}
      </n-tooltip>

      <n-dropdown :options="dropdownOptions" @select="handleDropdownSelect" trigger="click">
        <n-button quaternary circle size="tiny" @click.stop class="more-btn">
          <template #icon>
            <n-icon>
              <MoreIcon />
            </n-icon>
          </template>
        </n-button>
      </n-dropdown>
    </div>

    <!-- 收起状态的提示 -->
    <n-tooltip v-if="collapsed" placement="right">
      <template #trigger>
        <div class="collapsed-indicator"></div>
      </template>
      <div>
        <div class="tooltip-title">{{ project.name }}</div>
        <div class="tooltip-meta">{{ project.completedCount }}/{{ project.taskCount }} 任务</div>
      </div>
    </n-tooltip>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  StarOutline as StarIcon,
  Star as StarFilledIcon,
  EllipsisHorizontalOutline as MoreIcon,
  CreateOutline as EditIcon,
  TrashOutline as DeleteIcon,
  DuplicateOutline as DuplicateIcon,
} from '@vicons/ionicons5'

interface Project {
  id: string
  name: string
  description?: string
  color: string
  taskCount: number
  completedCount: number
  isFavorite: boolean
  createdAt: string
}

interface Props {
  project: Project
  collapsed: boolean
  active: boolean
}

interface Emits {
  (e: 'select', project: Project): void
  (e: 'favorite', projectId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 计算属性
const progressPercentage = computed(() => {
  if (props.project.taskCount === 0) return 0
  return Math.round((props.project.completedCount / props.project.taskCount) * 100)
})

const dropdownOptions = [
  {
    label: '编辑项目',
    key: 'edit',
  },
  {
    label: '复制项目',
    key: 'duplicate',
  },
  {
    type: 'divider',
    key: 'divider',
  },
  {
    label: '删除项目',
    key: 'delete',
  },
]

// 事件处理
const handleDropdownSelect = (key: string) => {
  switch (key) {
    case 'edit':
      // TODO: 打开编辑对话框
      break
    case 'duplicate':
      // TODO: 复制项目
      break
    case 'delete':
      // TODO: 删除项目
      break
  }
}
</script>

<style scoped>
.project-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  margin: 2px 0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.project-item:hover {
  background: var(--n-button-color-hover);
}

.project-item.active {
  background: var(--n-button-color-pressed);
  color: var(--n-button-text-color-pressed);
}

.project-item.collapsed {
  justify-content: center;
  padding: 12px 8px;
}

.project-indicator {
  display: flex;
  align-items: center;
  margin-right: 12px;
}

.project-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.project-info {
  flex: 1;
  min-width: 0;
}

.project-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--n-text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.project-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-count {
  font-size: 12px;
  color: var(--n-text-color-2);
  white-space: nowrap;
}

.progress-bar {
  flex: 1;
  height: 3px;
  background: var(--n-border-color);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.3s ease;
}

.project-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.project-item:hover .project-actions {
  opacity: 1;
}

.favorite-btn,
.more-btn {
  flex-shrink: 0;
}

.collapsed-indicator {
  position: absolute;
  inset: 0;
}

.tooltip-title {
  font-weight: 500;
  margin-bottom: 2px;
}

.tooltip-meta {
  font-size: 12px;
  color: var(--n-text-color-2);
}
</style>
