import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Tooltip from 'primevue/tooltip'
import router from './router'
import App from './App.vue'
import './assets/styles/main.css'
import 'primeicons/primeicons.css'

const TealPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#f4f9fa',
      100: '#e6f1f3',
      200: '#c8e0e5',
      300: '#a8ced6',
      400: '#84bac5',
      500: '#549fae',
      600: '#1b7f93',
      700: '#176c7d',
      800: '#135967',
      900: '#0f4651',
      950: '#0b333b'
    }
  }
})

const app = createApp(App)

app.use(PrimeVue, {
  theme: {
    preset: TealPreset,
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
app.directive('tooltip', Tooltip)

app.mount('#app')
