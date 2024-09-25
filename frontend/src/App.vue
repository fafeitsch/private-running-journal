<script lang="ts" setup>
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useJournalStore } from "./store/journal-store";
import { useTrackStore } from "./store/track-store";
import ConfirmDialog from "primevue/confirmdialog";
import { useSettingsStore } from "./store/settings-store";
import { EventsOn } from "../wailsjs/runtime";
import { useToast } from "primevue/usetoast";
import Toast from "primevue/toast";
import Tabs from "primevue/tabs";
import Tab from "primevue/tab";
import TabList from "primevue/tablist";

const { t, locale } = useI18n();

const journalStoreRef = storeToRefs(useJournalStore());
const trackStoreRef = storeToRefs(useTrackStore());
const { settings, loaded: settingsLoaded } = storeToRefs(useSettingsStore());

watch(
  () => settings.value.language,
  (language) => {
    locale.value = language;
  },
);

const navItems = computed(() => [
  {
    label: t("sidenav.journal"),
    icon: "pi pi-list",
    link: `/journal/${encodeURIComponent(journalStoreRef.selectedEntryId.value || "")}`,
    testId: "journal-tab",
  },
  {
    label: t("sidenav.tracks"),
    icon: "pi pi-directions",
    link: `/tracks/${encodeURIComponent(trackStoreRef.selectedTrackId.value || "")}`,
    testId: "tracks-tab",
  },
  {
    label: t("sidenav.settings"),
    icon: "pi pi-cog",
    link: `/settings/`,
    testId: "settings-tab",
  },
  {
    label: t("sidenav.about"),
    icon: "pi pi-info-circle",
    link: `/about/`,
    testId: "about-tab",
  },
]);
const active = ref(0);

const route = useRoute();
const router = useRouter();
watch(
  () => route.fullPath,
  (value) => {
    const item = navItems.value.findIndex((i) => i.link.startsWith(value));
    if (item > -1) {
      active.value = item;
    }
  },
);

const toast = useToast();

onMounted(() => {
  EventsOn("git-error", (stack: string) => {
    toast.add({ closable: true, detail: stack, severity: "error", summary: t("shared.gitError") });
  });
});
</script>

<template>
  <div v-if="settingsLoaded" class="h-full flex flex-col">
    <Tabs v-model:value="active" class="shrink-0">
      <TabList>
        <Tab
          v-for="(tab, index) in navItems"
          :key="tab.link"
          :value="index"
          :data-testid="'nav-' + tab.testId"
          class="!font-normal"
          @click="router.push(tab.link)"
        >
          <router-link :to="tab.link" v-slot="{ href, navigate }" custom
            ><a :href="href" @click="navigate" class="flex gap-2 items-center"
              ><i :class="tab.icon"></i>{{ tab.label }}</a
            >
          </router-link>
        </Tab>
      </TabList>
    </Tabs>
    <div class="grow shrink w-full overflow-hidden">
      <router-view class="h-full w-full"></router-view>
    </div>
  </div>
  <ConfirmDialog group="leave"></ConfirmDialog>
  <Toast></Toast>
</template>

<style></style>
