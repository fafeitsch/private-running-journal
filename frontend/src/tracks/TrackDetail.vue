<script setup lang="ts">
import { useRoute } from "vue-router";
import { computed, ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { trackEditor, tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import TrackEditor from "./TrackEditor.vue";
import Button from "primevue/button";
import ConfirmPopup from "primevue/confirmpopup";
import { useConfirm } from "primevue/useconfirm";
import { useLeaveConfirmation } from "../shared/use-leave-confirmation";
import MoveTrackOverlay from "./MoveTrackOverlay.vue";
import CreateTrackOverlay from "./CreateTrackOverlay.vue";
import SaveTrack = tracks.SaveTrack;
import Coordinates = tracks.Coordinates;
import TrackDto = trackEditor.TrackDto;
import CoordinateDto = trackEditor.CoordinateDto;

const route = useRoute();
const tracksStore = useTrackStore();
const { selectedTrack } = storeToRefs(tracksStore);
const { t, n } = useI18n();
const dirty = ref(false);

const track = ref<Omit<TrackDto, "convertValues"> | undefined>(undefined);
const tracksApi = useTracksApi();

watch(selectedTrack, (value) => {
  dirty.value = false;
  if (!value) {
    track.value = undefined;
    return;
  }
  track.value = { ...value };
});

const gpxData = ref<{
  waypoints: CoordinateDto[];
  distanceMarkers: (CoordinateDto & { distance: number })[];
}>({ waypoints: [], distanceMarkers: [] });
const editedWaypoints = ref<Coordinates[]>([]);
const length = ref(0);
const trackEditDirection = ref<"forward" | "drag" | "backward">("drag");

watch(
  () => route.params.trackId as string,
  async (trackId) => {
    if (trackId !== "new") {
      selectedTrack.value = await tracksApi.getTrack(trackId);
      gpxData.value = {
        waypoints: selectedTrack.value.waypoints,
        distanceMarkers: selectedTrack.value.distanceMarkers,
      };
    } else {
      selectedTrack.value = new TrackDto({
        id: "new",
        length: 0,
        name: "",
        usages: [],
        parents: [],
        waypoints: [],
        distanceMarkers: [],
      });
    }
    editedWaypoints.value = selectedTrack.value.waypoints;
    length.value = selectedTrack.value.length;
    trackEditDirection.value = "drag";
  },
  { immediate: true },
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
  if (track.value.usages.length > 0) {
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
    await tracksStore.loadTracks();
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
    await tracksStore.loadTracks();
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
        v-if="selectedTrack?.id !== 'new'"
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
        v-if="selectedTrack?.id !== 'new'"
        icon="pi pi-trash"
        @click="deleteTrack"
        v-tooltip="{ value: t('shared.delete'), showDelay: 500 }"
        :aria-label="t('shared.delete')"
      ></Button>
      <MoveTrackOverlay v-if="selectedTrack?.id !== 'new'"></MoveTrackOverlay>
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
        <InputText id="usagesInput" disabled :value="`${track!.usages.length}`"></InputText>
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
        v-if="selectedTrack"
        class="h-full w-full"
        :waypoints="selectedTrack.waypoints"
        :polyline-meta="selectedTrack"
        :edit-direction="trackEditDirection"
        @change-track="trackChanged"
      ></TrackEditor>
    </div>
  </div>
</template>

<style scoped></style>
