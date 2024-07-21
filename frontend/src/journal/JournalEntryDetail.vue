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
import { journal, tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import LeafletMap from "./LeafletMap.vue";
import TrackTimeResult from "./TrackTimeResult.vue";
import TrackSelection from "./TrackSelection.vue";
import { useJournalStore } from "../store/journal-store";
import { storeToRefs } from "pinia";
import { useLeaveConfirmation } from "../shared/use-leave-confirmation";
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from "primevue/confirmpopup";

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
  selectedEntry.value = undefined;
  loadError.value = false;
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
    let length = undefined;
    if (customLengthEnabled.value) {
      length = selectedEntry.value.customLength;
    } else {
      length = selectedEntry.value.track?.length || undefined;
    }
    journalEntryLength.value = (length || 0) / 1000;
    selectedDate.value = new Date(Date.parse(selectedEntry.value.date));
  } catch (e) {
    console.error(e);
    loadError.value = true;
  } finally {
    loading.value = false;
  }
}

const gpxData = ref<tracks.GpxData | undefined>(undefined);

watch(
  selectedEntry,
  async () => {
    if (!selectedEntry.value || !selectedEntry.value.track) {
      gpxData.value = undefined;
      return;
    }
    try {
      gpxData.value = await tracksApi.getGpxData(selectedEntry.value.track.id);
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
    await journalApi.saveEntry(value);
    const length = customLengthEnabled.value
      ? value.customLength!
      : value.track!.length * value.laps;
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

useLeaveConfirmation(dirty);
</script>

<template>
  <div class="flex flex-column">
    <div v-if="loading" class="flex w-full flex-grow-1 justify-content-center align-items-center">
      <ProgressSpinner></ProgressSpinner>
    </div>
    <div v-else-if="loadError" class="px-2">
      <Message severity="error" :closable="false"
        ><div class="flex align-items-center">
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
    <div
      v-else-if="selectedEntry"
      class="flex flex-column gap-2 w-full p-2 flex-grow-1 flex-shrink-1"
    >
      <div class="flex gap-2">
        <Button
          icon="pi pi-save"
          :disabled="!dirty"
          :aria-label="t('shared.save')"
          v-tooltip="{ value: t('shared.delete'), showDelay: 500 }"
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
        <!--suppress TypeScriptValidateTypes -->
        <InputText disabled :value="d(selectedDate, 'long')"></InputText>
      </InputGroup>
      <TrackSelection
        v-model="selectedEntry!.track"
        :linked-track="selectedEntry!.linkedTrack"
        @update:model-value="() => (dirty = true)"
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
      <div class="flex-shrink-1 flex-grow-1">
        <LeafletMap class="h-full w-full" :gpx-data="gpxData"></LeafletMap>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
