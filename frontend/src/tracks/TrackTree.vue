<script setup lang="ts">
import { computed, onMounted, Ref } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import { TreeNode } from "primevue/treenode";
import { useRouter } from "vue-router";

const trackStore = useTrackStore();
const { availableTracks }: { availableTracks: Ref<TreeNode[]> } = storeToRefs(trackStore);

onMounted(async () => {
  try {
    await trackStore.loadTracks();
  } catch (e) {
    // TODO error handling
    console.error(e);
  }
});

const selectableTracks = computed(() => {
  const makeSelectable: (node: TreeNode) => TreeNode = (node: TreeNode) => ({
    ...node,
    selectable: true,
    children: node.children?.map(makeSelectable),
  });
  return availableTracks.value.map(makeSelectable);
});

const router = useRouter();

function selectNode(node: TreeNode) {
  router.push(`/tracks/${encodeURIComponent(node.key!)}`);
}
</script>

<template>
  <Tree
    class="h-full overflow-auto"
    :value="selectableTracks"
    selection-mode="single"
    @node-select="selectNode"
  ></Tree>
</template>

<style scoped></style>
