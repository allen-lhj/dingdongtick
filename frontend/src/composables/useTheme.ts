import { ref, computed, watch } from 'vue'
import { darkTheme, type GlobalTheme } from 'naive-ui'

// 主题状态
const isDark = ref(false)

// 从 localStorage 读取主题设置
const loadThemeFromStorage = () => {
  const stored = localStorage.getItem('theme-preference')
  if (stored) {
    isDark.value = stored === 'dark'
  } else {
    // 如果没有存储的偏好，检查系统偏好
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
}

// 保存主题设置到 localStorage
const saveThemeToStorage = (theme: string) => {
  localStorage.setItem('theme-preference', theme)
}

// 计算当前主题对象
const theme = computed<GlobalTheme | null>(() => {
  return isDark.value ? darkTheme : null
})

// 切换主题
const toggleTheme = () => {
  isDark.value = !isDark.value
  saveThemeToStorage(isDark.value ? 'dark' : 'light')
}

// 设置特定主题
const setTheme = (dark: boolean) => {
  isDark.value = dark
  saveThemeToStorage(dark ? 'dark' : 'light')
}

// 监听系统主题变化
const watchSystemTheme = () => {
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  
  const handleChange = (e: MediaQueryListEvent) => {
    // 只有在没有用户偏好设置时才跟随系统
    const stored = localStorage.getItem('theme-preference')
    if (!stored) {
      isDark.value = e.matches
    }
  }
  
  mediaQuery.addEventListener('change', handleChange)
  
  // 返回清理函数
  return () => {
    mediaQuery.removeEventListener('change', handleChange)
  }
}

// 初始化主题
const initTheme = () => {
  loadThemeFromStorage()
  return watchSystemTheme()
}

// 导出主题管理功能
export const useTheme = () => {
  return {
    isDark: computed(() => isDark.value),
    theme,
    toggleTheme,
    setTheme,
    initTheme
  }
}

// 导出全局状态（用于跨组件共享）
export { isDark, theme, toggleTheme, setTheme, initTheme }
