<script setup lang="ts">
import TrackTree from "./TrackTree.vue";
import SplitterPanel from "primevue/splitterpanel";
import TrackDetail from "./TrackDetail.vue";
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useTrackStore } from "../store/track-store";

const route = useRoute();
const selectedTrack = computed(() => route.params.trackId);
const trackStore = useTrackStore();

onMounted(async () => {
  try {
    await trackStore.loadTracks();
  } catch (e) {
    // TODO error handling
    console.error(e);
  }
});
</script>

<template>
  <Splitter>
    <SplitterPanel class="flex flex-col p-2" :size="20">
      <div class="h-full overflow-hidden">
        <TrackTree />
      </div>
    </SplitterPanel>
    <SplitterPanel class="flex items-center justify-center" :size="80">
      <TrackDetail v-if="selectedTrack"></TrackDetail>
    </SplitterPanel>
  </Splitter>
</template>

<style scoped></style>
