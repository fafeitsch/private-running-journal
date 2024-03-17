<script setup lang="ts">
import { useRoute } from "vue-router";
import { computed, ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import TrackEditor from "./TrackEditor.vue";
import GpxData = tracks.GpxData;

const route = useRoute();
const tracksStore = useTrackStore();
const { selectedTrackId, selectedTrack } = storeToRefs(tracksStore);
const { t } = useI18n();

watch(
  () => route.params.trackId as string,
  (trackId) => {
    selectedTrackId.value = trackId;
  },
  { immediate: true },
);

const prefix = computed(() => {
  if (!selectedTrack) {
    return "";
  }
  return selectedTrack.value?.parentNames.join("/") || "";
});

const gpxData = ref<GpxData | undefined>(undefined);
const tracksApi = useTracksApi();

watch(
  selectedTrack,
  async () => {
    if (!selectedTrack.value) {
      gpxData.value = undefined;
      return;
    }
    try {
      gpxData.value = await tracksApi.getGpxData(selectedTrack.value.id);
      length.value = selectedTrack.value.length;
    } catch (e) {
      console.error(e);
    }
  },
  { deep: true },
);

const length = ref(0);
const formattedLength = computed(() => (length.value / 1000).toFixed(1));

function trackChanged(newLength: number) {
  length.value = newLength;
}

const trackEditDirection = ref<"forward" | "drag" | "backward">("drag");
</script>

<template>
  <div v-if="selectedTrack" class="w-full p-2 flex flex-column h-full gap-2">
    <div class="flex gap-2">
      <InputGroup>
        <InputGroupAddon>
          <label for="nameInput">{{ t("tracks.name") }}</label>
        </InputGroupAddon>
        <InputText
          v-if="prefix"
          :value="prefix"
          class="flex-shrink-0 flex-grow-0 w-auto"
          disabled
          :pt="{ root: { size: prefix.length } }"
        ></InputText>
        <InputText id="nameInput" v-model="selectedTrack!.name"></InputText>
      </InputGroup>
      <InputGroup class="flex-grow-0 w-2">
        <InputGroupAddon>
          <label for="nameInput">{{ t("tracks.usages") }}</label>
        </InputGroupAddon>
        <InputText id="nameInput" v-model="selectedTrack!.usages"></InputText>
      </InputGroup>
    </div>
    <div class="flex gap-2 align-items-center">
      <InputGroup class="flex-shrink-1 flex-grow-1 w-auto">
        <InputGroupAddon>
          <label for="lengthInput">{{ t("journal.details.length") }}</label>
        </InputGroupAddon>
        <InputText id="lengthInput" v-model="formattedLength" disabled></InputText>
        <InputGroupAddon>km</InputGroupAddon>
      </InputGroup>
      <label for="editModeInput">{{ t("tracks.editMode") }}</label>
      <SelectButton
        id="editModeInput"
        :options="[
          { name: t('tracks.editModeBackward'), value: 'backward' },
          { name: t('tracks.editModeDrag'), value: 'drag' },
          { name: t('tracks.editModeForward'), value: 'forward' },
        ]"
        v-model="trackEditDirection"
        option-value="value"
        option-label="name"
      ></SelectButton>
    </div>
    <div class="flex-shrink-1 flex-grow-1">
      <TrackEditor
        class="h-full w-full"
        :gpx-data="gpxData"
        :edit-direction="trackEditDirection"
        @change-track="trackChanged"
      ></TrackEditor>
    </div>
  </div>
</template>

<style scoped></style>
