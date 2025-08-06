<template>
  <n-modal
    :show="show"
    @update:show="$emit('update:show', $event)"
    preset="dialog"
    title="创建新任务"
    style="width: 600px"
  >
    <n-form ref="formRef" :model="form" :rules="rules">
      <n-form-item label="任务标题" path="title">
        <n-input v-model:value="form.title" placeholder="输入任务标题" />
      </n-form-item>

      <n-form-item label="任务描述" path="description">
        <n-input
          v-model:value="form.description"
          type="textarea"
          placeholder="输入任务描述（可选）"
          :rows="3"
        />
      </n-form-item>

      <n-form-item label="优先级" path="priority">
        <n-select
          v-model:value="form.priority"
          :options="priorityOptions"
          placeholder="选择优先级"
        />
      </n-form-item>

      <n-form-item label="截止日期" path="dueDate">
        <n-date-picker
          v-model:value="form.dueDate"
          type="datetime"
          clearable
          placeholder="选择截止日期（可选）"
        />
      </n-form-item>

      <n-form-item label="标签" path="tags">
        <n-dynamic-tags v-model:value="form.tags" />
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
  projectId: string | null
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'created', task: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref()
const form = reactive({
  title: '',
  description: '',
  priority: 'medium',
  dueDate: null as number | null,
  tags: [] as string[],
})

const priorityOptions = [
  { label: '低优先级', value: 'low' },
  { label: '中优先级', value: 'medium' },
  { label: '高优先级', value: 'high' },
]

const rules = {
  title: {
    required: true,
    message: '请输入任务标题',
    trigger: 'blur',
  },
}

const handleCreate = async () => {
  try {
    await formRef.value?.validate()
    const task = {
      id: Date.now().toString(),
      title: form.title,
      description: form.description,
      completed: false,
      priority: form.priority,
      dueDate: form.dueDate ? new Date(form.dueDate).toISOString() : undefined,
      tags: form.tags,
      projectId: props.projectId || '',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      subtasks: [],
    }
    emit('created', task)
    emit('update:show', false)
    // 重置表单
    Object.assign(form, {
      title: '',
      description: '',
      priority: 'medium',
      dueDate: null,
      tags: [],
    })
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const handleCancel = () => {
  emit('update:show', false)
}
</script>
