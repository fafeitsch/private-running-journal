<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import { tracks } from "../../wailsjs/go/models";
import { useMap } from "../shared/use-map";
import GpxData = tracks.GpxData;

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

async function handleTrackEditEvent(length: number) {
  emit("change-track", length);
}

watch(editDirection, (value) => mapApi.changeEditDirection(value), { immediate: true });

const emit = defineEmits<{ (e: "change-track", length: number): void }>();
</script>

<template>
  <div ref="mapContainer" :id="mapId"></div>
</template>
