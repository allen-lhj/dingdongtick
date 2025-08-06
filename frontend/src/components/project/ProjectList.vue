<template>
  <div class="project-list-container">
    <!-- 头部操作栏 -->
    <div class="project-header">
      <div class="header-title">
        <n-icon size="18">
          <FolderIcon />
        </n-icon>
        <span v-if="!collapsed">项目</span>
      </div>

      <div class="header-actions" v-if="!collapsed">
        <n-tooltip>
          <template #trigger>
            <n-button quaternary circle size="small" @click="showCreateProject = true">
              <template #icon>
                <n-icon>
                  <AddIcon />
                </n-icon>
              </template>
            </n-button>
          </template>
          新建项目
        </n-tooltip>
      </div>
    </div>

    <!-- 搜索框 -->
    <div class="search-section" v-if="!collapsed">
      <n-input v-model:value="searchQuery" placeholder="搜索项目..." clearable size="small">
        <template #prefix>
          <n-icon>
            <SearchIcon />
          </n-icon>
        </template>
      </n-input>
    </div>

    <!-- 项目列表 -->
    <div class="project-list-content">
      <n-scrollbar>
        <!-- 收藏项目 -->
        <div class="project-section" v-if="favoriteProjects.length > 0">
          <div class="section-title" v-if="!collapsed">
            <n-icon size="14">
              <StarIcon />
            </n-icon>
            收藏
          </div>
          <div class="project-items">
            <ProjectItem
              v-for="project in favoriteProjects"
              :key="project.id"
              :project="project"
              :collapsed="collapsed"
              :active="selectedProjectId === project.id"
              @select="handleProjectSelect"
              @favorite="handleToggleFavorite"
            />
          </div>
        </div>

        <!-- 所有项目 -->
        <div class="project-section">
          <div class="section-title" v-if="!collapsed">
            <n-icon size="14">
              <FolderIcon />
            </n-icon>
            所有项目
          </div>
          <div class="project-items">
            <ProjectItem
              v-for="project in filteredProjects"
              :key="project.id"
              :project="project"
              :collapsed="collapsed"
              :active="selectedProjectId === project.id"
              @select="handleProjectSelect"
              @favorite="handleToggleFavorite"
            />
          </div>
        </div>
      </n-scrollbar>
    </div>

    <!-- 创建项目对话框 -->
    <CreateProjectModal v-model:show="showCreateProject" @created="handleProjectCreated" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  FolderOutline as FolderIcon,
  AddOutline as AddIcon,
  SearchOutline as SearchIcon,
  StarOutline as StarIcon,
} from '@vicons/ionicons5'
import ProjectItem from './ProjectItem.vue'
import CreateProjectModal from './CreateProjectModal.vue'

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
  collapsed: boolean
}

interface Emits {
  (e: 'project-select', projectId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 状态管理
const searchQuery = ref('')
const selectedProjectId = ref<string | null>(null)
const showCreateProject = ref(false)
const projects = ref<Project[]>([])

// 计算属性
const filteredProjects = computed(() => {
  if (!searchQuery.value) return projects.value.filter((p) => !p.isFavorite)

  return projects.value
    .filter((p) => !p.isFavorite)
    .filter(
      (project) =>
        project.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        project.description?.toLowerCase().includes(searchQuery.value.toLowerCase()),
    )
})

const favoriteProjects = computed(() => {
  return projects.value.filter((p) => p.isFavorite)
})

// 事件处理
const handleProjectSelect = (project: Project) => {
  selectedProjectId.value = project.id
  emit('project-select', project.id)
}

const handleToggleFavorite = (projectId: string) => {
  const project = projects.value.find((p) => p.id === projectId)
  if (project) {
    project.isFavorite = !project.isFavorite
  }
}

const handleProjectCreated = (project: Project) => {
  projects.value.unshift(project)
}

// 数据加载
const loadProjects = async () => {
  // TODO: 实现项目数据加载
  projects.value = [
    {
      id: '1',
      name: '个人任务',
      description: '个人日常任务管理',
      color: '#18a058',
      taskCount: 12,
      completedCount: 8,
      isFavorite: true,
      createdAt: '2024-01-01',
    },
    {
      id: '2',
      name: '工作项目',
      description: '工作相关任务',
      color: '#2080f0',
      taskCount: 25,
      completedCount: 15,
      isFavorite: false,
      createdAt: '2024-01-02',
    },
    {
      id: '3',
      name: '学习计划',
      description: '技能提升和学习',
      color: '#f0a020',
      taskCount: 8,
      completedCount: 3,
      isFavorite: true,
      createdAt: '2024-01-03',
    },
  ]
}

onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-list-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--n-color);
}

.project-header {
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--n-border-color);
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

.search-section {
  padding: 12px 16px;
  border-bottom: 1px solid var(--n-border-color);
}

.project-list-content {
  flex: 1;
  overflow: hidden;
}

.project-section {
  padding: 8px 0;
}

.section-title {
  padding: 8px 16px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--n-text-color-2);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.project-items {
  padding: 0 8px;
}
</style>
