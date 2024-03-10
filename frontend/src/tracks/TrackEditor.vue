<script setup lang="ts">
import { getCurrentInstance, onMounted, ref, toRefs, watch } from "vue";
import * as L from "leaflet";
//@ts-expect-error
import * as E from "leaflet-editable/src/Leaflet.Editable";
import { tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";

const use = E;
import GpxData = tracks.GpxData;

const mapId = ref(`map${getCurrentInstance()?.uid}`);
let map: L.Map;
const mapContainer = ref();
const tracksApi = useTracksApi();

onMounted(() => {
  map = L.map(mapId.value, { editable: true }).setView([49, 9], 13);
  L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
  }).addTo(map);
  new ResizeObserver(invalidateSize).observe(mapContainer.value);
});

const props = defineProps<{
  gpxData: GpxData | undefined;
  editDirection: "forward" | "backward" | "drag";
}>();
const { gpxData, editDirection } = toRefs(props);

let trackLayer: L.Polyline | undefined = undefined;
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
  console.log(trackLayer);
  trackLayer.enableEdit();
  trackLayer.addEventListener({
    "editable:vertex:dragend": handleTrackEditEvent,
    "editable:vertex:new": handleTrackEditEvent,
    "editable:vertex:deleted": handleTrackEditEvent,
  });
  map.setView(
    L.latLng(gpxData.value!.waypoints[0].latitude, gpxData.value!.waypoints[0].longitude),
  );
  distanceMarkerLayer = L.layerGroup(
    gpxData.value?.distanceMarkers.map((dm) => createDistanceMarker(dm)),
  ).addTo(map);
});

async function handleTrackEditEvent() {
  const props = await tracksApi.ComputePolylineProps(
    (trackLayer!.getLatLngs() as L.LatLng[]).map((latLng) => ({
      latitude: latLng.lat,
      longitude: latLng.lng,
    })),
  );
  distanceMarkerLayer?.removeFrom(map);

  distanceMarkerLayer = L.layerGroup(props.distanceMarkers.map((dm) => createDistanceMarker(dm))).addTo(map);
  emit("change-track", props.length);
}

function createDistanceMarker(dm: tracks.DistanceMarker) {
  const icon = L.divIcon({html: `${(dm.distance / 1000).toFixed(0)}`, className: "distance-marker", iconAnchor: [18, 18]})
  return L.marker(L.latLng(dm.latitude, dm.longitude), {title: dm.distance.toString(), icon});
}

function invalidateSize() {
  map.invalidateSize();
}

watch(
  editDirection,
  (value) => {
    if (!trackLayer) {
      return;
    }
    console.log(value);
    //@ts-expect-error
    trackLayer.editor.disable();
    //@ts-expect-error
    trackLayer.editor.enable();
    if (value === "forward") {
      //@ts-expect-error
      trackLayer.editor.continueForward();
    } else if (value === "backward") {
      //@ts-expect-error
      trackLayer.editor.continueBackward();
    } else if (value === "drag") {
      //@ts-expect-error
      trackLayer.editor.reset();
    } else {
      console.error("invalid editDirection: ", value);
    }
  },
  { immediate: true },
);

const emit = defineEmits<{ (e: "change-track", length: number): void }>();
</script>

<template>
  <div ref="mapContainer" :id="mapId"></div>
</template>

<style>
.distance-marker {
  background-color: white;
  min-width: 36px;
  min-height: 36px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2em;
  border-color: #10b981;
  border-width: 4px;
  border-style: solid;
}
</style>
