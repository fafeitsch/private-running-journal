import * as L from "leaflet";
import { MaybeRefOrGetter, ref, Ref, toValue, watch } from "vue";
import { tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
// @ts-expect-error
// noinspection ES6UnusedImports needed to make editable work
import * as E from "leaflet-editable/src/Leaflet.Editable";
import { useSettingsStore } from "../store/settings-store";
import GpxData = tracks.GpxData;
import Coordinates = tracks.Coordinates;

const use = E;

export const useMap = () => {
  let map: L.Map | undefined = undefined;
  const settingsStore = useSettingsStore();

  let gpxData = ref<GpxData>();

  function initMap(id: MaybeRefOrGetter<string>, mapContainer: Ref) {
    const mapSettings = settingsStore.settings.mapSettings;
    //@ts-expect-error
    map = L.map(toValue(id), { editable: true }).setView(
      mapSettings.center as [number, number],
      mapSettings.zoomLevel,
    );
    L.tileLayer(`http://127.0.0.1:${settingsStore.settings.httpPort}/tiles/{z}/{x}/{y}`, {
      maxZoom: 19,
      attribution: mapSettings.attribution,
    }).addTo(map);
    new ResizeObserver(() => map?.invalidateSize()).observe(mapContainer.value);
  }

  function createDistanceMarker(dm: tracks.DistanceMarker) {
    const icon = L.divIcon({
      html: `<span data-testid="distance-marker">${(dm.distance / 1000).toFixed(0)}</span>`,
      className: "distance-marker",
      iconAnchor: [18, 18],
    });
    return L.marker(L.latLng(dm.latitude, dm.longitude), { title: dm.distance.toString(), icon });
  }

  let trackLayer: L.Polyline | undefined = undefined;
  let distanceMarkerLayer: L.Layer | undefined = undefined;

  watch(gpxData, (newValue) => {
    if (!map) {
      return;
    }
    trackLayer?.removeFrom(map);
    distanceMarkerLayer?.removeFrom(map);
    if (!newValue) {
      return;
    }
    trackLayer = L.polyline(
      newValue.waypoints.map((wp) => L.latLng(wp.latitude, wp.longitude)),
      { color: "red" },
    ).addTo(map);
    enableEditing();
    if (newValue.waypoints.length) {
      map.setView(
        L.latLng(gpxData.value!.waypoints[0].latitude, gpxData.value!.waypoints[0].longitude),
      );
    }
    distanceMarkerLayer = L.layerGroup(
      newValue.distanceMarkers.map((dm) => createDistanceMarker(dm)),
    ).addTo(map);
  });

  let editEnabled = false;
  let editTrackHandler: (props: { length: number; waypoints: Coordinates[] }) => void = () =>
    void 0;

  function createTrackLayerIfNecessary() {
    if(!map) {
      throw new Error('map is not initialized yet')
    }
    if (!trackLayer) {
      trackLayer = L.polyline([], { color: "red" }).addTo(map);
    }
  }

  function enableEditing(
    value: boolean = editEnabled,
    handler: (props: { length: number; waypoints: Coordinates[] }) => void = editTrackHandler,
  ) {
    editTrackHandler = handler;
    editEnabled = value;
    createTrackLayerIfNecessary();
    trackLayer!.removeEventListener({
      //@ts-expect-error
      "editable:vertex:dragend": handleTrackEditEvent,
      "editable:vertex:new": handleTrackEditEvent,
      "editable:vertex:deleted": handleTrackEditEvent,
    });
    if (value) {
      //@ts-expect-error
      trackLayer!.enableEdit();
      trackLayer!.addEventListener({
        //@ts-expect-error
        "editable:vertex:dragend": handleTrackEditEvent,
        "editable:vertex:new": handleTrackEditEvent,
        "editable:vertex:deleted": handleTrackEditEvent,
      });
    } else {
      //@ts-expect-error
      trackLayer!.disableEdit();
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
    createTrackLayerIfNecessary()
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
