<script setup lang="ts">
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import JournalList from './JournalList.vue';
</script>

<template>
  <Splitter>
    <SplitterPanel class="flex flex-column p-2" :size="20">
      <JournalList class="h-full overflow-hidden"></JournalList>
    </SplitterPanel>
    <SplitterPanel class="flex align-items-center justify-content-center" :size="80"> Panel 2 </SplitterPanel>
  </Splitter>
</template>

<style scoped></style>
