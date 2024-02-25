<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { computed, onMounted, ref, watch } from "vue";
import { useJournalApi } from "../api/journal";
import ProgressSpinner from "primevue/progressspinner";
import Message from "primevue/message";
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputGroup from "primevue/inputgroup";
import TreeSelect from "primevue/treeselect";
import { backend } from "../../wailsjs/go/models";
import { useTracksApi } from "../api/tracks";
import { TreeNode } from "primevue/treenode";
import Calendar from "primevue/calendar";

const { t, d, locale } = useI18n();
const route = useRoute();
const journalApi = useJournalApi();
const router = useRouter();

const loading = ref(false);
const error = ref(false);
const selectedEntryId = computed(() => route.params["entryId"]);
const selectedEntry = ref<backend.JournalEntry | undefined>(undefined);
const selectedDate = ref<Date>(new Date());
const tracks = ref<backend.Track[]>([]);

const selectedTrack = ref<Record<string, boolean>>({});

const tracksApi = useTracksApi();

onMounted(async () => {
  tracks.value = await tracksApi.getTracks();
  // todo error handling
});

watch(selectedEntryId, () => loadEntry(), { immediate: true });

async function loadEntry() {
  selectedEntry.value = undefined;
  error.value = false;
  loading.value = true;
  if (!selectedEntryId.value || typeof selectedEntryId.value !== "string") {
    await router.replace("/journal");
    return;
  }
  try {
    selectedEntry.value = await journalApi.getListEntry(selectedEntryId.value);
    const trackKey = `${selectedEntry.value.track.baseId}-${selectedEntry.value.track.id}`;
    selectedTrack.value = { [trackKey]: true };
    selectedDate.value = new Date(Date.parse(selectedEntry.value.date));
  } catch (e) {
    console.error(e);
    error.value = true;
  } finally {
    loading.value = false;
  }
}

const length = computed(() => ((selectedEntry.value?.track.length || 0) / 1000).toFixed(1));

const trackSelection = computed<TreeNode[]>(() => {
  const names: Record<string, string> = {};
  const parentsGroups = tracks.value.reduce(
    (acc, curr) => {
      if (!acc[curr.baseId]) {
        names[curr.baseId] = curr.baseName;
        acc[curr.baseId] = [];
      }
      acc[curr.baseId].push(curr);
      return acc;
    },
    {} as Record<string, backend.Track[]>,
  );
  return Object.entries(parentsGroups).map(([key, group]) => ({
    key,
    selectable: false,
    label: names[key],
    children: group.map((track) => ({
      key: key + "-" + track.id,
      label: track.variant,
      data: track,
      selectedLabel: `${names[key]} – ${track.variant}`,
    })),
  }));
});

const pace = computed(() => {
  if (!selectedEntry.value || !/\d\d:\d\d:\d\d/.test(selectedEntry.value.time)) {
    return "";
  }
  const [hours, minutes, seconds] = selectedEntry.value.time.split(":").map((part) => Number(part));
  const secondsTotal = hours * 60 * 60 + minutes * 60 + seconds;
  let rawPace = secondsTotal / (selectedEntry.value.track.length / 1000);
  const paceHours = Math.floor(rawPace / (60 * 60));
  rawPace = rawPace % (60 * 60);
  const paceMinutes = Math.floor(rawPace / 60);
  console.log(paceMinutes, rawPace, selectedEntry.value.track.length)
  rawPace = rawPace % 60;
  console.log(rawPace)
  return `${paceHours.toString().padStart(2, "0")}:${paceMinutes.toString().padStart(2, "0")}:${rawPace.toFixed(0).toString().padStart(2, "0")}`;
});
</script>

<template>
  <div v-if="loading">
    <ProgressSpinner></ProgressSpinner>
  </div>
  <div v-else-if="error">
    <Message severity="error" :closable="false"
      ><div class="flex align-items-center">
        <span>{{ t("journal.loadEntryError") }}</span>
        <Button severity="danger" rounded text icon="pi pi-replay" @click="loadEntry"></Button></div
    ></Message>
  </div>
  <div v-else-if="selectedEntry" class="flex flex-column gap-2 w-full p-2">
    <InputGroup>
      <InputGroupAddon>
        <label for="date">{{ t("journal.details.date") }}</label>
      </InputGroupAddon>
      <Calendar
        id="date"
        v-model="selectedDate"
        :date-format="locale === 'de' ? 'dd.mm.yy' : 'yyyy/mm/dd'"
      ></Calendar>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="length">{{ t("journal.details.length") }}</label>
      </InputGroupAddon>
      <InputText id="length" :value="length" :disabled="true"></InputText>
      <InputGroupAddon>km</InputGroupAddon>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="time">{{ t("journal.details.time") }}</label>
      </InputGroupAddon>
      <InputText id="time" v-model="selectedEntry!.time"></InputText>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="pace">{{ t("journal.details.pace") }}</label>
      </InputGroupAddon>
      <InputText id="pace" :value="pace" disabled></InputText>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="track">{{ t("journal.details.track") }}</label>
      </InputGroupAddon>
      <TreeSelect
        id="track"
        v-model="selectedTrack"
        selection-mode="single"
        :options="trackSelection"
        placeholder="Select Item"
        class="md:w-20rem w-full"
        @node-select="(node) => (selectedEntry!.track = node.data)"
      >
        <template #value="props">
          {{ props.value[0].selectedLabel }}
        </template>
      </TreeSelect>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="comment">{{ t("journal.details.comment") }}</label>
      </InputGroupAddon>
      <InputText id="comment" v-model="selectedEntry!.comment"></InputText>
    </InputGroup>
  </div>
</template>

<style scoped></style>
