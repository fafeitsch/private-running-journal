<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRoute, useRouter } from "vue-router";
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
import Calendar from "primevue/calendar";
import LeafletMap from "./LeafletMap.vue";
import TrackTimeResult from "./TrackTimeResult.vue";
import TrackSelection from "./TrackSelection.vue";
import ConfirmDialog from "primevue/confirmdialog";
import { useConfirm } from "primevue/useconfirm";
import { useJournalStore } from "../store/journal-store";
import { storeToRefs } from "pinia";

const { t, d, locale } = useI18n();
const route = useRoute();
const journalApi = useJournalApi();
const router = useRouter();

const loading = ref(false);
const error = ref(false);
const dirty = ref(false);
const selectedEntry = ref<journal.Entry | undefined>(undefined);
const selectedDate = ref<Date>(new Date());
const journalStore = useJournalStore();
const { selectedEntryId } = storeToRefs(journalStore);

const tracksApi = useTracksApi();

watch(
  () => route.params.entryId as string,
  (value: string) => loadEntry(value),
  { immediate: true },
);

async function loadEntry(entryId: string) {
  selectedEntry.value = undefined;
  error.value = false;
  loading.value = true;
  selectedEntryId.value = entryId;
  if (!entryId) {
    await router.replace("/journal");
    return;
  }
  try {
    selectedEntry.value = await journalApi.getListEntry(entryId);
    selectedDate.value = new Date(Date.parse(selectedEntry.value.date));
  } catch (e) {
    console.error(e);
    error.value = true;
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
    }
  },
  { deep: true },
);

async function saveEntry() {
  if (!selectedEntry.value) {
    return;
  }
  try {
    await journalApi.saveEntry(selectedEntry.value);
    dirty.value = false;
  } catch (e) {
    console.error(e);
  }
}

const confirm = useConfirm();

onBeforeRouteLeave(() => handleRouteLeave());
onBeforeRouteUpdate(() => handleRouteLeave());

function handleRouteLeave(): Promise<boolean> {
  if (!dirty.value) {
    return Promise.resolve(true);
  }
  let resolveFn: (result: boolean) => void;
  const result = new Promise<boolean>((resolve) => (resolveFn = resolve));
  confirm.require({
    header: t("shared.confirm.header"),
    accept: () => resolveFn(true),
    reject: () => resolveFn(false),
    message: t("shared.confirm.message"),
    rejectLabel: t("shared.cancel"),
    acceptLabel: t("shared.confirm.discard"),
  });
  return result;
}
</script>

<template>
  <div class="flex flex-column">
    <div v-if="loading" class="flex w-full flex-grow-1 justify-content-center align-items-center">
      <ProgressSpinner></ProgressSpinner>
    </div>
    <div v-else-if="error">
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
      <div class="flex">
        <Button icon="pi pi-save" :disabled="!dirty" @click="saveEntry"></Button>
        <ConfirmDialog></ConfirmDialog>
      </div>
      <InputGroup>
        <InputGroupAddon>
          <label for="date">{{ t("journal.details.date") }}</label>
        </InputGroupAddon>
        <Calendar
          id="date"
          v-model="selectedDate"
          :date-format="locale === 'de' ? 'dd.mm.yy' : 'yyyy/mm/dd'"
          :disabled="true"
        ></Calendar>
      </InputGroup>
      <TrackSelection
        v-model="selectedEntry!.track"
        :linked-track="selectedEntry!.linkedTrack"
        @update:model-value="() => (dirty = true)"
      ></TrackSelection>
      <TrackTimeResult
        v-model:laps="selectedEntry!.laps"
        v-model:time="selectedEntry!.time"
        :track-length="selectedEntry!.track?.length || undefined"
        @update:laps="() => (dirty = true)"
        @update:time="() => (dirty = true)"
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
        <LeafletMap class="h-full w-full":gpx-data="gpxData"></LeafletMap>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
