<script setup lang="ts">
import Popover from "primevue/popover";
import Button from "primevue/button";
import { computed,  ref, watch } from "vue";
import Message from "primevue/message";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { useTracksApi } from "../api/tracks";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { useTrackStore } from "../store/track-store";
import TreeSelect from "primevue/treeselect";
import { storeToRefs } from "pinia";
import { TreeNode } from "primevue/treenode";
import { tracksToTreeNodes } from "../shared/track-utils";

const {  t } = useI18n();

const overlayPanel = ref();
const error = ref<boolean>(false);

const tracksApi = useTracksApi();
const tracksStore = useTrackStore();
const { trackTree, selectedTrack } = storeToRefs(tracksStore);
const folders = computed(() => {
  const filterFolders = (folder: TreeNode) => !!folder.children?.length;
  const recursiveMap: (folder: TreeNode) => TreeNode = (folder: TreeNode) => ({
    ...folder,
    selectable: true,
    children: folder.children!.filter(filterFolders).map(recursiveMap),
  });
  const tracks = tracksToTreeNodes(trackTree.value).filter(filterFolders).map(recursiveMap);
  tracks.push({
    selectable: true,
    label: t("tracks.createFolder"),
    children: [],
    key: "//new-folder",
  });
  return tracks;
});
const folderSelection = ref<{}>({});
const selectedFolder = computed(() => Object.keys(folderSelection.value)[0]);

watch(selectedFolder, (value) => {
  if (selectedFolder.value === "//new-folder") {
    return "";
  }
  return (folderName.value = value || "");
});

const folderName = ref("");
const forbiddenFolderName = computed(() => folderName.value.includes(".."));

async function moveTrack() {
  if (folderName.value === undefined || !selectedTrack.value) {
    return;
  }
  error.value = false;

  try {
    let newPath = folderName.value.split("/").filter(s => !!s)
    await tracksApi.saveTrack({...selectedTrack.value, parents: newPath});
    await tracksStore.loadTracks()
    overlayPanel.value.hide();
  } catch (e) {
    error.value = true;
    console.error(e);
  }
}
</script>

<template>
  <Button icon="pi pi-arrow-right" @click="(event) => overlayPanel.toggle(event)"></Button>
  <Popover ref="overlayPanel">
    <div v-focustrap class="flex flex-col gap-2 overlay">
      <InputGroup class="flex w-full">
        <InputGroupAddon>
          <label for="parentInput" class="px-2">{{ t("tracks.folder") }}</label>
        </InputGroupAddon>
        <TreeSelect
          id="parentInput"
          v-model="folderSelection"
          selection-mode="single"
          :options="folders"
          class="w-full"
        >
        </TreeSelect>
      </InputGroup>
      <InputGroup v-if="selectedFolder === '//new-folder'" class="flex w-full">
        <InputGroupAddon>
          <label for="newFolderName">{{ t("tracks.folderName") }}</label>
        </InputGroupAddon>
        <InputText
          class="grow"
          id="newFolderName"
          v-model="folderName"
          autofocus
        ></InputText>
      </InputGroup>
      <div class="flex gap-2">
        <Message v-if="error" class="grow shrink" severity="error">{{
          t("tracks.moveFailed")
        }}</Message>
        <span v-else class="grow"></span>
        <Button
          :label="t('tracks.move')"
          @click="moveTrack"
          :disabled="forbiddenFolderName"
        ></Button>
      </div>
    </div>
  </Popover>
</template>

<style scoped>
.overlay {
  width: 400px;
}
</style>
