<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import ProgressSpinner from "primevue/progressspinner";
import Message from "primevue/message";
import Button from "primevue/button";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";
import CreateEntryOverlay from "./CreateEntryOverlay.vue";
import { useJournalStore } from "../store/journal-store";
import { storeToRefs } from "pinia";
import { MenuItem } from "primevue/menuitem";
import { useJournalApi } from "../api/journal";

const { t, d } = useI18n();
const loading = ref(false);
const error = ref<boolean>(false);
const store = useJournalStore();
const journalApi = useJournalApi();
const { listEntries } = storeToRefs(store);

const route = useRoute();

onMounted(async () => {
  await load();
});

async function load() {
  loading.value = true;
  error.value = false;
  try {
    await store.loadEntries();
  } catch (e) {
    error.value = true;
    console.error(e);
  } finally {
    loading.value = false;
  }
}

const entries = computed(() => {
  const result = listEntries.value.map((entry) => ({
    ...entry,
    link: encodeURIComponent(entry.id),
    trackError: !entry.trackName && !entry.length,
    length: (entry.length / 1000).toFixed(1),
  }));
  result.sort((a, b) => a.date.localeCompare(b.date));
  return result;
});
</script>

<template>
  <div class="h-full overflow-hidden flex flex-column">
    <header class="flex justify-content-between align-items-center">
      <span class="text-2xl">{{ $t("journal.entries") }}</span
      ><CreateEntryOverlay></CreateEntryOverlay>
    </header>
    <div v-if="loading" class="flex-grow-1 flex-shrink-1 flex flex-column justify-content-center">
      <ProgressSpinner></ProgressSpinner>
    </div>
    <Message v-else-if="error" severity="error" :closable="false"
      ><div class="flex align-items-center">
        <span>{{ t("journal.loadingError") }}</span>
        <Button severity="danger" rounded text icon="pi pi-replay" @click="load"></Button></div
    ></Message>
    <ul v-else class="list-none p-0 mt-3 flex-grow-1 flex-shrink-1 overflow-auto">
      <li
        v-for="entry of entries"
        :key="entry.id"
        class="list-none p-0 m-0 flex"
        v-tooltip="{ value: entry.trackName, showDelay: 500 }"
      >
        <RouterLink
          v-ripple
          class="flex-grow-1 flex-shrink-1 flex gap-1 cursor-pointer p-ripple transition-colors hover:surface-100 transition-duration-150 text-700 py-3 px-1"
          :to="{ path: '/journal/' + entry.link }"
          active-class="surface-200"
          ><span class="font-medium">{{ d(entry.date, 'default') }}</span>
          <template v-if="!entry.trackError">
            <span
              class="font-normal flex-shrink-1 flex-grow-1 text-overflow-ellipsis overflow-hidden white-space-nowrap"
              >{{ entry.trackName }}</span
            >
            <span class="font-medium">{{ entry.length }}&nbsp;km</span>
          </template>
          <div v-else class="flex overflow-hidden ml-2 flex gap-2 text-red-500">
            <span class="pi pi-exclamation-triangle"></span><span class="white-space-nowrap overflow-hidden text-overflow-ellipsis">{{ t("journal.noTrack") }}</span>
          </div>
        </RouterLink>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.variant {
  flex-basis: 0;
}
</style>
