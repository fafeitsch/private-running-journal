import { defineStore } from "pinia";
import { settings } from "../../wailsjs/go/models";
import { ref } from "vue";
import { useSettingsApi } from "../api/settings";
import AppSettings = settings.AppSettings;

export const defaultSettings = {
  mapSettings: {
    attribution: "",
    cacheTiles: false,
    tileServer: "",
    center: [0, 0],
    zoomLevel: 0,
  },
  httpPort: 47836,
};

const settingsApi = useSettingsApi();
export const useSettingsStore = defineStore("settings", () => {
  const settings = ref<AppSettings>(new AppSettings(defaultSettings));

  async function loadSettings() {
    settings.value = await settingsApi.getSettings();
    console.log(settings.value.mapSettings.tileServer);
  }

  return { settings, loadSettings };
});
