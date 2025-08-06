<template>
  <n-modal
    :show="show"
    @update:show="$emit('update:show', $event)"
    preset="dialog"
    title="创建新项目"
  >
    <n-form ref="formRef" :model="form" :rules="rules">
      <n-form-item label="项目名称" path="name">
        <n-input v-model:value="form.name" placeholder="输入项目名称" />
      </n-form-item>

      <n-form-item label="项目描述" path="description">
        <n-input
          v-model:value="form.description"
          type="textarea"
          placeholder="输入项目描述（可选）"
          :rows="3"
        />
      </n-form-item>

      <n-form-item label="项目颜色" path="color">
        <div class="color-picker">
          <div
            v-for="color in colorOptions"
            :key="color"
            class="color-option"
            :class="{ active: form.color === color }"
            :style="{ backgroundColor: color }"
            @click="form.color = color"
          ></div>
        </div>
      </n-form-item>
    </n-form>

    <template #action>
      <n-button @click="handleCancel">取消</n-button>
      <n-button type="primary" @click="handleCreate">创建</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'

interface Props {
  show: boolean
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'created', project: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref()
const form = reactive({
  name: '',
  description: '',
  color: '#18a058',
})

const colorOptions = [
  '#18a058',
  '#2080f0',
  '#f0a020',
  '#d03050',
  '#7c3aed',
  '#059669',
  '#dc2626',
  '#ea580c',
]

const rules = {
  name: {
    required: true,
    message: '请输入项目名称',
    trigger: 'blur',
  },
}

const handleCreate = async () => {
  try {
    await formRef.value?.validate()
    const project = {
      id: Date.now().toString(),
      ...form,
      taskCount: 0,
      completedCount: 0,
      isFavorite: false,
      createdAt: new Date().toISOString(),
    }
    emit('created', project)
    emit('update:show', false)
    // 重置表单
    Object.assign(form, {
      name: '',
      description: '',
      color: '#18a058',
    })
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const handleCancel = () => {
  emit('update:show', false)
}
</script>

<style scoped>
.color-picker {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.color-option {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  border-color: var(--n-color-primary);
  transform: scale(1.1);
}
</style>
