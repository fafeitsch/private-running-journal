import { settings } from "../../wailsjs/go/models";
import { GetSettings } from "../../wailsjs/go/backend/App";
import AppSettings = settings.AppSettings;

export const useSettingsApi = () => {
  function getSettings(): Promise<AppSettings> {
    return GetSettings();
  }
  return { getSettings };
};
