<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import type { TreeNode } from "primevue/treenode";
import { useRouter } from "vue-router";
import { TreeSelectionKeys } from "primevue/tree";
import { MenuItem } from "primevue/menuitem";
import { useI18n } from "vue-i18n";
import { tracksToTreeNodes } from "../shared/track-utils";
import CreateTrackOverlay from "./CreateTrackOverlay.vue";
import Button from "primevue/button";
import { useTracksApi } from "../api/tracks";
import { tracks } from "../../wailsjs/go/models";
import OverlayPanel from "primevue/overlaypanel";
import CreateEntryOverlay from "../journal/CreateEntryOverlay.vue";
import Track = tracks.Track;

const trackStore = useTrackStore();
const { availableTracks, selectedTrackId } = storeToRefs(trackStore);
const { t } = useI18n();
const tracksApi = useTracksApi();

const selectableTracks = computed(() => tracksToTreeNodes(availableTracks.value, true));

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
const moreMenuOpenedOn = ref<{ track: Track; event: any } | undefined>(undefined);
const addClickedOn = ref<{ parentId: string; target: HTMLElement } | undefined>(undefined);
const deleteConfirm = ref();
const deleteConfirmMessage = ref("");
const menuItems = ref<MenuItem>([
  {
    label: t("shared.delete"),
    icon: "pi pi-trash",
    command: async (event: Event) => {
      if (!moreMenuOpenedOn.value) {
        return;
      }
      const countChildren: (acc: number, track: Track) => number = (acc: number, track: Track) => {
        return 0;
      };
      const children = countChildren(-1, moreMenuOpenedOn.value.track);
      deleteConfirmMessage.value = t("tracks.deleteConfirmation", {
        children,
        count: moreMenuOpenedOn.value.track.usages,
      });
      setTimeout(() =>
        deleteConfirm.value.show(new Event("click"), moreMenuOpenedOn.value!.event.target),
      );
    },
  },
]);

async function deleteTrack() {
  if (!moreMenuOpenedOn.value) {
    return;
  }
  try {
    await tracksApi.deleteTrack(moreMenuOpenedOn.value.track.id);
    trackStore.deleteTrack(moreMenuOpenedOn.value.track.id);
    deleteConfirm.value.hide();
    moreMenuOpenedOn.value = undefined;
  } catch (e) {
    console.error(e);
  }
}

function openMoreMenu(event: Event, track: Track, isFolder: boolean) {
  if (moreMenuOpenedOn.value && moreMenuOpenedOn.value.track !== track) {
    treeNodeMenu.value.hide();
    moreMenuOpenedOn.value = undefined;
    return;
  }
  moreMenuOpenedOn.value = { track, event };
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
    :pt="{ label: { class: 'w-full flex align-items-center overflow-hidden' } }"
  >
    <template #default="slotProps">
      <span
        class="flex-grow-1 flex-shrink-1 white-space-nowrap text-overflow-ellipsis overflow-hidden"
        >{{ slotProps.node.label }}</span
      >
      <Button
        class="flex-shrink-0"
        text
        rounded
        icon="pi pi-ellipsis-v"
        @click.stop.prevent="
          openMoreMenu($event, slotProps.node.data, !!slotProps.node.children?.length)
        "
      ></Button>
    </template>
  </Tree>
  <Menu ref="treeNodeMenu" :model="menuItems" :popup="true"></Menu>
  <OverlayPanel ref="deleteConfirm">
    <div class="flex align-items-center gap-2">
      {{ deleteConfirmMessage }}
      <Button @click="deleteTrack()">{{ t("shared.delete") }}</Button>
    </div>
  </OverlayPanel>
</template>

<style scoped></style>
