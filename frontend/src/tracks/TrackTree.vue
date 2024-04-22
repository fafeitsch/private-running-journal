<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import type { TreeNode } from "primevue/treenode";
import { useRouter } from "vue-router";
import { TreeSelectionKeys } from "primevue/tree";
import { useI18n } from "vue-i18n";
import { tracksToTreeNodes } from "../shared/track-utils";
import CreateTrackOverlay from "./CreateTrackOverlay.vue";

const trackStore = useTrackStore();
const { availableTracks, selectedTrackId } = storeToRefs(trackStore);
const { t } = useI18n();

const selectableTracks = computed(() => tracksToTreeNodes(availableTracks.value, true));

const selection = ref<TreeSelectionKeys>({});

watch(
  selectedTrackId,
  (value) => {
    selection.value = value ? { [value]: true } : selection.value;
  },
  { immediate: true },
);

const expansion = ref<TreeSelectionKeys>({ root: true });

watch(
  () => ({ trackId: selectedTrackId.value, tracks: selectableTracks.value }),
  ({ tracks, trackId }) => {
    expansion.value = { root: true };
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
  <div class="flex h-full flex-column overflow-hidden gap-2">
    <header class="flex justify-content-between align-items-center">
      <span class="text-2xl">{{ $t("tracks.title") }}</span
      ><CreateTrackOverlay></CreateTrackOverlay>
    </header>
    <Tree
      class="flex-grow-1 flex-shrink-1 overflow-auto"
      :value="selectableTracks"
      v-model:selection-keys="selection"
      v-model:expanded-keys="expansion"
      selection-mode="single"
      @node-select="selectNode"
      :pt="{
        label: { class: 'w-full flex align-items-center white-space-nowrap', 'data-testid': 'track-tree-node' },
        toggler: {'data-testid': 'track-tree-node-toggler'}
      }"
    >
    </Tree>
  </div>
</template>

<style scoped></style>
