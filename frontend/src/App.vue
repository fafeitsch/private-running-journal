<script lang="ts" setup>
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { TabMenuChangeEvent } from "primevue/tabmenu";
import { useRoute, useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useJournalStore } from "./store/journal-store";
import { useTrackStore } from "./store/track-store";
import ConfirmDialog from "primevue/confirmdialog";
import { useSettingsStore } from "./store/settings-store";

const { t, locale } = useI18n();

const journalStoreRef = storeToRefs(useJournalStore());
const trackStoreRef = storeToRefs(useTrackStore());
const { settings } = storeToRefs(useSettingsStore());

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
  },
  {
    label: t("sidenav.tracks"),
    icon: "pi pi-directions",
    link: `/tracks/${encodeURIComponent(trackStoreRef.selectedTrackId.value || "")}`,
  },
  {
    label: t("sidenav.settings"),
    icon: "pi pi-cog",
    link: `/settings/`,
  },
  {
    label: t("sidenav.about"),
    icon: "pi pi-info-circle",
    link: `/about/`,
  },
]);
const active = ref(0);

const router = useRouter();
function tabChangeEvent(event: TabMenuChangeEvent) {
  router.push(navItems.value[event.index].link);
}

const route = useRoute();
watch(
  () => route.fullPath,
  (value) => {
    const item = navItems.value.findIndex((i) => i.link.startsWith(value));
    if (item > -1) {
      active.value = item;
    }
  },
);
</script>

<template>
  <div class="h-full flex flex-column">
    <TabMenu
      class="flex-shrink-0"
      :model="navItems"
      v-model:active-index="active"
      @tab-change="tabChangeEvent"
    ></TabMenu>
    <div class="flex-grow-1 flex-shrink-1 w-full overflow-hidden">
      <router-view class="h-full w-full"></router-view>
    </div>
  </div>
  <ConfirmDialog group="leave"></ConfirmDialog>
</template>

<style></style>
