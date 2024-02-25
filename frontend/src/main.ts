import { createApp } from "vue";
import PrimeVue from "primevue/config";
import App from "./App.vue";
import "./style.scss";
import "primevue/resources/themes/aura-light-green/theme.css";
import { createI18n } from "vue-i18n";
import { en, de } from "./locales";
import { createRouter, createWebHashHistory } from "vue-router";
import Ripple from "primevue/ripple";
import JournalPage from "./journal/JournalPage.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", redirect: 'journal'},
    {
      path: "/journal/:entryId?",
      component: JournalPage,
    },
  ],
});

const i18n = createI18n({
  legacy: false,
  locale: "de",
  fallbackLocale: "en", // set fallback locale
  messages: { en, de },
});


createApp(App).use(PrimeVue).use(i18n).use(router).directive("ripple", Ripple).mount("#app");
