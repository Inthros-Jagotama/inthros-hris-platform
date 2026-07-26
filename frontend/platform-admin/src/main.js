import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Tooltip from 'primevue/tooltip'
import router from './router'
import App from './App.vue'
import './assets/styles/main.css'

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
// to catch PrimeVue v4 DOM errors (parentNode null, isUnstyled)
// when tooltip elements are manipulated during dynamic rendering.
//
// Ref: https://github.com/primefaces/primevue/issues/6478
const SafeTooltip = {
  mounted(el, binding) {
    try {
      Tooltip.mounted(el, binding)
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
  }
}

app.directive('tooltip', SafeTooltip)

app.mount('#app')
