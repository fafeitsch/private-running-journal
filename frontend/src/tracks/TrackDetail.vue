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
import Button from "primevue/button";
import ConfirmPopup from "primevue/confirmpopup";
import { useConfirm } from "primevue/useconfirm";
import { useLeaveConfirmation } from "../shared/use-leave-confirmation";
import MoveTrackOverlay from "./MoveTrackOverlay.vue";
import CreateTrackOverlay from "./CreateTrackOverlay.vue";
import GpxData = tracks.GpxData;
import SaveTrack = tracks.SaveTrack;
import Coordinates = tracks.Coordinates;

const route = useRoute();
const tracksStore = useTrackStore();
const { selectedTrackId, selectedTrack } = storeToRefs(tracksStore);
const { t, n } = useI18n();
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

const gpxData = ref<GpxData>(new GpxData({waypoints: [], distanceMarkers: []}));
const tracksApi = useTracksApi();
const editedWaypoints = ref<Coordinates[]>([]);
const trackEditDirection = ref<"forward" | "drag" | "backward">("drag");

const length = ref(0);

watch(
  selectedTrack,
  async () => {
    if (!selectedTrack.value || selectedTrackId.value === "new") {
      gpxData.value = new GpxData({ waypoints: [], distanceMarkers: [] });
      editedWaypoints.value = [];
      trackEditDirection.value = 'forward'
      length.value = 0;
      return;
    }
    try {
      gpxData.value = await tracksApi.getGpxData(selectedTrack.value.id);
      editedWaypoints.value = gpxData.value.waypoints;
      length.value = selectedTrack.value.length;
      trackEditDirection.value = 'drag'
    } catch (e) {
      console.error(e);
    }
  },
  { deep: true, immediate: true },
);

const formattedLength = computed(() =>
  n(length.value / 1000, { maximumFractionDigits: 1, minimumFractionDigits: 1 }),
);

function trackChanged(props: { length: number; waypoints: Coordinates[] }) {
  length.value = props.length;
  editedWaypoints.value = props.waypoints;
  dirty.value = true;
}

const confirm = useConfirm();

async function saveTrack(event: any) {
  trackEditDirection.value = "drag";
  if (!track.value) {
    return;
  }
  let choice = true;
  if (track.value.usages > 0) {
    let resolveFn: (result: boolean) => void;
    const result = new Promise<boolean>((resolve) => (resolveFn = resolve));
    confirm.require({
      target: event.currentTarget,
      group: "track",
      header: t("shared.confirm.header"),
      accept: () => resolveFn(true),
      reject: () => resolveFn(false),
      message: t("tracks.changeJournalWarning", { count: track.value.usages }),
      rejectLabel: t("shared.cancel"),
      acceptLabel: t("shared.save"),
    });
    choice = await result;
  }
  if (!choice) {
    return;
  }
  try {
    const updated = await tracksApi.saveTrack(
      new SaveTrack({
        id: track.value.id,
        name: track.value.name,
        parents: [],
        waypoints: editedWaypoints.value,
      }),
    );
    dirty.value = false;
    tracksStore.updateTrack(updated);
  } catch (e) {
    // TODO error handling
    console.error(e);
  }
}

async function deleteTrack(event: Event) {
  if (!track.value) {
    return;
  }
  let resolveFn: (result: boolean) => void;
  const result = new Promise<boolean>((resolve) => (resolveFn = resolve));
  confirm.require({
    target: event.currentTarget as HTMLElement,
    group: "track",
    header: t("shared.confirm.header"),
    accept: () => resolveFn(true),
    reject: () => resolveFn(false),
    message: t("tracks.deleteConfirmation", { count: track.value.usages }),
    rejectLabel: t("shared.cancel"),
    acceptLabel: t("shared.delete"),
  });
  let choice = await result;
  if (!choice) {
    return;
  }
  try {
    await tracksApi.deleteTrack(track.value.id);
    tracksStore.deleteTrack(track.value.id);
  } catch (e) {
    console.error(e);
  }
}

useLeaveConfirmation(dirty);
</script>

<template>
  <div v-if="track" class="w-full p-2 flex flex-col h-full gap-2">
    <div class="flex gap-2">
      <Button
        v-if="selectedTrackId !== 'new'"
        icon="pi pi-save"
        :disabled="!dirty"
        @click="saveTrack"
        v-tooltip="{ value: t('shared.save'), showDelay: 500 }"
        :aria-label="t('shared.save')"
      ></Button>
      <CreateTrackOverlay
        v-else
        :name="track!.name"
        :waypoints="editedWaypoints"
        @track-created="dirty = false"
      ></CreateTrackOverlay>
      <ConfirmPopup group="track">
        <template #message="{ message }">
          <div style="max-width: 330px" class="p-2" data-testid="delete-track-confirmation">
            {{ message.message }}
          </div>
        </template>
      </ConfirmPopup>
      <Button
        v-if="selectedTrackId !== 'new'"
        icon="pi pi-trash"
        @click="deleteTrack"
        v-tooltip="{ value: t('shared.delete'), showDelay: 500 }"
        :aria-label="t('shared.delete')"
      ></Button>
      <MoveTrackOverlay v-if="selectedTrackId !== 'new'"></MoveTrackOverlay>
    </div>
    <div class="flex gap-2">
      <InputGroup>
        <InputGroupAddon>
          <label for="nameInput">{{ t("tracks.name") }}</label>
        </InputGroupAddon>
        <InputText
          id="nameInput"
          v-model="track!.name"
          @update:model-value="dirty = true"
        ></InputText>
      </InputGroup>
      <InputGroup class="grow-0 w-2/12">
        <InputGroupAddon>
          <label for="usagesInput">{{ t("tracks.usages") }}</label>
        </InputGroupAddon>
        <InputText id="usagesInput" disabled :value="`${track!.usages}`"></InputText>
      </InputGroup>
    </div>
    <div class="flex gap-2 items-center">
      <InputGroup class="shrink grow w-auto">
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
        :allow-empty="false"
      ></SelectButton>
    </div>
    <div class="shrink grow">
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
