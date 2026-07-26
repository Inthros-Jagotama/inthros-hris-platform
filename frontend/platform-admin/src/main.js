import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Tooltip from 'primevue/tooltip'
import router from './router'
import App from './App.vue'
import './assets/styles/main.css'
import './assets/styles/slug-animation.css'

const app = createApp(App)

app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      prefix: 'p',
      darkModeSelector: '.p-dark',
      cssLayer: false
    }
  },
  ripple: true,
  inputStyle: 'filled'
})

app.use(router)
app.use(ToastService)
app.use(ConfirmationService)

// Safe Tooltip directive — wraps PrimeVue Tooltip lifecycle
// with retry logic to handle PrimeVue v4 DOM errors (parentNode null, isUnstyled)
// when tooltip elements are dynamically rendered inside DataTable virtual scrolling.
// Retries mount on next tick if initial attempt fails.
// Ref: https://github.com/primefaces/primevue/issues/6478
const SafeTooltip = {
  mounted(el, binding) {
    try {
      // PrimeVue 4 requires element to be fully in DOM before mounting tooltip
      // DataTable virtual scrolling can cause premature mount attempts
      if (!document.body.contains(el)) {
        // Element not yet in DOM — retry after DOM stabilizes
        requestAnimationFrame(() => {
          try { Tooltip.mounted(el, binding) } catch (e) { /* ignore */ }
        })
        return
      }
      // Also try rendering in microtask to ensure DOM readiness
      requestAnimationFrame(() => {
        try { Tooltip.mounted(el, binding) } catch (e) { /* ignore */ }
      })
    } catch (e) {
      // Silently ignore PrimeVue internal errors during mount
    }
  },
  unmounted(el, binding) {
    try {
      Tooltip.unmounted(el, binding)
    } catch (e) {
      // Ignore DOM cleanup errors when element is already gone
    }
  },
  updated(el, binding) {
    // Re-apply tooltip when binding value changes (e.g., locale switch)
    try {
      // Destroy old tooltip first
      try { Tooltip.unmounted(el) } catch (e) { /* ignore */ }
      // Re-mount with new binding value
      requestAnimationFrame(() => {
        try { Tooltip.mounted(el, binding) } catch (e) { /* ignore */ }
      })
    } catch (e) { /* ignore */ }
  }
}

app.directive('tooltip', SafeTooltip)

app.mount('#app')
