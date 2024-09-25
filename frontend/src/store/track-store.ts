import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useTracksApi } from "../api/tracks";
import {trackEditor, tracks} from "../../wailsjs/go/models";
import Track = tracks.Track;
import TrackTreeNode = tracks.TrackTreeNode;
import TrackDto = trackEditor.TrackDto;

export const useTrackStore = defineStore("tacks", () => {
  const trackApi = useTracksApi();
  const trackTree = ref<TrackTreeNode>(new TrackTreeNode({ name: "", tracks: [], nodes: [] }));
  const selectedTrack = ref<TrackDto | undefined>(undefined);

  const selectedTrackId = computed(() => selectedTrack.value?.id);

  async function loadTracks() {
    trackTree.value = await trackApi.getTrackTree();
  }

  return {
    loadTracks,
    trackTree,
    selectedTrackId,
    selectedTrack,
  };
});
