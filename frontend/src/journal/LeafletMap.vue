<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import * as L from "leaflet";
import { tracks } from "../../wailsjs/go/models";
import { useMap } from "../shared/use-map";
import GpxData = tracks.GpxData;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
let map: L.Map;
const mapContainer = ref();
const mapApi = useMap();

onMounted(() => {
  mapApi.initMap(mapId, mapContainer);
});

const props = defineProps<{ gpxData: GpxData | undefined }>();
const { gpxData } = toRefs(props);

watch(gpxData, (value) => (mapApi.gpxData.value = value));
</script>

<template>
  <div ref="mapContainer" :id="mapId" style="z-index: 0"></div>
</template>

