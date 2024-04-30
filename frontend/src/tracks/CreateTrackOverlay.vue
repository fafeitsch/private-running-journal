<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel";
import InlineMessage from "primevue/inlinemessage";
import Button from "primevue/button";
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useJournalStore } from "../store/journal-store";
import { useRouter } from "vue-router";
import { useTracksApi } from "../api/tracks";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { useTrackStore } from "../store/track-store";
import TreeSelect from "primevue/treeselect";
import { storeToRefs } from "pinia";
import { TreeNode } from "primevue/treenode";
import { tracksToTreeNodes } from "../shared/track-utils";

const { locale, t } = useI18n();

const overlayPanel = ref();
const name = ref<string>("");
const error = ref<boolean>(false);

const tracksApi = useTracksApi();
const tracksStore = useTrackStore();
const store = useJournalStore();
const { availableTracks, selectedTrack } = storeToRefs(tracksStore);
const folders = computed(() => {
  const filterFolders = (folder: TreeNode) => !!folder.children?.length;
  const recursiveMap: (folder: TreeNode) => TreeNode = (folder: TreeNode) => ({
    ...folder,
    selectable: true,
    children: folder.children!.filter(filterFolders).map(recursiveMap),
  });
  const tracks = tracksToTreeNodes(availableTracks.value).filter(filterFolders).map(recursiveMap);
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

watch(selectedTrack, (value) => {
  if (!value) {
    return;
  }
  folderSelection.value = { ["/" + value.hierarchy.join("/")]: true };
});

watch(selectedFolder, (value) => {
  if (selectedFolder.value === "//new-folder") {
    return "";
  }
  return (folderName.value = value || "");
});

const folderName = ref("");
const forbiddenFolderName = computed(() => folderName.value.includes(".."));

const router = useRouter();

async function createEntry() {
  if (!name.value) {
    return;
  }
  error.value = false;

  try {
    const track = await tracksApi.createTrack({
      name: name.value,
      parent: folderName.value.startsWith("/") ? folderName.value.substring(1) : folderName.value,
    });
    tracksStore.addTrack(track);
    router.push("/tracks/" + encodeURIComponent(track.id));
    overlayPanel.value.hide();
  } catch (e) {
    error.value = true;
    console.error(e);
  }
}
</script>

<template>
  <Button
    icon="pi pi-plus"
    @click="(event) => overlayPanel.toggle(event)"
    :aria-label="t('shared.add')"
    :v-tooltip="t('shared.add')"
  ></Button>
  <OverlayPanel ref="overlayPanel">
    <div v-focustrap class="flex flex-column gap-2 overlay">
      <InputGroup class="flex w-full">
        <InputGroupAddon>
          <label for="newTrackName">{{ t("tracks.name") }}</label>
        </InputGroupAddon>
        <InputText class="flex-grow-1" id="newTrackName" v-model="name" autofocus></InputText>
      </InputGroup>
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
          class="flex-grow-1"
          id="newFolderName"
          v-model="folderName"
          autofocus
        ></InputText>
      </InputGroup>
      <div class="flex gap-2">
        <InlineMessage v-if="error" class="flex-grow-1 flex-shrink-1" severity="error">{{
          t("journal.createEntryError")
        }}</InlineMessage>
        <span v-else class="flex-grow-1"></span>
        <Button
          :label="t('shared.add')"
          @click="createEntry"
          :disabled="!name || forbiddenFolderName"
          data-testid="create-empty-track-button"
        ></Button>
      </div>
    </div>
  </OverlayPanel>
</template>

<style scoped>
.overlay {
  width: 400px;
}
</style>
