<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel";
import InlineMessage from "primevue/inlinemessage";
import Button from "primevue/button";
import { ref } from "vue";
import Calendar from "primevue/calendar";
import { useI18n } from "vue-i18n";
import { useJournalApi } from "../api/journal";
import TrackSelection from "./TrackSelection.vue";
import { tracks } from "../../wailsjs/go/models";
import { useJournalStore } from "../store/journal-store";
import { useRouter } from "vue-router";

const { locale, t } = useI18n();

const overlayPanel = ref();
const selectedDate = ref<Date>(new Date());
const error = ref<boolean>(false);
const selectedTrack = ref<tracks.Track | undefined>(undefined);

const journalApi = useJournalApi();
const store = useJournalStore();
const router = useRouter();

async function createEntry() {
  if (!selectedTrack.value) {
    return;
  }
  error.value = false;
  const date = selectedDate.value;
  const dateString = `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")}`;
  try {
    const entry = await journalApi.createJournalEntry(dateString, selectedTrack.value.id);
    store.addEntryToList(entry);
    router.push("/journal/" + encodeURIComponent(entry.id));
    overlayPanel.value.hide();
  } catch (e) {
    error.value = true;
    console.error(e);
  }
}
</script>

<template>
  <Button
    icon="pi pi-plus"
    @click="(event) => overlayPanel.toggle(event)"
    :aria-label="t('shared.add')"
    :v-tooltip="t('shared.add')"
  ></Button>
  <OverlayPanel ref="overlayPanel">
    <div class="flex flex-column gap-2 overlay">
      <TrackSelection v-model="selectedTrack"></TrackSelection>
      <Calendar
        id="date"
        inline
        v-model="selectedDate"
        show-button-bar
        :date-format="locale === 'de' ? 'dd.mm.yy' : 'yyyy/mm/dd'"
        :pt="{
          todayButton: {
            root: { 'data-testid': 'create-entry-today-button' },
          },
        }"
      ></Calendar>
      <div class="flex gap-2">
        <InlineMessage v-if="error" class="flex-grow-1 flex-shrink-1" severity="error">{{
          t("journal.createEntryError")
        }}</InlineMessage>
        <span v-else class="flex-grow-1"></span>
        <Button
          :label="t('journal.createEntry')"
          @click="createEntry"
          :disabled="!selectedTrack"
          data-testid="create-entry-button"
        ></Button>
      </div>
    </div>
  </OverlayPanel>
</template>

<style scoped>
.overlay {
  width: 400px;
}
</style>
