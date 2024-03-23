<script setup lang="ts">
import OverlayPanel from "primevue/overlaypanel";
import InlineMessage from "primevue/inlinemessage";
import Button from "primevue/button";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useJournalStore } from "../store/journal-store";
import { useRouter } from "vue-router";
import { useTracksApi } from "../api/tracks";
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';

const { locale, t } = useI18n();

const overlayPanel = ref();
const name = ref<string>("");
const error = ref<boolean>(false);

const tracksApi = useTracksApi();
const store = useJournalStore();
const router = useRouter();

async function createEntry() {
  if (!name.value) {
    return;
  }
  error.value = false;

  try {
    const track = await tracksApi.createTrack({ name: name.value, parent: "" });
    router.push("/tracks/" + encodeURIComponent(track.id));
    overlayPanel.value.hide();
  } catch (e) {
    error.value = true;
    console.error(e);
  }
}
</script>

<template>
  <Button icon="pi pi-plus" @click="(event) => overlayPanel.toggle(event)"></Button>
  <OverlayPanel ref="overlayPanel">
    <div class="flex flex-column gap-2 overlay">
      <InputGroup class="flex w-full">
        <InputGroupAddon>
          <label for="newTrackName">{{ t("tracks.name") }}</label>
        </InputGroupAddon>
        <InputText class="flex-grow-1" id="newTrackName" v-model="name"></InputText>
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
