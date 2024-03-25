import { settings } from "../../wailsjs/go/models";
import {GetSettings, SaveSettings} from "../../wailsjs/go/backend/App";
import AppSettings = settings.AppSettings;

export const useSettingsApi = () => {
  function getSettings(): Promise<AppSettings> {
    return GetSettings();
  }
  function saveSettings(settings: AppSettings): Promise<void> {
    return SaveSettings(settings);
  }
  return { getSettings, saveSettings };
};
