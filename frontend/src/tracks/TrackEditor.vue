<script setup lang="ts">
import {getCurrentInstance, onMounted, ref, toRefs, watch, watchEffect} from "vue";
import { trackEditor, tracks } from "../../wailsjs/go/models";
import { useMap } from "../shared/use-map";
import Coordinates = tracks.Coordinates;
import CoordinateDto = trackEditor.CoordinateDto;
import PolylineMeta = trackEditor.PolylineMeta;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
const mapContainer = ref();
const mapApi = useMap();

const props = defineProps<{
  polylineMeta: PolylineMeta;
  waypoints: CoordinateDto[];
  editDirection: "forward" | "backward" | "drag";
}>();

const { polylineMeta, waypoints, editDirection } = toRefs(props);

onMounted(() => {
  mapApi.initMap(mapId.value, mapContainer);
  mapApi.enableEditing(true, handleTrackEditEvent);
  mapApi.waypoints.value = waypoints.value;
  mapApi.polylineMeta.value = polylineMeta.value;
  setTimeout(() => mapApi.changeEditDirection(editDirection.value));
});


watch(polylineMeta, (value) => (mapApi.polylineMeta.value = value));
watch(waypoints, (value) => (mapApi.waypoints.value = value));

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
