<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel";
import InlineMessage from "primevue/inlinemessage";
import Button from "primevue/button";
import { ref, toRefs, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useJournalStore } from "../store/journal-store";
import { useRouter } from "vue-router";
import { useTracksApi } from "../api/tracks";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import { useTrackStore } from "../store/track-store";

const { locale, t } = useI18n();

const props = defineProps<{
  showEvent: { parentId: string; target: HTMLElement } | undefined;
}>();
const { showEvent } = toRefs(props);

watch(showEvent, (value) => {
  if (value) {
    setTimeout(() => overlayPanel.value.show(new Event("click"), value.target));
  } else {
    overlayPanel.value.hide();
  }
});

const overlayPanel = ref();
const name = ref<string>("");
const error = ref<boolean>(false);

const tracksApi = useTracksApi();
const tracksStore = useTrackStore();
const store = useJournalStore();
const router = useRouter();

async function createEntry() {
  if (!name.value || !showEvent.value) {
    return;
  }
  error.value = false;

  try {
    let parentId = showEvent.value.parentId === "root" ? "" : showEvent.value.parentId;
    const track = await tracksApi.createTrack({
      name: name.value,
      parent: parentId,
    });
    tracksStore.addTrack(track, parentId);
    router.push("/tracks/" + encodeURIComponent(track.id));
    overlayPanel.value.hide();
  } catch (e) {
    error.value = true;
    console.error(e);
  }
}
</script>

<template>
  <OverlayPanel ref="overlayPanel">
    <div v-focustrap class="flex flex-column gap-2 overlay">
      <InputGroup class="flex w-full">
        <InputGroupAddon>
          <label for="newTrackName">{{ t("tracks.name") }}</label>
        </InputGroupAddon>
        <InputText class="flex-grow-1" id="newTrackName" v-model="name" autofocus></InputText>
      </InputGroup>
      <div class="flex gap-2">
        <InlineMessage v-if="error" class="flex-grow-1 flex-shrink-1" severity="error">{{
          t("journal.createEntryError")
        }}</InlineMessage>
        <span v-else class="flex-grow-1"></span>
        <Button :label="t('journal.createEntry')" @click="createEntry" :disabled="!name"></Button>
      </div>
    </div>
  </OverlayPanel>
</template>

<style scoped>
.overlay {
  width: 400px;
}
</style>
