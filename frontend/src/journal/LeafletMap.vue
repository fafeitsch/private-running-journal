<script setup lang="ts">
import { getCurrentInstance, onMounted, ref } from "vue";
import * as L from 'leaflet'

const mapId = ref(`map${getCurrentInstance()?.uid}`);
let map: L.Map
const mapContainer = ref()

onMounted(() => {
  map = L.map(mapId.value).setView([49, 9], 13);
  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map);
  new ResizeObserver(invalidateSize).observe(mapContainer.value)
});

function invalidateSize() {
  map.invalidateSize()
}
</script>

<template>
  <div ref="mapContainer" class="h-full w-full" :id="mapId"></div>
</template>

<style scoped></style>
