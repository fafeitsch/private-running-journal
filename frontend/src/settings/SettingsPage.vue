<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import { defaultSettings, useSettingsStore } from "../store/settings-store";
import { settings as settingsType } from "../../wailsjs/go/models";
import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { storeToRefs } from "pinia";
import Button from "primevue/button";
import { useLeaveConfirmation } from "../shared/use-leave-confirmation";
import { useSettingsApi } from "../api/settings";
import GitIntegrationSettings from "./GitIntegrationSettings.vue";
import AppSettings = settingsType.AppSettings;

const settingsStore = useSettingsStore();
const { t } = useI18n();

const settings = ref<AppSettings>(new AppSettings(defaultSettings));

watch(
  storeToRefs(settingsStore).settings,
  () => {
    settings.value = new AppSettings(JSON.parse(JSON.stringify(settingsStore.settings)));
  },
  { immediate: true },
);

const dirty = ref(false);

useLeaveConfirmation(dirty);
const settingsApi = useSettingsApi();

async function saveSettings() {
  try {
    await settingsApi.saveSettings(settings.value);
    settingsStore.settings = settings.value;
    dirty.value = false;
  } catch (e) {
    console.error(e);
  }
}
</script>

<template>
  <div class="px-4 flex flex-column gap-2">
    <header class="flex gap-2 align-items-center">
      <h2 class="text-xl">{{ t("sidenav.settings") }}</h2>
      <Button icon="pi pi-save" :disabled="!dirty" @click="saveSettings" data-testid="save-settings-button"></Button>
    </header>
    <div class="flex flex-column gap-2 flex-grow-1 flex-shrink-1 overflow-auto">
      <Panel
        :header="t('settings.general.header')"
        :pt="{ content: { class: 'flex align-items-baseline gap-2' } }"
      >
        <InputGroup class="flex-grow-1 w-2">
          <InputGroupAddon>
            <label for="languageInput">{{ t("settings.general.language.label") }}</label>
          </InputGroupAddon>
          <Dropdown
            id="languageInput"
            v-model="settings.language"
            :options="['de', 'en']"
            :option-label="(option: string) => t('settings.general.language.' + option)"
            @update:model-value="dirty = true"
            data-testid="language-input"
          ></Dropdown>
        </InputGroup>
        <div class="flex flex-column gap-2 flex-grow-1">
          <InputGroup>
            <InputGroupAddon>
              <label for="portInput">{{ t("settings.general.port.label") }}</label>
            </InputGroupAddon>
            <InputNumber
              id="portInput"
              v-model="settings.httpPort"
              :use-grouping="false"
              @update:model-value="dirty = true"
              :max="65535"
            ></InputNumber>
          </InputGroup>
          <div class="text-xs">
            {{ t("settings.general.port.help") }}
          </div>
        </div>
      </Panel>
      <Panel
        :header="t('settings.mapSettings.header')"
        :pt="{ content: { class: 'flex-column flex gap-4' } }"
      >
        <div class="flex flex-column gap-2">
          <InputGroup>
            <InputGroupAddon>
              <label for="tileServerInput">{{ t("settings.mapSettings.tileServer.label") }}</label>
            </InputGroupAddon>
            <InputText
              id="tileServerInput"
              v-model="settings.mapSettings.tileServer"
              @update:model-value="dirty = true"
            ></InputText>
          </InputGroup>
          <i18n-t
            class="text-xs"
            keypath="settings.mapSettings.tileServer.help.template"
            tag="span"
            for="settings.mapSettings.tileServer.help.link"
          >
            <a
              href="https://wiki.openstreetmap.org/wiki/Raster_tile_providers"
              target="_blank"
              class="inline-flex align-items-baseline"
            >
              <span class="pi pi-external-link"></span>&nbsp;
              {{ $t("settings.mapSettings.tileServer.help.link") }}
            </a>
          </i18n-t>
        </div>
        <div class="flex flex-column gap-2">
          <InputGroup>
            <InputGroupAddon>
              <label for="tileAttributionInput">{{
                t("settings.mapSettings.attribution.label")
              }}</label>
            </InputGroupAddon>
            <InputText
              id="tileAttributionInput"
              v-model="settings.mapSettings.attribution"
              @update:model-value="dirty = true"
            ></InputText>
          </InputGroup>
          <div class="text-xs">{{ t("settings.mapSettings.attribution.help") }}</div>
        </div>
        <div class="flex flex-column gap-2">
          <InputGroup>
            <InputGroupAddon>
              <Checkbox
                id="cacheTilesInput"
                v-model="settings.mapSettings.cacheTiles"
                :disabled="false"
                :binary="true"
                @update:model-value="dirty = true"
              ></Checkbox>
            </InputGroupAddon>
            <InputGroupAddon class="flex-grow-1 justify-content-start">
              <label for="cacheTilesInput">{{ t("settings.mapSettings.cache.label") }}</label>
            </InputGroupAddon>
          </InputGroup>
          <div class="text-xs">{{ t("settings.mapSettings.cache.help") }}</div>
        </div>
        <div class="flex flex-column gap-2">
          <h4 class="m-0 p-0">{{ t("settings.mapSettings.mapPosition.general") }}</h4>
          <div class="flex gap-2">
            <InputGroup>
              <InputGroupAddon>
                <label for="latInput">{{ t("settings.mapSettings.mapPosition.latitude") }}</label>
              </InputGroupAddon>
              <InputNumber
                v-model="settings.mapSettings.center[0]"
                id="latInput"
                @update:model-value="dirty = true"
              ></InputNumber>
            </InputGroup>
            <InputGroup>
              <InputGroupAddon>
                <label for="lonInput">{{ t("settings.mapSettings.mapPosition.longitude") }}</label>
              </InputGroupAddon>
              <InputNumber
                v-model="settings.mapSettings.center[1]"
                id="lonInput"
                @update:model-value="dirty = true"
              ></InputNumber>
            </InputGroup>
            <InputGroup>
              <InputGroupAddon>
                <label for="zoomInput">{{ t("settings.mapSettings.mapPosition.zoom") }}</label>
              </InputGroupAddon>
              <InputNumber
                v-model="settings.mapSettings.zoomLevel"
                id="zoomInput"
                :min="1"
                :max="19"
                @update:model-value="dirty = true"
              ></InputNumber>
            </InputGroup>
          </div>
        </div>
      </Panel>
      <GitIntegrationSettings v-model="settings.gitSettings" @update:model-value="dirty = true" />
    </div>
  </div>
</template>

<style scoped>
a {
  text-decoration: underline;
}
</style>
