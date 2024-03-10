<script setup lang="ts">
import TrackTree from "./TrackTree.vue";
import SplitterPanel from "primevue/splitterpanel";
import TrackDetail from "./TrackDetail.vue";
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const selectedTrack = computed(() => route.params.trackId);
</script>

<template>
  <Splitter>
    <SplitterPanel class="flex flex-column p-2" :size="20">
      <div class="h-full overflow-hidden">
        <TrackTree />
      </div>
    </SplitterPanel>
    <SplitterPanel class="flex align-items-center justify-content-center" :size="80">
      <TrackDetail v-if="selectedTrack"></TrackDetail>
    </SplitterPanel>
  </Splitter>
</template>

<style scoped></style>
