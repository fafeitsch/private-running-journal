<script setup lang="ts">
import Popover from "primevue/popover";
import Message from "primevue/message";
import Button from "primevue/button";
import { computed, ref, toRefs, watch } from "vue";
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
import { tracks } from "../../wailsjs/go/models";
import CreateTrack = tracks.CreateTrack;
import Coordinates = tracks.Coordinates;

const {  t } = useI18n();

const props = defineProps<{ name: string; waypoints: Coordinates[] }>();

const { name, waypoints } = toRefs(props);

const overlayPanel = ref();
const error = ref<boolean>(false);

const tracksApi = useTracksApi();
const tracksStore = useTrackStore();
const { trackTree } = storeToRefs(tracksStore);
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

const router = useRouter();

const emit = defineEmits<{ trackCreated: [] }>();

async function createEntry() {
  if (!name.value) {
    return;
  }
  error.value = false;

  try {
    const track = await tracksApi.createTrack(
      new CreateTrack({
        name: name.value,
        parent: folderName.value.startsWith("/") ? folderName.value.substring(1) : folderName.value,
        waypoints: waypoints.value,
      }),
    );
    emit("trackCreated");
    await tracksStore.loadTracks()
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
    icon="pi pi-save"
    @click="(event) => overlayPanel.toggle(event)"
    :aria-label="t('shared.save')"
    :disabled="!name || name === 'new'"
    v-tooltip="{ value: t('shared.save'), showDelay: 500 }"
  ></Button>
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
        <InputText class="grow" id="newFolderName" v-model="folderName" autofocus></InputText>
      </InputGroup>
      <div class="flex gap-2">
        <Message v-if="error" class="grow shrink" severity="error"
          >{{ t("journal.createEntryError") }}
        </Message>
        <span v-else class="grow"></span>
        <Button
          :label="t('shared.add')"
          @click="createEntry"
          :disabled="!name || forbiddenFolderName"
          data-testid="create-empty-track-button"
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
