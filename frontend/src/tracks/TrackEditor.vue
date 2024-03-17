<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import { tracks } from "../../wailsjs/go/models";
import { useMap } from "../shared/use-map";
import GpxData = tracks.GpxData;
import Coordinates = tracks.Coordinates;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
const mapContainer = ref();
const mapApi = useMap();

onMounted(() => {
  mapApi.initMap(mapId.value, mapContainer);
  mapApi.enableEditing(true, handleTrackEditEvent);
});

const props = defineProps<{
  gpxData: GpxData | undefined;
  editDirection: "forward" | "backward" | "drag";
}>();

const { gpxData, editDirection } = toRefs(props);

watch(gpxData, (value) => (mapApi.gpxData.value = value), { immediate: true });

async function handleTrackEditEvent(props: { length: number; waypoints: Coordinates[]}) {
  emit("change-track", props);
}

watch(editDirection, (value) => mapApi.changeEditDirection(value), { immediate: true });

const emit = defineEmits<{
  (e: "change-track", props: { length: number; waypoints: Coordinates[]}): void;
}>();
</script>

<template>
  <div ref="mapContainer" :id="mapId"></div>
</template>
