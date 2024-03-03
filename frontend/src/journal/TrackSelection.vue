<script setup lang="ts">
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import TreeSelect from "primevue/treeselect";
import { useI18n } from "vue-i18n";
import { onMounted, ref, toRefs, watch } from "vue";
import { tracks } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import { TreeNode } from "primevue/treenode";
import Track = tracks.Track;

const { t } = useI18n();
const tracksApi = useTracksApi();
const selectedTrack = defineModel<Track | undefined>();

const props = defineProps<{ linkedTrack?: string }>();
const { linkedTrack } = toRefs(props);

const availableTracks = ref<TreeNode[]>([]);

onMounted(async () => {
  try {
    const trackToListEntry: (tracks: tracks.Track, parentNames: string) => TreeNode = (
      track: tracks.Track,
      parentNames: string,
    ) => {
      const name = parentNames ? `${parentNames} / ${track.name}` : track.name;
      return {
        key: track.id,
        label: track.name,
        data: track,
        children: track.variants.map((entry) => trackToListEntry(entry, name)),
        selectable: track.length > 0,
        selectedLabel: name,
      };
    };
    const rawTracks = await tracksApi.getTracks();
    availableTracks.value = rawTracks.map((entry) => trackToListEntry(entry, ""));
    console.log("available tracks", availableTracks.value);
  } catch (e) {
    // todo error handling
    console.error(e);
  }
});

const selectedEntry = ref<Record<string, boolean>>({});

watch(
  selectedTrack,
  () => {
    if (!selectedTrack.value) {
      selectedEntry.value = {};
      return;
    }
    selectedEntry.value = {
      [selectedTrack.value.id]: true,
    };
  },
  { immediate: true },
);
</script>

<template>
  <InputGroup>
    <InputGroupAddon class="flex gap-2">
      <label for="track">{{ t("journal.details.track") }}</label>
      <span
        v-if="!selectedTrack && linkedTrack"
        class="text-red-500 pi pi-exclamation-triangle"
        v-tooltip="{
          value: t('journal.details.trackNotFound', { link: linkedTrack }),
          showDelay: 500,
        }"
      ></span>
    </InputGroupAddon>
    <TreeSelect
      id="track"
      v-model="selectedEntry"
      selection-mode="single"
      :options="availableTracks"
      placeholder="Select Item"
      class="md:w-20rem w-full"
      @node-select="(node) => (selectedTrack = node.data)"
    >
      <template #value="props">
        {{ props.value[0]?.selectedLabel }}
      </template>
    </TreeSelect>
  </InputGroup>
</template>

<style scoped></style>
