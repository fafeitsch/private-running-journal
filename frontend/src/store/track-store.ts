import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { TreeNode } from "primevue/treenode";
import { useTracksApi } from "../api/tracks";
import { tracks } from "../../wailsjs/go/models";

export const useTrackStore = defineStore("tacks", () => {
  const trackApi = useTracksApi();
  const availableTracks = ref<TreeNode[]>([]);
  const selectedTrackId = ref<string | undefined>(undefined);

  async function loadTracks() {
    availableTracks.value = await trackApi.getTracks();
  }

  const selectedTrack = computed(() => {
    if (!selectedTrackId.value) {
      return undefined;
    }
    const findTrack: (acc: tracks.Track | undefined, node: TreeNode) => tracks.Track | undefined = (
      acc: tracks.Track | undefined,
      node: TreeNode,
    ) => {
      if (acc) {
        return acc;
      }
      if (node.key === selectedTrackId.value) {
        return node.data;
      }
      return node.children?.reduce(findTrack, undefined);
    };
    return availableTracks.value.reduce(findTrack, undefined);
  });

  return { loadTracks, availableTracks, selectedTrackId, selectedTrack };
});
