import {createApp} from 'vue'
import PrimeVue from 'primevue/config'
import App from './App.vue'
import './style.scss';
import 'primevue/resources/themes/aura-light-green/theme.css'
import {createI18n} from 'vue-i18n';
import {en,de} from './locales'

const i18n = createI18n({
  legacy: false,
  locale: 'ja', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages: {en, de}
})

createApp(App).use(PrimeVue).use(i18n).mount('#app')
