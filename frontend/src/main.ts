import { createApp } from "vue";
import PrimeVue from "primevue/config";
import App from "./App.vue";
import "./style.scss";
import "primevue/resources/themes/aura-light-green/theme.css";
import { createI18n } from "vue-i18n";
import { de, en } from "./locales";
import { createRouter, createWebHashHistory } from "vue-router";
import Ripple from "primevue/ripple";
import Tooltip from "primevue/tooltip";
import JournalPage from "./journal/JournalPage.vue";
import { createPinia } from "pinia";
import ConfirmationService from "primevue/confirmationservice";
import TrackPage from "./tracks/TrackPage.vue";
import FocusTrap from 'primevue/focustrap';

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", redirect: "journal" },
    {
      path: "/journal/:entryId?",
      component: JournalPage,
    },
    { path: "/tracks/:trackId?", component: TrackPage },
  ],
});

const i18n = createI18n({
  legacy: false,
  locale: "de",
  fallbackLocale: "en", // set fallback locale
  messages: { en, de },
});

createApp(App)
  .use(PrimeVue)
  .use(i18n)
  .use(router)
  .use(createPinia())
  .use(ConfirmationService)
  .directive('focustrap', FocusTrap)
  .directive("tooltip", Tooltip)
  .directive("ripple", Ripple)
  .mount("#app");
