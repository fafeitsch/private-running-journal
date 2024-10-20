<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import * as L from "leaflet";
import { trackEditor } from "../../wailsjs/go/models";
import { useMap } from "../shared/use-map";
import CoordinateDto = trackEditor.CoordinateDto;
import PolylineMeta = trackEditor.PolylineMeta;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
let map: L.Map;
const mapContainer = ref();
const mapApi = useMap();

onMounted(() => {
  mapApi.initMap(mapId, mapContainer);
  mapApi.waypoints.value = waypoints.value;
  mapApi.polylineMeta.value = polylineMeta.value;
});

const props = withDefaults(
  defineProps<{ waypoints: CoordinateDto[]; polylineMeta: PolylineMeta }>(),
  { waypoints: () => [], polylineMeta: () => new PolylineMeta({ length: 0, distanceMarkers: [] }) },
);
const { waypoints, polylineMeta } = toRefs(props);

watch(waypoints, (value) => (mapApi.waypoints.value = value));
watch(polylineMeta, (value) => (mapApi.polylineMeta.value = value));
</script>

<template>
  <div ref="mapContainer" :id="mapId" style="z-index: 0"></div>
</template>
