<template>
  <div class="layout-demo">
    <h1>Naive UI 主页布局设计方案</h1>
    
    <div class="demo-description">
      <h2>设计概述</h2>
      <p>基于您的需求，我设计了一个四区域的响应式布局：</p>
      
      <div class="layout-features">
        <div class="feature-item">
          <h3>1. 固定侧边栏 (MainSidebar)</h3>
          <ul>
            <li>用户头像和信息展示</li>
            <li>主要导航菜单（首页、今天、项目、任务等）</li>
            <li>主题切换和设置按钮</li>
            <li>支持展开/收起状态</li>
          </ul>
        </div>
        
        <div class="feature-item">
          <h3>2. 项目列表 (ProjectList)</h3>
          <ul>
            <li>可展开/收起的项目导航面板</li>
            <li>项目搜索功能</li>
            <li>收藏项目和所有项目分组显示</li>
            <li>项目进度可视化</li>
            <li>项目操作菜单（编辑、删除、收藏等）</li>
          </ul>
        </div>
        
        <div class="feature-item">
          <h3>3. 任务列表 (TaskList)</h3>
          <ul>
            <li>根据选中项目动态显示任务</li>
            <li>任务搜索和筛选（全部、待完成、已完成）</li>
            <li>按日期分组显示任务</li>
            <li>任务优先级和截止时间显示</li>
            <li>任务快速操作</li>
          </ul>
        </div>
        
        <div class="feature-item">
          <h3>4. 任务详情 (TaskDetail)</h3>
          <ul>
            <li>选中任务的详细信息展示</li>
            <li>任务编辑功能</li>
            <li>子任务管理</li>
            <li>任务属性设置（优先级、截止时间、标签）</li>
            <li>任务操作菜单</li>
          </ul>
        </div>
      </div>
    </div>
    
    <div class="technical-details">
      <h2>技术实现特点</h2>
      
      <div class="tech-features">
        <div class="tech-item">
          <h4>响应式设计</h4>
          <p>使用 Naive UI 的 n-layout 组件系列，支持移动端自适应</p>
        </div>
        
        <div class="tech-item">
          <h4>可调整面板</h4>
          <p>使用 n-split 组件实现任务列表和详情面板的可调整分割</p>
        </div>
        
        <div class="tech-item">
          <h4>状态管理</h4>
          <p>组件间通过 props 和 emit 进行通信，支持 Pinia 状态管理</p>
        </div>
        
        <div class="tech-item">
          <h4>主题支持</h4>
          <p>完整支持 Naive UI 的明暗主题切换</p>
        </div>
      </div>
    </div>
    
    <div class="component-structure">
      <h2>组件结构</h2>
      <pre class="code-block">
layouts/
├── MainLayout.vue          # 主布局容器
components/
├── layout/
│   └── MainSidebar.vue     # 主侧边栏
├── project/
│   ├── ProjectList.vue     # 项目列表
│   ├── ProjectItem.vue     # 项目项
│   └── CreateProjectModal.vue # 创建项目对话框
├── task/
│   ├── TaskList.vue        # 任务列表
│   ├── TaskItem.vue        # 任务项
│   ├── TaskDetail.vue      # 任务详情
│   └── CreateTaskModal.vue # 创建任务对话框
└── settings/
    └── SettingsPanel.vue   # 设置面板
      </pre>
    </div>
    
    <div class="next-steps">
      <h2>后续开发建议</h2>
      <ol>
        <li><strong>数据层集成</strong>：集成 API 服务和 Pinia 状态管理</li>
        <li><strong>功能完善</strong>：实现任务的 CRUD 操作和项目管理</li>
        <li><strong>性能优化</strong>：添加虚拟滚动和懒加载</li>
        <li><strong>用户体验</strong>：添加加载状态、错误处理和动画效果</li>
        <li><strong>测试覆盖</strong>：编写单元测试和集成测试</li>
      </ol>
    </div>
  </div>
</template>

<script setup lang="ts">
// 演示页面，无需额外逻辑
</script>

<style scoped>
.layout-demo {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

h1 {
  color: var(--n-text-color);
  text-align: center;
  margin-bottom: 40px;
  font-size: 32px;
  font-weight: 700;
}

h2 {
  color: var(--n-text-color);
  margin: 32px 0 16px 0;
  font-size: 24px;
  font-weight: 600;
  border-bottom: 2px solid var(--n-color-primary);
  padding-bottom: 8px;
}

h3 {
  color: var(--n-text-color);
  margin: 20px 0 12px 0;
  font-size: 18px;
  font-weight: 600;
}

h4 {
  color: var(--n-text-color);
  margin: 16px 0 8px 0;
  font-size: 16px;
  font-weight: 600;
}

p {
  color: var(--n-text-color-2);
  line-height: 1.6;
  margin-bottom: 16px;
}

.demo-description {
  background: var(--n-card-color);
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 32px;
  border: 1px solid var(--n-border-color);
}

.layout-features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.feature-item {
  background: var(--n-modal-color);
  padding: 20px;
  border-radius: 6px;
  border: 1px solid var(--n-border-color);
}

.feature-item ul {
  margin: 12px 0;
  padding-left: 20px;
}

.feature-item li {
  color: var(--n-text-color-2);
  margin-bottom: 8px;
  line-height: 1.5;
}

.technical-details {
  background: var(--n-card-color);
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 32px;
  border: 1px solid var(--n-border-color);
}

.tech-features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
  margin-top: 20px;
}

.tech-item {
  background: var(--n-modal-color);
  padding: 16px;
  border-radius: 6px;
  border: 1px solid var(--n-border-color);
}

.component-structure {
  background: var(--n-card-color);
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 32px;
  border: 1px solid var(--n-border-color);
}

.code-block {
  background: var(--n-code-color);
  color: var(--n-text-color);
  padding: 16px;
  border-radius: 6px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  overflow-x: auto;
  border: 1px solid var(--n-border-color);
}

.next-steps {
  background: var(--n-card-color);
  padding: 24px;
  border-radius: 8px;
  border: 1px solid var(--n-border-color);
}

.next-steps ol {
  margin: 16px 0;
  padding-left: 20px;
}

.next-steps li {
  color: var(--n-text-color-2);
  margin-bottom: 12px;
  line-height: 1.6;
}

.next-steps strong {
  color: var(--n-text-color);
  font-weight: 600;
}

@media (max-width: 768px) {
  .layout-demo {
    padding: 20px 16px;
  }
  
  h1 {
    font-size: 24px;
  }
  
  .layout-features,
  .tech-features {
    grid-template-columns: 1fr;
  }
}
</style>
