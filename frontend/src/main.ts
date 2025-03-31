import { createApp } from "vue";
import PrimeVue, { defaultOptions } from "primevue/config";
import App from "./App.vue";
import "./style.scss";
import { createI18n } from "vue-i18n";
import { de, en } from "./locales";
import { createRouter, createWebHashHistory } from "vue-router";
import Tooltip from "primevue/tooltip";
import JournalPage from "./journal/JournalPage.vue";
import { createPinia } from "pinia";
import ConfirmationService from "primevue/confirmationservice";
import TrackPage from "./tracks/TrackPage.vue";
import FocusTrap from "primevue/focustrap";
import SettingsPage from "./settings/SettingsPage.vue";
import { useSettingsStore } from "./store/settings-store";
import AboutPage from "./about/AboutPage.vue";
import ToastService from "primevue/toastservice";
import Aura from "@primevue/themes/aura";
import DashboardPage from "./dashboard/DashboardPage.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", redirect: "dashboard" },
    { path: "/dashboard", component: DashboardPage },
    {
      path: "/journal/:entryId?",
      component: JournalPage,
    },
    { path: "/tracks/:trackId?", component: TrackPage },
    { path: "/settings", component: SettingsPage },
    {
      path: "/about",
      component: AboutPage,
    },
  ],
});

const i18n = createI18n({
  legacy: false,
  locale: "de",
  fallbackLocale: "en", // set fallback locale
  messages: { en, de },
  datetimeFormats: {
    en: {
      default: {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      },
      long: {
        weekday: "long",
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      },
    },
    de: {
      default: {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      },
      long: {
        weekday: "long",
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      },
    },
  },
});

createApp(App)
  .use(PrimeVue, {
    theme: { preset: Aura, prefix: "p" },
    locale: { ...defaultOptions.locale, firstDayOfWeek: 1 },
    ripple: true,
    zIndex: {
      overlay: 1000,
    },
  })
  .use(i18n)
  .use(router)
  .use(createPinia())
  .use(ConfirmationService)
  .use(ToastService)
  .directive("focustrap", FocusTrap)
  .directive("tooltip", Tooltip)
  .mount("#app");

useSettingsStore().loadSettings();
