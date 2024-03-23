<script setup lang="ts">
import { computed, ref, Ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import type { TreeNode } from "primevue/treenode";
import { useRouter } from "vue-router";
import { TreeSelectionKeys } from "primevue/tree";
import CreateTrackOverlay from "./CreateTrackOverlay.vue";
import { MenuItem } from "primevue/menuitem";
import { useI18n } from "vue-i18n";
import { useTracksApi } from "../api/tracks";
import { shared } from "../../wailsjs/go/models";
import Track = shared.Track;
import {tracksToTreeNodes} from '../shared/track-utils';

const tracksApi = useTracksApi();
const trackStore = useTrackStore();
const { availableTracks, selectedTrackId } = storeToRefs(trackStore);
const { t } = useI18n();

const selectableTracks = computed(() => tracksToTreeNodes(availableTracks.value));

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

const treeNodeMenu = ref();
const clickedTrack = ref<Track | undefined>(undefined);

const menuItems = ref<MenuItem>([
  {
    label: t("shared.add"),
    icon: "pi pi-plus",
    command: async () => {
      if (!clickedTrack.value) {
        return;
      }
      const name = `${clickedTrack.value.name} ${t("tracks.variant")}`;
      try {
        const track = await tracksApi.createTrack({ parent: clickedTrack.value.id, name });
        trackStore.addTrack(track, clickedTrack.value.id);
        router.push("/tracks/" + encodeURIComponent(track.id));
      } catch (error) {
        //TODO error handling
        console.error(error);
      }
    },
  },
]);

function showContextMenu(track: Track, event: any) {
  clickedTrack.value = track;
  treeNodeMenu.value.show(event);
}
</script>

<template>
  <header class="flex justify-content-between align-items-center">
    <span class="text-2xl">{{ $t("tracks.title") }}</span
    ><CreateTrackOverlay></CreateTrackOverlay>
  </header>
  <Tree
    class="h-full overflow-auto"
    :value="selectableTracks"
    v-model:selection-keys="selection"
    v-model:expanded-keys="expansion"
    selection-mode="single"
    @node-select="selectNode"
  >
    <template #default="slotProps">
      <span @contextmenu="showContextMenu(slotProps.node.data, $event)">{{
        slotProps.node.label
      }}</span>
    </template>
  </Tree>
  <ContextMenu ref="treeNodeMenu" :model="menuItems"></ContextMenu>
</template>

<style scoped></style>
