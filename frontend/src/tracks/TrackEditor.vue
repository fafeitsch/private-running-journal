<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import { tracks } from "../../wailsjs/go/models";
import { useMap } from "../shared/use-map";
import GpxData = tracks.GpxData;
import Coordinates = tracks.Coordinates;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
const mapContainer = ref();
const mapApi = useMap();

const props = defineProps<{
  gpxData: GpxData;
  editDirection: "forward" | "backward" | "drag";
}>();

const { gpxData, editDirection } = toRefs(props);

onMounted(() => {
  mapApi.initMap(mapId.value, mapContainer);
  mapApi.enableEditing(true, handleTrackEditEvent);
  mapApi.gpxData.value = gpxData.value;
  setTimeout(() => mapApi.changeEditDirection(editDirection.value));
});


watch(gpxData, (value) => (mapApi.gpxData.value = value));

async function handleTrackEditEvent(props: { length: number; waypoints: Coordinates[] }) {
  emit("change-track", props);
}

watch(editDirection, (value) => mapApi.changeEditDirection(value));

const emit = defineEmits<{
  (e: "change-track", props: { length: number; waypoints: Coordinates[] }): void;
}>();
</script>

<template>
  <div ref="mapContainer" :id="mapId" style="z-index: 0"></div>
</template>
