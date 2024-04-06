import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useTracksApi } from "../api/tracks";
import { tracks } from "../../wailsjs/go/models";
import Track = tracks.Track;

export const useTrackStore = defineStore("tacks", () => {
  const trackApi = useTracksApi();
  const availableTracks = ref<Track[]>([]);
  const selectedTrackId = ref<string | undefined>(undefined);

  async function loadTracks() {
    availableTracks.value = await trackApi.getTracks();
  }

  const selectedTrack = computed(() => {
    if (!selectedTrackId.value) {
      return undefined;
    }

    return availableTracks.value.find((t) => t.id === selectedTrackId.value);
  });

  function addTrack(track: Track) {
    availableTracks.value.push(track);
  }

  function updateTrack(track: Track) {
    const existing = availableTracks.value.find((t) => t.id === track.id);
    if (!existing) {
      return;
    }
    existing.name = track.name;
    existing.length = track.length;
    existing.usages = track.usages;
  }

  function deleteTrack(id: string) {
    availableTracks.value = availableTracks.value.filter((track) => track.id !== id);
  }

  return {
    loadTracks,
    availableTracks,
    selectedTrackId,
    selectedTrack,
    addTrack,
    updateTrack,
    deleteTrack,
  };
});
