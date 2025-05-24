import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useTracksApi } from "../api/tracks";
import {projection, trackEditor} from "../../wailsjs/go/models";
import TrackDto = trackEditor.TrackDto;
import TrackTreeNode = projection.TrackTreeNode;

export const useTrackStore = defineStore("tacks", () => {
  const trackApi = useTracksApi();
  const trackTree = ref<TrackTreeNode>(new TrackTreeNode({ name: "", tracks: [], nodes: [] }));
  const selectedTrack = ref<TrackDto | undefined>(undefined);

  const selectedTrackId = computed(() => selectedTrack.value?.id);

  async function loadTracks(filter: string = "") {
    trackTree.value = await trackApi.getTrackTree(filter);
  }

  return {
    loadTracks,
    trackTree,
    selectedTrackId,
    selectedTrack,
  };
});
