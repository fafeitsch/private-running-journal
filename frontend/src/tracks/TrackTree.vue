<script setup lang="ts">
import { computed, ref, Ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import type { TreeNode } from "primevue/treenode";
import { useRouter } from "vue-router";
import { TreeSelectionKeys } from "primevue/tree";

const trackStore = useTrackStore();
const { selectedTrackId } = storeToRefs(trackStore);
const { availableTracks }: { availableTracks: Ref<TreeNode[]> } = storeToRefs(trackStore);

const selectableTracks = computed(() => {
  const makeSelectable: (node: TreeNode) => TreeNode = (node: TreeNode) => ({
    ...node,
    selectable: node.data.length > 0,
    children: node.children?.map(makeSelectable),
  });
  return availableTracks.value.map(makeSelectable);
});

const selection = ref<TreeSelectionKeys>({});

watch(
  selectedTrackId,
  (value) => {
    selection.value = value ? { [value]: true } : {};
  },
  { immediate: true },
);

const expansion = ref<TreeSelectionKeys>({});

watch(
  () => ({ trackId: selectedTrackId.value, tracks: selectableTracks.value }),
  ({ tracks, trackId }) => {
    expansion.value = {};
    if (!trackId || !tracks) {
      return;
    }
    const setParent: (acc: Record<string, string>, node: TreeNode) => Record<string, string> = (
      acc: Record<string, string>,
      node: TreeNode,
    ) => {
      node.children?.forEach((child) => (acc[child.key as string] = node.key as string));
      return node.children ? node.children.reduce(setParent, acc) : acc;
    };
    const parents = tracks.reduce(setParent, {} as Record<string, string>);
    let parent = parents[trackId];
    while (parent) {
      expansion.value[parent as string] = true;
      parent = parents[parent];
    }
  },
);

const router = useRouter();

function selectNode(node: TreeNode) {
  router.push(`/tracks/${encodeURIComponent(node.key!)}`);
}
</script>

<template>
  <Tree
    class="h-full overflow-auto"
    :value="selectableTracks"
    v-model:selection-keys="selection"
    v-model:expanded-keys="expansion"
    selection-mode="single"
    @node-select="selectNode"
  ></Tree>
</template>

<style scoped></style>
