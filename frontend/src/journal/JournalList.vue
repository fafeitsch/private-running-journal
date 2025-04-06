<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import ProgressSpinner from "primevue/progressspinner";
import Message from "primevue/message";
import Button from "primevue/button";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { useJournalStore } from "../store/journal-store";
import { storeToRefs } from "pinia";
import MonthChooser from '../shared/MonthChooser.vue';

const { t, d, n, locale } = useI18n();
const loading = ref(false);
const error = ref<boolean>(false);
const store = useJournalStore();
const { listEntries, selectedMonth } = storeToRefs(store);

const router = useRouter();

const route = useRoute();

onMounted(async () => {
  await load();
});

watch(selectedMonth, () => load());

async function load() {
  loading.value = true;
  error.value = false;
  try {
    let month = selectedMonth.value;
    if (!month) {
      return;
    }
    const start = month.toISOString();
    const clone = new Date(month.getTime());
    const end = new Date(clone.setMonth(month.getMonth() + 1)).toISOString();
    await store.loadEntries(start, end);
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
    length: n(entry.length / 1000, { maximumFractionDigits: 1, minimumFractionDigits: 1 }),
  }));
  result.sort((a, b) => a.date.localeCompare(b.date));
  return result;
});
</script>

<template>
  <div class="h-full overflow-hidden flex flex-col gap-2">
    <header class="flex justify-between items-center">
      <span class="text-2xl">{{ t("journal.entries") }}</span>
      <Button
        icon="pi pi-plus"
        @click="router.push('/journal/new')"
        :aria-label="t('shared.add')"
        v-tooltip="{ value: t('shared.add'), showDelay: 500 }"
      ></Button>
    </header>
    <MonthChooser v-model="selectedMonth"></MonthChooser>
    <div v-if="loading" class="grow shrink flex flex-col justify-center">
      <ProgressSpinner></ProgressSpinner>
    </div>
    <Message v-else-if="error" severity="error" :closable="false">
      <div class="flex items-center">
        <span>{{ t("journal.loadingError") }}</span>
        <Button severity="danger" rounded text icon="pi pi-replay" @click="load"></Button>
      </div>
    </Message>
    <ul v-else class="list-none p-0 grow shrink overflow-y-auto overflow-x-hidden">
      <li
        v-for="entry of entries"
        :key="entry.id"
        class="list-none p-0 m-0 flex overflow-hidden"
        v-tooltip="{ value: entry.trackName, showDelay: 500 }"
        data-testid="journal-entry-item"
      >
        <RouterLink
          class="grow shrink flex items-center gap-1 cursor-pointer overflow-hidden p-ripple transition-colors hover:bg-surface-100 dark:hover:bg-surface-700 duration-150 text-surface-700 dark:text-surface-100 py-4 px-1"
          :to="{ path: '/journal/' + entry.link }"
          active-class="bg-surface-200 dark:bg-surface-600"
          ><span class="font-medium">{{ d(entry.date, "default") }}</span>
          <template v-if="!entry.trackError">
            <span class="font-normal shrink grow text-ellipsis overflow-hidden whitespace-nowrap">{{
              entry.trackName
            }}</span>
            <span class="font-medium">{{ entry.length }}&nbsp;km</span>
          </template>
          <div v-else class="flex overflow-hidden ml-2 items-center gap-2 text-red-500">
            <span class="pi pi-exclamation-triangle"></span
            ><span class="whitespace-nowrap overflow-hidden text-ellipsis">{{
              t("journal.noTrack")
            }}</span>
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
