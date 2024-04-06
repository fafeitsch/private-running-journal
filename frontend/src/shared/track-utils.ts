import { TreeNode } from "primevue/treenode";
import { tracks } from "../../wailsjs/go/models";
import Track = tracks.Track;

export function tracksToTreeNodes(tracks: Track[], allSelectable = false): TreeNode[] {
  const roots: TreeNode[] = [];
  const folders: Record<string, TreeNode> = {};
  tracks.forEach((track) => {
    const hierarchy = [...track.hierarchy];
    while (hierarchy.length) {
      const path = "/" + hierarchy.join("/");
      const name = hierarchy.pop();
      folders[path] = folders[path] || {
        key: path,
        label: name,
        selectable: false,
        children: [],
      };
    }
    const trackNode = {
      key: track.id,
      label: track.name,
      selectable: true,
      data: track,
    };
    if(track.hierarchy.length) {
      folders["/" + track.hierarchy.join("/")].children!.push(trackNode);
    }else {
      roots.push(trackNode)
    }
  });
  Object.values(folders).forEach((folder) => {
    const parentId = folder.key!.substring(0, folder.key!.lastIndexOf("/"));
    if (parentId === "") {
      roots.push(folder);
    } else {
      folders[parentId].children!.push(folder);
    }
  });
  const sortFolders = (folder1: TreeNode, folder2: TreeNode) => {
    if ((folder1.children && folder2.children) || (!folder1.children && !folder2.children)) {
      return folder1.label!.localeCompare(folder2.label!);
    }
    return folder1.children ? -1 : 1;
  };
  Object.values(folders).forEach((folder) => folder.children!.sort(sortFolders));
  roots.sort(sortFolders);
  return roots;
}
