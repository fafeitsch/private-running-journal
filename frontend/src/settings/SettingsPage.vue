<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import { defaultSettings, useSettingsStore } from "../store/settings-store";
import { settings as settingsType } from "../../wailsjs/go/models";
import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { storeToRefs } from "pinia";
import Button from "primevue/button";
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

function saveSettings() {}
</script>

<template>
  <div class="px-4 flex flex-column gap-2">
    <header class="flex gap-2 align-items-center">
      <h2 class="text-xl">{{ t("sidenav.settings") }}</h2>
      <Button icon="pi pi-save" :disabled="!dirty" @click="saveSettings"></Button>
    </header>
    <Panel :header="t('settings.general.header')">
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
        <InputGroupAddon>
          {{ t("settings.general.port.help") }}
        </InputGroupAddon>
      </InputGroup>
    </Panel>
    <Panel
      :header="t('settings.mapSettings.header')"
      :pt="{ content: { class: 'flex-column flex gap-2' } }"
    >
      <InputGroup>
        <InputGroupAddon>
          <label for="tileServerInput">{{ t("settings.mapSettings.tileServer.label") }}</label>
        </InputGroupAddon>
        <InputText
          id="tileServerInput"
          v-model="settings.mapSettings.tileServer"
          @update:model-value="dirty = true"
        ></InputText>
        <InputGroupAddon>
          <i18n-t
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
        </InputGroupAddon>
      </InputGroup>
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
        <InputGroupAddon>
          {{ t("settings.mapSettings.cache.help") }}
        </InputGroupAddon>
      </InputGroup>
      <InputGroup>
        <InputGroupAddon>
          {{ t("settings.mapSettings.mapPosition.general") }}
        </InputGroupAddon>
        <InputGroupAddon>
          <label for="latInput">{{ t("settings.mapSettings.mapPosition.latitude") }}</label>
        </InputGroupAddon>
        <InputText
          v-model="settings.mapSettings.center[0]"
          id="latInput"
          @update:model-value="dirty = true"
        ></InputText>
        <InputGroupAddon>
          <label for="lonInput">{{ t("settings.mapSettings.mapPosition.longitude") }}</label>
        </InputGroupAddon>
        <InputText
          v-model="settings.mapSettings.center[1]"
          id="lonInput"
          @update:model-value="dirty = true"
        ></InputText>
        <InputGroupAddon>
          <label for="zoomInput">{{ t("settings.mapSettings.mapPosition.zoom") }}</label>
        </InputGroupAddon>
        <InputText
          v-model="settings.mapSettings.zoomLevel"
          id="zoomInput"
          :min="1"
          :max="19"
          @update:model-value="dirty = true"
        ></InputText>
      </InputGroup>
    </Panel>
  </div>
</template>

<style scoped>
a {
  text-decoration: underline;
}
</style>
