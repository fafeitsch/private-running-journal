import { TreeNode } from "primevue/treenode";
import {projection} from '../../wailsjs/go/models';
import TrackTreeNode = projection.TrackTreeNode;

export function tracksToTreeNodes(root: TrackTreeNode, allSelectable = false, hierarchy = "/"): TreeNode[] {
  const roots: TreeNode[] = [];

  root.nodes.forEach((node) => {
    const key = [hierarchy, node.name].join("/");
    roots.push({
      key,
      label: node.name,
      selectable: false,
      children: tracksToTreeNodes(node, allSelectable, key),
      data: key
    });
  });

  root.tracks.forEach((track) => {
    roots.push({ key: track.id, label: track.name, data: track, selectable: true });
  });

  const sortFolders = (folder1: TreeNode, folder2: TreeNode) => {
    if ((folder1.children && folder2.children) || (!folder1.children && !folder2.children)) {
      return folder1.label!.localeCompare(folder2.label!);
    }
    return folder1.children ? -1 : 1;
  };
  roots.sort(sortFolders);
  return roots;
}
