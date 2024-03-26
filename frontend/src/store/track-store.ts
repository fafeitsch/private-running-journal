import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { TreeNode } from "primevue/treenode";
import { useTracksApi } from "../api/tracks";
import { tracks } from "../../wailsjs/go/models";
import Track = tracks.Track;
import SaveTrack = tracks.SaveTrack;

export const useTrackStore = defineStore("tacks", () => {
  const trackApi = useTracksApi();
  const availableTracks = ref<Track[]>([]);
  const selectedTrackId = ref<string | undefined>(undefined);

  async function loadTracks() {
    availableTracks.value = await trackApi.getTracks();
  }

  function findTrack(id: string) {
    const findTrack: (acc: tracks.Track | undefined, node: Track) => tracks.Track | undefined = (
      acc: tracks.Track | undefined,
      node: Track,
    ) => {
      if (acc) {
        return acc;
      }
      if (node.id === id) {
        return node;
      }
      return node.variants.reduce(findTrack, undefined);
    };
    return findTrack;
  }

  const selectedTrack = computed(() => {
    if (!selectedTrackId.value) {
      return undefined;
    }

    return availableTracks.value.reduce(findTrack(selectedTrackId.value), undefined);
  });

  function addTrack(track: Track, parentKey?: string) {
    if (!parentKey) {
      availableTracks.value.push(track);
      availableTracks.value.sort((a, b) => a.name.localeCompare(b.name));
      return;
    }
    const parent = availableTracks.value.reduce(findTrack(parentKey), undefined);
    if (!parent) {
      return;
    }
    parent.variants.push(track);
    parent.variants.sort((a, b) => a.name.localeCompare(b.name));
  }

  function updateTrack(track: Track) {
    const existing = availableTracks.value.reduce(findTrack(track.id!), undefined)
    if(!existing) {
      return
    }
    existing.name = track.name
    existing.length = track.length
    existing.parentNames = track.parentNames
    existing.variants = track.variants
    existing.usages = track.usages
  }

  return { loadTracks, availableTracks, selectedTrackId, selectedTrack, addTrack, updateTrack };
});
