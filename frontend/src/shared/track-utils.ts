import {TreeNode} from 'primevue/treenode';
import { tracks} from '../../wailsjs/go/models';
import Track =tracks.Track;

export function tracksToTreeNodes(tracks: Track[], allSelectable = false): TreeNode[] {
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
      selectable: allSelectable || track.length > 0,
      selectedLabel: name,
      icon: track.length > 0 ? "pi pi-directions" : undefined,
    };
  };
  return tracks
    .map((entry) => trackToListEntry(entry, ""))
    .sort((t1, t2) => t1.label!.localeCompare(t2.label!));
}
