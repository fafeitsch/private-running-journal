<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { computed, onMounted, ref, watch } from "vue";
import { useJournalApi } from "../api/journal";
import ProgressSpinner from "primevue/progressspinner";
import Message from "primevue/message";
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputGroup from "primevue/inputgroup";
import TreeSelect from "primevue/treeselect";
import { journal, tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import { TreeNode } from "primevue/treenode";
import Calendar from "primevue/calendar";
import LeafletMap from "./LeafletMap.vue";
import TrackTimeResult from "./TrackTimeResult.vue";

const { t, d, locale } = useI18n();
const route = useRoute();
const journalApi = useJournalApi();
const router = useRouter();

const loading = ref(false);
const error = ref(false);
const selectedEntryId = computed(() => route.params["entryId"]);
const selectedEntry = ref<journal.Entry | undefined>(undefined);
const selectedDate = ref<Date>(new Date());
const availableTracks = ref<tracks.Track[]>([]);

const selectedTrack = ref<Record<string, boolean>>({});

const tracksApi = useTracksApi();

onMounted(async () => {
  try {
    availableTracks.value = await tracksApi.getTracks();
    console.log(availableTracks.value);
  } catch (e) {
    // todo error handling
    console.error(e);
  }
});

watch(selectedEntryId, () => loadEntry(), { immediate: true });

async function loadEntry() {
  selectedEntry.value = undefined;
  error.value = false;
  loading.value = true;
  if (!selectedEntryId.value || typeof selectedEntryId.value !== "string") {
    await router.replace("/journal");
    return;
  }
  try {
    selectedEntry.value = await journalApi.getListEntry(selectedEntryId.value);
    if (!selectedEntry.value.track) {
      console.error("no track found ");
      return;
    }
    selectedTrack.value = { [selectedEntry.value.track.id]: true };
    selectedDate.value = new Date(Date.parse(selectedEntry.value.date));
  } catch (e) {
    console.error(e);
    error.value = true;
  } finally {
    loading.value = false;
  }
}

const trackSelection = computed<TreeNode[]>(() => {
  const trackToListEntry: (tracks: tracks.Track, parentNames: string) => TreeNode = (
    track: tracks.Track,
    parentNames: string,
  ) => {
    const name = parentNames ? `${parentNames} / ${track.name}` : track.name;
    return {
      key: track.id,
      label: track.name,
      data: track,
      children: track.variants.map((entry) => trackToListEntry(entry, name)),
      selectable: track.length > 0,
      selectedLabel: name,
    };
  };

  return availableTracks.value.map((entry) => trackToListEntry(entry, ""));
});

const gpxData = ref<tracks.GpxData | undefined>(undefined);

watch(
  selectedEntry,
  async () => {
    if (!selectedEntry.value) {
      gpxData.value = undefined;
      return;
    }
    try {
      gpxData.value = await tracksApi.getGpxData(selectedEntry.value.track!.id);
    } catch (e) {
      console.error(e);
    }
  },
  { deep: true },
);
</script>

<template>
  <div class="flex flex-column">
    <div v-if="loading">
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
            @click="loadEntry"
          ></Button></div
      ></Message>
    </div>
    <div v-else-if="selectedEntry" class="flex flex-column gap-2 w-full p-2">
      <InputGroup>
        <InputGroupAddon>
          <label for="date">{{ t("journal.details.date") }}</label>
        </InputGroupAddon>
        <Calendar
          id="date"
          v-model="selectedDate"
          :date-format="locale === 'de' ? 'dd.mm.yy' : 'yyyy/mm/dd'"
        ></Calendar>
      </InputGroup>
      <TrackTimeResult
        v-model:laps="selectedEntry!.laps"
        v-model:time="selectedEntry!.time"
        :track-length="selectedEntry!.track!.length"
      ></TrackTimeResult>
      <InputGroup>
        <InputGroupAddon>
          <label for="track">{{ t("journal.details.track") }}</label>
        </InputGroupAddon>
        <TreeSelect
          id="track"
          v-model="selectedTrack"
          selection-mode="single"
          :options="trackSelection"
          placeholder="Select Item"
          class="md:w-20rem w-full"
          @node-select="(node) => (selectedEntry!.track = node.data)"
        >
          <template #value="props">
            {{ props.value[0]?.selectedLabel }}
          </template>
        </TreeSelect>
      </InputGroup>
      <InputGroup>
        <InputGroupAddon>
          <label for="comment">{{ t("journal.details.comment") }}</label>
        </InputGroupAddon>
        <InputText id="comment" v-model="selectedEntry!.comment"></InputText>
      </InputGroup>
    </div>
    <LeafletMap :gpx-data="gpxData" class="flex-grow-1"></LeafletMap>
  </div>
</template>

<style scoped></style>
