<script setup lang="ts">
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRoute } from "vue-router";
import { computed, ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import TrackEditor from "./TrackEditor.vue";
import Button from "primevue/button";
import ConfirmDialog from "primevue/confirmdialog";
import ConfirmPopup from "primevue/confirmpopup";
import { useConfirm } from "primevue/useconfirm";
import GpxData = tracks.GpxData;
import SaveTrack = tracks.SaveTrack;
import Coordinates = tracks.Coordinates;

const route = useRoute();
const tracksStore = useTrackStore();
const { selectedTrackId, selectedTrack } = storeToRefs(tracksStore);
const { t } = useI18n();
const dirty = ref(false);

const track = ref<Omit<tracks.Track, "convertValues"> | undefined>(undefined);

watch(selectedTrack, (value) => {
  dirty.value = false;
  if (!value) {
    track.value = undefined;
    return;
  }
  track.value = { ...value };
});

watch(
  () => route.params.trackId as string,
  (trackId) => {
    selectedTrackId.value = trackId;
  },
  { immediate: true },
);

const prefix = computed(() => {
  if (!track) {
    return "";
  }
  return track.value?.parentNames.join("/") || "";
});

const gpxData = ref<GpxData | undefined>(undefined);
const tracksApi = useTracksApi();
const editedWaypoints = ref<Coordinates[]>([]);

watch(
  selectedTrack,
  async () => {
    if (!selectedTrack.value) {
      gpxData.value = undefined;
      editedWaypoints.value = [];
      return;
    }
    try {
      gpxData.value = await tracksApi.getGpxData(selectedTrack.value.id);
      editedWaypoints.value = gpxData.value.waypoints;
      length.value = selectedTrack.value.length;
    } catch (e) {
      console.error(e);
    }
  },
  { deep: true },
);

const length = ref(0);
const formattedLength = computed(() => (length.value / 1000).toFixed(1));

function trackChanged(props: { length: number; waypoints: Coordinates[] }) {
  length.value = props.length;
  editedWaypoints.value = props.waypoints;
  dirty.value = true;
}

const trackEditDirection = ref<"forward" | "drag" | "backward">("drag");

async function saveTrack(event: any) {
  if (!track.value) {
    return;
  }
  let choice = true;
  if (track.value.usages > 0) {
    let resolveFn: (result: boolean) => void;
    const result = new Promise<boolean>((resolve) => (resolveFn = resolve));
    confirm.require({
      target: event.currentTarget,
      group: "journal",
      header: t("shared.confirm.header"),
      accept: () => resolveFn(true),
      reject: () => resolveFn(false),
      message: t("tracks.changeJournalWarning", { count: track.value.usages }),
      rejectLabel: t("shared.cancel"),
      acceptLabel: t("shared.save"),
    });
    choice = await result;
  }
  if (choice) {
    try {
      await tracksApi.saveTrack(
        new SaveTrack({
          id: track.value.id,
          name: track.value.name,
          parents: [],
          waypoints: editedWaypoints.value,
        }),
      );
    } catch (e) {
      // TODO error handling
      console.error(e);
    }
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
    group: "leave",
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
  <div v-if="track" class="w-full p-2 flex flex-column h-full gap-2">
    <div class="flex">
      <Button icon="pi pi-save" :disabled="!dirty" @click="saveTrack"></Button>
      <ConfirmPopup group="journal">
        <template #message="{ message }">
          <div style="max-width: 330px" class="p-2">{{ message.message }}</div>
        </template>
      </ConfirmPopup>
      <ConfirmDialog group="leave"></ConfirmDialog>
    </div>
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
        <InputText
          id="nameInput"
          v-model="track!.name"
          @update:model-value="dirty = true"
        ></InputText>
      </InputGroup>
      <InputGroup class="flex-grow-0 w-2">
        <InputGroupAddon>
          <label for="nameInput">{{ t("tracks.usages") }}</label>
        </InputGroupAddon>
        <InputText id="nameInput" v-model="track!.usages"></InputText>
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
