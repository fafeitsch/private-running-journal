<script setup lang="ts">
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import JournalList from "./JournalList.vue";
import { useRoute } from "vue-router";
import { computed } from "vue";
import JournalEntryDetail from "./JournalEntryDetail.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const route = useRoute();

const selectedEntry = computed(() => route.params["entryId"]);
</script>

<template>
  <Splitter>
    <SplitterPanel class="flex flex-col p-2" :size="20">
      <div class="h-full overflow-hidden">
        <JournalList></JournalList>
      </div>
    </SplitterPanel>
    <SplitterPanel class="flex items-center justify-center" :size="80">
      <JournalEntryDetail v-if="selectedEntry" class="h-full w-full"> </JournalEntryDetail>
      <div v-else class="p-2" data-testid="empty-journal-state">{{ t("journal.emptyState") }}</div>
    </SplitterPanel>
  </Splitter>
</template>

<style scoped></style>
