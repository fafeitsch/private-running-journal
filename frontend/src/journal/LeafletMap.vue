<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import * as L from "leaflet";
import { tracks} from "../../wailsjs/go/models";
import GpxData = tracks.GpxData;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
let map: L.Map;
const mapContainer = ref();

onMounted(() => {
  map = L.map(mapId.value).setView([49, 9], 13);
  L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
  }).addTo(map);
  new ResizeObserver(invalidateSize).observe(mapContainer.value);
});

const props = defineProps<{ gpxData: GpxData | undefined }>();
const { gpxData } = toRefs(props);

let trackLayer: L.Layer | undefined = undefined;
let distanceMarkerLayer: L.Layer | undefined = undefined;

watch(gpxData, () => {
  trackLayer?.removeFrom(map);
  distanceMarkerLayer?.removeFrom(map);
  if (!gpxData.value) {
    return;
  }
  trackLayer = L.polyline(
    gpxData.value!.waypoints.map((wp) => L.latLng(wp.latitude, wp.longitude)),
    { color: "red" },
  ).addTo(map);

  distanceMarkerLayer = L.layerGroup(
    gpxData.value?.distanceMarkers.map((dm) => L.marker(L.latLng(dm.latitude, dm.longitude), {title: dm.distance.toString()})),
  ).addTo(map);
});

function invalidateSize() {
  map.invalidateSize();
}
</script>

<template>
  <div ref="mapContainer" class="h-full w-full" :id="mapId"></div>
</template>

<style scoped></style>
