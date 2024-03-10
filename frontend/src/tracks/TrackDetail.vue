<script setup lang="ts">
import { useRoute } from "vue-router";
import { computed, watch } from "vue";
import { useTrackStore } from "../store/track-store";
import { storeToRefs } from "pinia";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";

const route = useRoute();
const tracksStore = useTrackStore();
const { selectedTrackId, selectedTrack } = storeToRefs(tracksStore);
const { t } = useI18n();

watch(
  () => route.params.trackId as string,
  (trackId) => {
    selectedTrackId.value = trackId;
  },
  { immediate: true },
);

const prefix = computed(() => {
  if (!selectedTrack) {
    return "";
  }
  return selectedTrack.value?.parentNames.join("/");
});
</script>

<template>
  <div v-if="selectedTrack" class="w-full p-2">
    <InputGroup>
      <InputGroupAddon>
        <label for="nameInput">{{ t("tracks.name") }}</label>
      </InputGroupAddon>

      <InputText id="nameInput" v-model="prefix" class="flex-shrink-0" disabled :pt="{root: {size: prefix.length}}"></InputText>
      <InputText id="nameInput" v-model="selectedTrack!.name"></InputText>
    </InputGroup>
  </div>
</template>

<style scoped></style>
