<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { ref, watch } from "vue";
import { useJournalApi } from "../api/journal";
import ProgressSpinner from "primevue/progressspinner";
import Message from "primevue/message";
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputGroup from "primevue/inputgroup";
import { journal, trackEditor } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import LeafletMap from "./LeafletMap.vue";
import TrackTimeResult from "./TrackTimeResult.vue";
import TrackSelection from "./TrackSelection.vue";
import { useJournalStore } from "../store/journal-store";
import { storeToRefs } from "pinia";
import { useLeaveConfirmation } from "../shared/use-leave-confirmation";
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from "primevue/confirmpopup";
import DatePicker from "primevue/datepicker";
import Entry = journal.Entry;
import TrackDto = trackEditor.TrackDto;

const { t, d, locale } = useI18n();
const route = useRoute();
const journalApi = useJournalApi();
const router = useRouter();

const loading = ref(false);
const loadError = ref(false);
const editError = ref<string | undefined>(undefined);
const dirty = ref(false);
const customLengthEnabled = ref(false);
const journalEntryLength = ref<number | undefined>(undefined);
const selectedEntry = ref<journal.Entry | undefined>(undefined);
const selectedDate = ref<Date>(new Date());
const journalStore = useJournalStore();
const { selectedEntryId } = storeToRefs(journalStore);

const tracksApi = useTracksApi();

watch(
  () => {
    return route.params.entryId as string;
  },
  (value: string) => {
    dirty.value = false;
    return loadEntry(value);
  },
  { immediate: true },
);

async function loadEntry(entryId: string | undefined) {
  loadError.value = false;
  if (entryId === "new") {
    selectedEntry.value = new Entry({
      id: entryId,
      date: new Date().toISOString(),
      comment: "",
      laps: 1,
      time: "",
      track: undefined,
    });
    selectedDate.value = new Date();
    return;
  }
  selectedEntry.value = undefined;
  loading.value = true;
  selectedEntryId.value = entryId;
  if (!entryId) {
    await router.replace("/journal");
    selectedEntry.value = undefined;
    selectedDate.value = new Date();
    return;
  }
  try {
    selectedEntry.value = await journalApi.getListEntry(entryId);
    customLengthEnabled.value = !!selectedEntry.value.customLength;
    calculateLength();
    selectedDate.value = new Date(Date.parse(selectedEntry.value.date));
  } catch (e) {
    console.error(e);
    loadError.value = true;
  } finally {
    loading.value = false;
  }
}

const gpxData = ref<TrackDto | undefined>(undefined);

watch(
  selectedEntry,
  async () => {
    if (!selectedEntry.value || !selectedEntry.value.track) {
      gpxData.value = undefined;
      return;
    }
    try {
      gpxData.value = await tracksApi.getTrack(selectedEntry.value.track.id);
    } catch (e) {
      console.error(e);
      loadError.value = true;
    }
  },
  { deep: true },
);

async function saveEntry() {
  let value = selectedEntry.value;
  if (!value) {
    return;
  }
  if (customLengthEnabled.value) {
    value.customLength = journalEntryLength.value! * 1000;
  } else {
    value.customLength = undefined;
  }
  editError.value = undefined;
  try {
    const length = customLengthEnabled.value
      ? value.customLength!
      : value.track!.length * value.laps;
    if (value.id === "new") {
      const year = `${selectedDate.value.getFullYear()}`.padStart(4, "0");
      const month = `${selectedDate.value.getMonth() + 1}`.padStart(2, "0");
      const day = `${selectedDate.value.getDate()}`.padStart(2, "0");
      value.date = `${year}-${month}-${day}`;
      const result = await journalApi.createJournalEntry(value);
      dirty.value = false;
      journalStore.addEntryToList({
        date: value.date,
        trackName: result.trackName,
        length: result.length,
        trackError: false,
        id: result.id,
      });
      router.replace("/journal/" + encodeURIComponent(result.id));
      return;
    }
    await journalApi.saveEntry(value);
    journalStore.updateEntry({
      ...value,
      trackName: value.track!.name,
      length,
    });
    dirty.value = false;
  } catch (e) {
    console.error(e);
    editError.value = "journal.saveError";
  }
}

const confirm = useConfirm();

