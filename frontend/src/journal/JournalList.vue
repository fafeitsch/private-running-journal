<script setup lang="ts">
import { useJournalApi } from "../api/journal";
import { computed, onMounted, ref } from "vue";
import { backend } from "../../wailsjs/go/models";
import ProgressSpinner from "primevue/progressspinner";
import Message from "primevue/message";
import Button from "primevue/button";
import { useI18n } from "vue-i18n";
import JournalListEntry = backend.JournalListEntry;
import {useRoute} from 'vue-router';

const { t, d } = useI18n();
const loading = ref(false);
const entries = ref<JournalListEntry[]>([]);
const error = ref<boolean>(false);
const { getListEntries } = useJournalApi();

const route = useRoute()
const selectedEntry = computed(() => route.params['routeId'])

onMounted(async () => {
  await loadEntries();
});

async function loadEntries() {
  loading.value = true;
  error.value = false;
  entries.value = [];
  try {
    entries.value = await getListEntries();
  } catch (e) {
    error.value = true;
    console.error(e);
  } finally {
    loading.value = false;
  }
}

const formattedEntries = computed(() =>
  entries.value.map((entry) => ({ ...entry, length: (entry.length / 1000).toFixed(1) })),
);
</script>

<template>
  <header class="text-2xl">{{ $t("journal.entries") }}</header>
  <div v-if="loading" class="h-full flex flex-column justify-content-center">
    <ProgressSpinner></ProgressSpinner>
  </div>
  <Message v-else-if="error" severity="error" :closable="false"
    ><div class="flex align-items-center">
      <span>{{ t("journal.loadingError") }}</span>
      <Button severity="danger" rounded text icon="pi pi-replay" @click="loadEntries"></Button></div
  ></Message>
  <ul v-else class="list-none p-0 mt-3 h-full overflow-auto">
    <li
      v-for="entry of formattedEntries"
      :key="entry.id"
      class="list-none p-0 m-0 flex flex-column"
    >
      <RouterLink
        v-ripple
        class="flex gap-1 cursor-pointer p-ripple transition-colors hover:surface-100 transition-duration-150 text-700 p-3"
        :to="{path: '/journal/' + entry.id}"
        active-class="surface-200"
        ><span class="font-medium">{{ d(entry.date) }}</span>
        <span
          class="font-normal flex-shrink-1 overflow-hidden white-space-nowrap text-overflow-ellipsis"
          >{{ entry.trackBaseName }}</span
        ><span
          class="variant font-light flex-grow-1 flex-shrink-1 text-overflow-ellipsis overflow-hidden white-space-nowrap"
          >&nbsp;–&nbsp;{{ entry.trackVariant }}</span
        >
        <span class="font-medium">{{ entry.length }}&nbsp;km</span></RouterLink
      >
    </li>
  </ul>
</template>

<style scoped>
.variant {
  flex-basis: 0
}
</style>
