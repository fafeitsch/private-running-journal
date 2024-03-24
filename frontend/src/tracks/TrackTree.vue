<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import type { TreeNode } from "primevue/treenode";
import { useRouter } from "vue-router";
import { TreeSelectionKeys } from "primevue/tree";
import { MenuItem } from "primevue/menuitem";
import { useI18n } from "vue-i18n";
import { shared } from "../../wailsjs/go/models";
import { tracksToTreeNodes } from "../shared/track-utils";
import Track = shared.Track;
import CreateTrackOverlay from "./CreateTrackOverlay.vue";

const trackStore = useTrackStore();
const { availableTracks, selectedTrackId } = storeToRefs(trackStore);
const { t } = useI18n();

const selectableTracks = computed(() => [
  {
    label: t("tracks.title"),
    selectable: false,
    key: "root",
    type: "root",
    expandedIcon: "",
    children: tracksToTreeNodes(availableTracks.value, true),
  },
]);

const selection = ref<TreeSelectionKeys>({});

watch(
  selectedTrackId,
  (value) => {
    selection.value = value ? { [value]: true } : {};
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

const treeNodeMenu = ref();
const contextMenuOpenedOn = ref<{ track: Track | "root"; event: any } | undefined>(undefined);
const addClickedOn = ref<{ parentId: string; target: HTMLElement } | undefined>(undefined);
const menuItems = ref<MenuItem>([
  {
    label: t("shared.add"),
    icon: "pi pi-plus",
    command: async (event: Event) => {
      if (!contextMenuOpenedOn.value) {
        return;
      }
      let clickedTrack = contextMenuOpenedOn.value.track;
      addClickedOn.value = {
        target: contextMenuOpenedOn.value.event.target,
        parentId: clickedTrack === "root" ? clickedTrack : clickedTrack.id,
      };
    },
  },
]);

function showContextMenu(track: Track | "root", event: any) {
  contextMenuOpenedOn.value = { track, event };
  treeNodeMenu.value.show(event);
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
    :pt="{ label: { class: 'w-full flex' } }"
  >
    <template #default="slotProps">
      <span @contextmenu="showContextMenu(slotProps.node.data, $event)" class="w-full">{{
        slotProps.node.label
      }}</span>
    </template>
    <template #root="{ node }">
      <span @contextmenu="showContextMenu('root', $event)" class="text-2xl w-full text-color">{{
        node.label
      }}</span>
    </template>
  </Tree>
  <ContextMenu ref="treeNodeMenu" :model="menuItems"></ContextMenu>
  <CreateTrackOverlay :show-event="addClickedOn"></CreateTrackOverlay>
</template>

<style scoped></style>
