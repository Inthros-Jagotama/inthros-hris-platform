import { reactive } from 'vue'

const THEME_KEY = 'hris_theme'

function detectInitialTheme() {
  // 1. Cek localStorage
  const stored = localStorage.getItem(THEME_KEY)
  if (stored === 'dark' || stored === 'light') return stored
  // 2. Cek system preference
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark'
  // 3. Default light
  return 'light'
}

function applyTheme(theme) {
  const root = document.documentElement
  if (theme === 'dark') {
    root.classList.add('dark')
    root.classList.add('p-dark')
  } else {
    root.classList.remove('dark')
    root.classList.remove('p-dark')
  }
}

// Apply on load immediately (before Vue mounts)
applyTheme(detectInitialTheme())

const state = reactive({
  theme: detectInitialTheme()
})

export function useTheme() {
  function setTheme(theme) {
    if (theme !== 'dark' && theme !== 'light') return
    state.theme = theme
    localStorage.setItem(THEME_KEY, theme)
    applyTheme(theme)
  }

  function toggleTheme() {
    setTheme(state.theme === 'dark' ? 'light' : 'dark')
  }

  function isDark() {
    return state.theme === 'dark'
  }

  function isLight() {
    return state.theme === 'light'
  }

  return { state, setTheme, toggleTheme, isDark, isLight }
}
