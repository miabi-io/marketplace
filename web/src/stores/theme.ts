import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

// Same key and semantics as the Miabi console, so a user who set dark there
// lands in dark here too when both are served from one domain.
export const useThemeStore = defineStore('theme', () => {
  const stored = localStorage.getItem('miabi_theme') as ThemeMode | null
  const mode = ref<ThemeMode>(stored && ['light', 'dark', 'system'].includes(stored) ? stored : 'system')

  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const isDark = computed(() => (mode.value === 'system' ? mediaQuery.matches : mode.value === 'dark'))

  function apply() {
    document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  }

  function setMode(m: ThemeMode) {
    mode.value = m
  }

  function toggle() {
    mode.value = isDark.value ? 'light' : 'dark'
  }

  watch(
    mode,
    (val) => {
      localStorage.setItem('miabi_theme', val)
      apply()
    },
    { immediate: true },
  )

  mediaQuery.addEventListener('change', () => {
    if (mode.value === 'system') apply()
  })

  return { mode, isDark, toggle, setMode }
})
