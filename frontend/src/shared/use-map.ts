import * as L from "leaflet";
import { MaybeRefOrGetter, ref, Ref, toValue, watch } from "vue";
import { tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
// @ts-expect-error
// noinspection ES6UnusedImports needed to make editable work
import * as E from "leaflet-editable/src/Leaflet.Editable";
import GpxData = tracks.GpxData;
import Coordinates = tracks.Coordinates;

const use = E;

export const useMap = () => {
  let map: L.Map | undefined = undefined;

  let gpxData = ref<GpxData | undefined>();

  function initMap(id: MaybeRefOrGetter<string>, mapContainer: Ref) {
    map = L.map(toValue(id), { editable: true }).setView([49, 9], 13);
    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 19,
      attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    }).addTo(map);
    new ResizeObserver(() => map?.invalidateSize()).observe(mapContainer.value);
  }

  function createDistanceMarker(dm: tracks.DistanceMarker) {
    const icon = L.divIcon({
      html: `${(dm.distance / 1000).toFixed(0)}`,
      className: "distance-marker",
      iconAnchor: [18, 18],
    });
    return L.marker(L.latLng(dm.latitude, dm.longitude), { title: dm.distance.toString(), icon });
  }

  let trackLayer: L.Polyline | undefined = undefined;
  let distanceMarkerLayer: L.Layer | undefined = undefined;

  watch(gpxData, () => {
    if (!map) {
      return;
    }
    trackLayer?.removeFrom(map);
    distanceMarkerLayer?.removeFrom(map);
    if (!gpxData.value) {
      return;
    }
    trackLayer = L.polyline(
      gpxData.value!.waypoints.map((wp) => L.latLng(wp.latitude, wp.longitude)),
      { color: "red" },
    ).addTo(map);
    console.log(trackLayer)
    enableEditing();
    if(gpxData.value?.waypoints.length) {
      map.setView(
        L.latLng(gpxData.value!.waypoints[0].latitude, gpxData.value!.waypoints[0].longitude),
      );
    }
    distanceMarkerLayer = L.layerGroup(
      gpxData.value?.distanceMarkers.map((dm) => createDistanceMarker(dm)),
    ).addTo(map);
  });

  let editEnabled = false;
  let editTrackHandler: (props: { length: number; waypoints: Coordinates[] }) => void = () =>
    void 0;

  function enableEditing(
    value: boolean = editEnabled,
    handler: (props: { length: number; waypoints: Coordinates[] }) => void = editTrackHandler,
  ) {
    editTrackHandler = handler;
    editEnabled = value;
    if (!trackLayer) {
      console.log('no track layer')
      return;
    }
    trackLayer.removeEventListener({
      "editable:vertex:dragend": handleTrackEditEvent,
      "editable:vertex:new": handleTrackEditEvent,
      "editable:vertex:deleted": handleTrackEditEvent,
    });
    if (value) {
      trackLayer.enableEdit();
      trackLayer.addEventListener({
        "editable:vertex:dragend": handleTrackEditEvent,
        "editable:vertex:new": handleTrackEditEvent,
        "editable:vertex:deleted": handleTrackEditEvent,
      });
    } else {
      trackLayer.disableEdit();
    }
  }

  const tracksApi = useTracksApi();

  async function handleTrackEditEvent() {
    if (!map) {
      return;
    }
    const coordinates = (trackLayer!.getLatLngs() as L.LatLng[]).map((latLng) => ({
      latitude: latLng.lat,
      longitude: latLng.lng,
    }));
    const props = await tracksApi.ComputePolylineProps(coordinates);
    distanceMarkerLayer?.removeFrom(map);

    distanceMarkerLayer = L.layerGroup(
      props.distanceMarkers.map((dm) => createDistanceMarker(dm)),
    ).addTo(map);
    editTrackHandler({
      length: props.length,
      waypoints: coordinates,
    });
  }

  function changeEditDirection(value: "backward" | "drag" | "forward") {
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
  }

  return { initMap, gpxData, enableEditing, changeEditDirection };
};
