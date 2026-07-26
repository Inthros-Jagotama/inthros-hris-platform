import { reactive } from 'vue'

const LANG_KEY = 'hris_lang'

function detectInitialLang() {
  // 1. Cek localStorage
  const stored = localStorage.getItem(LANG_KEY)
  if (stored === 'id' || stored === 'en') return stored
  // 2. Cek browser language
  const navLang = navigator.language || navigator.userLanguage || ''
  if (navLang.startsWith('id') || navLang.startsWith('in')) return 'id'
  // 3. Default English
  return 'en'
}

const state = reactive({
  lang: detectInitialLang()
})

export function useLanguage() {
  function setLang(lang) {
    if (lang !== 'en' && lang !== 'id') return
    state.lang = lang
    localStorage.setItem(LANG_KEY, lang)
    // Update <html> lang attribute untuk aksesibilitas
    document.documentElement.lang = lang === 'id' ? 'id' : 'en'
  }

  function toggleLang() {
    setLang(state.lang === 'en' ? 'id' : 'en')
  }

  function isID() {
    return state.lang === 'id'
  }

  function isEN() {
    return state.lang === 'en'
  }

  return { state, setLang, toggleLang, isID, isEN }
}
