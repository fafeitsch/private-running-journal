import { defineStore } from "pinia";
import { settings } from "../../wailsjs/go/models";
import {ref, watch} from "vue";
import { useSettingsApi } from "../api/settings";
import AppSettings = settings.AppSettings;
import {useI18n} from 'vue-i18n';

export const defaultSettings = {
  mapSettings: {
    attribution: "",
    cacheTiles: false,
    tileServer: "",
    center: [0, 0],
    zoomLevel: 0,
  },
  gitSettings: {
    enabled: false,
    pullOnStartUp: false,
    pushAfterCommit: false,
  },
  language: "en",
  httpPort: 47836,
};

const settingsApi = useSettingsApi();
export const useSettingsStore = defineStore("settings", () => {
  const settings = ref<AppSettings>(new AppSettings(defaultSettings));

  const i18n = useI18n()

  async function loadSettings() {
    settings.value = await settingsApi.getSettings();
    i18n.locale.value = settings.value.language
  }

  watch(() => settings.value.language, (language) => {
    i18n.locale.value = language
  })

  return { settings, loadSettings };
});
