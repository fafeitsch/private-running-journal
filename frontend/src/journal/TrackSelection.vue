<script setup lang="ts">
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import TreeSelect from "primevue/treeselect";
import { useI18n } from "vue-i18n";
import { onMounted, ref, toRefs, watch, watchEffect } from "vue";
import { projection } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import type { TreeNode } from "primevue/treenode";
import { tracksToTreeNodes } from "../shared/track-utils";
import TrackTreeEntry = projection.TrackTreeEntry;

const { t } = useI18n();
const tracksApi = useTracksApi();
const selectedTrack = defineModel<TrackTreeEntry | undefined>();

const props = defineProps<{ linkedTrack?: string }>();
const { linkedTrack } = toRefs(props);

const availableTracks = ref<TreeNode[]>([]);
onMounted(async () => {
  try {
    const tracks = await tracksApi.getTrackTree();
    availableTracks.value = tracksToTreeNodes(tracks);
  } catch (e) {
    // todo error handling
    console.error(e);
  }
});

const selectedEntry = ref<Record<string, boolean>>({});

watchEffect(() => {
  const reducer = (acc: TrackTreeEntry | undefined, cur: TreeNode): TrackTreeEntry | undefined => {
    if (acc) {
      return acc;
    }
    if(cur.key === linkedTrack.value) {
      return cur.data
    }
    const track = cur.children?.find((c) => c.key === linkedTrack.value);
    if (track) {
      return track.data;
    }
    return cur.children?.reduce(reducer, undefined);
  };
  selectedTrack.value = availableTracks.value.reduce(reducer, undefined);
});

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
      <label for="track-select">{{ t("journal.details.track") }}*</label>
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
      id="track-select"
      v-model="selectedEntry"
      selection-mode="single"
      :options="availableTracks"
      class="md:w-80 w-full"
      @node-select="(node) => (selectedTrack = node.data)"
      data-testid="track-selection"
      :pt="{
        pcTree: {
          nodeLabel: {
            class: 'w-full flex items-center whitespace-nowrap',
            'data-testid': 'track-tree-selection-node',
          },
          nodeToggleButton: { 'data-testid': 'track-tree-selection-node-toggler' },
        },
      }"
    >
      <template #value="props">
        {{ props.value[0]?.label }}
      </template>
    </TreeSelect>
  </InputGroup>
</template>

<style scoped></style>
