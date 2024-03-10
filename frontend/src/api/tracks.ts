import { tracks } from "../../wailsjs/go/models";
import { GetGpxData, GetTracks } from "../../wailsjs/go/backend/App";
import { TreeNode } from "primevue/treenode";

export function useTracksApi() {
  async function getTracks(): Promise<TreeNode[]> {
    const trackToListEntry: (tracks: tracks.Track, parentNames: string) => TreeNode = (
      track: tracks.Track,
      parentNames: string,
    ) => {
      const name = parentNames ? `${parentNames} / ${track.name}` : track.name;
      return {
        key: track.id,
        label: track.name,
        data: track,
        children: track.variants.map((entry) => trackToListEntry(entry, name)),
        selectable: track.length > 0,
        selectedLabel: name,
        icon: track.length > 0 ? "pi pi-directions" : undefined,
      };
    };
    const rawTracks = await GetTracks();
    return rawTracks
      .map((entry) => trackToListEntry(entry, ""))
      .sort((t1, t2) => t1.label!.localeCompare(t2.label!));
  }
  function getGpxData(baseName: string): Promise<tracks.GpxData> {
    return GetGpxData(baseName);
  }
  return { getTracks, getGpxData };
}