async function deleteEntry(event: Event) {
  if (!selectedEntry.value) {
    return;
  }
  let resolveFn: (result: boolean) => void;
  const result = new Promise<boolean>((resolve) => (resolveFn = resolve));
  confirm.require({
    target: event.currentTarget as HTMLElement,
    group: "journal",
    header: t("shared.confirm.header"),
    accept: () => resolveFn(true),
    reject: () => resolveFn(false),
    message: t("journal.deleteConfirmation"),
    rejectLabel: t("shared.cancel"),
    acceptLabel: t("shared.delete"),
  });
  let choice = await result;
  if (!choice) {
    return;
  }
  editError.value = undefined;
  try {
    await journalApi.deleteEntry(selectedEntry.value.id);
    journalStore.deleteEntry(selectedEntry.value.id);
  } catch (e) {
    console.error(e);
    editError.value = "journal.deleteError";
  }
}

function onChangeCustomLengthEnabled() {
  if (!customLengthEnabled.value && selectedEntry.value?.track) {
    journalEntryLength.value = selectedEntry.value.track.length / 1000;
  }
  dirty.value = true;
}

function onTrackSelectionChanged() {
  calculateLength();
  dirty.value = true;
}

function calculateLength() {
  if (!selectedEntry.value) {
    return;
  }
  let length = undefined;
  if (customLengthEnabled.value) {
    length = selectedEntry.value.customLength;
  } else {
    length = selectedEntry.value.track?.length || undefined;
  }
  journalEntryLength.value = (length || 0) / 1000;
}

useLeaveConfirmation(dirty);
</script>

<template>
  <div class="flex flex-col">
    <div v-if="loading" class="flex w-full grow justify-center items-center">
      <ProgressSpinner></ProgressSpinner>
    </div>
    <div v-else-if="loadError" class="px-2" data-testid="journal-entry-load-error">
      <Message severity="error" :closable="false"
        ><div class="flex items-center">
          <span>{{ t("journal.loadEntryError") }}</span>
          <Button
            severity="danger"
            rounded
            text
            icon="pi pi-replay"
            @click="loadEntry(route.params.entryId as string)"
          ></Button></div
      ></Message>
    </div>
    <div v-else-if="selectedEntry" class="flex flex-col gap-2 w-full p-2 grow shrink">
      <div class="flex gap-2">
        <Button
          icon="pi pi-save"
          :aria-label="t('shared.save')"
          v-tooltip="{ value: t('shared.save'), showDelay: 500 }"
          :disabled="!dirty || !selectedEntry.track"
          @click="saveEntry"
        ></Button>
        <Button
          icon="pi pi-trash"
          @click="deleteEntry($event)"
          :aria-label="t('shared.delete')"
          v-tooltip="{ value: t('shared.delete'), showDelay: 500 }"
        ></Button>
        <ConfirmPopup group="journal">
          <template #message="{ message }">
            <div
              style="max-width: 330px"
              class="p-2"
              data-testid="delete-journal-entry-confirmation"
            >
              {{ message.message }}
            </div>
          </template>
        </ConfirmPopup>
        <Message class="m-0" v-if="editError" :severity="'error'" :closable="false">
          {{ t(editError) }}
        </Message>
      </div>
      <InputGroup>
        <InputGroupAddon>
          <label for="date">{{ t("journal.details.date") }}</label>
        </InputGroupAddon>
        <DatePicker
          id="date"
          v-model="selectedDate"
          show-button-bar
          :date-format="locale === 'de' ? 'dd.mm.yy' : 'yyyy/mm/dd'"
          data-testid="entry-date-input"
          :pt="{
            pcTodayButton: {
              root: { 'data-testid': 'journal-entry-today-button' },
            },
          }"
          :disabled="selectedEntry?.id !== 'new'"
        ></DatePicker>
      </InputGroup>
      <TrackSelection
        v-model="selectedEntry!.track"
        :linked-track="selectedEntry!.linkedTrack"
        @update:model-value="() => onTrackSelectionChanged()"
      ></TrackSelection>
      <TrackTimeResult
        v-model:laps="selectedEntry!.laps"
        v-model:time="selectedEntry!.time"
        v-model:track-length="journalEntryLength"
        v-model:custom-length="customLengthEnabled"
        @update:track-length="() => (dirty = true)"
        @update:laps="() => (dirty = true)"
        @update:time="() => (dirty = true)"
        @update:custom-length="onChangeCustomLengthEnabled"
      ></TrackTimeResult>
      <InputGroup>
        <InputGroupAddon>
          <label for="comment">{{ t("journal.details.comment") }}</label>
        </InputGroupAddon>
        <InputText
          id="comment"
          v-model="selectedEntry!.comment"
          @update:model-value="() => (dirty = true)"
        ></InputText>
      </InputGroup>
      <div v-if="gpxData" class="shrink grow">
        <LeafletMap
          class="h-full w-full"
          :waypoints="gpxData?.waypoints || []"
          :polyline-meta="gpxData"
        ></LeafletMap>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
