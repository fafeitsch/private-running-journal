<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { dashboard } from "../../wailsjs/go/models";
import { useDashboardApi } from "../api/dasboard";
import DashboardDto = dashboard.DashboardDto;
import { Splitter, Slider } from "primevue";
import DataView from "primevue/dataview";
import Card from "primevue/card";
import Divider from "primevue/divider";
import ProgressSpinner from "primevue/progressspinner";
import MonthChooser from "./MonthChooser.vue";
import TopTrack from "./TopTrack.vue";
import { useDashboardStore, getEndOfMonth } from "../store/dashboard-store";

const data = ref<DashboardDto | undefined>(undefined);
const failed = ref<boolean>(false);
const loading = ref<boolean>(true);

const dashboardStore = useDashboardStore();
const { selectedStartDate, selectedEndDate, topTracksCount } = storeToRefs(dashboardStore);

onMounted(async () => {
  refresh();
});

// Still as watcher because on update:model-value it called
// refresh not only on step change
watch(topTracksCount, (value) => refresh());

async function refresh() {
  loading.value = true;
  const start = selectedStartDate.value;
  const end = getEndOfMonth(selectedEndDate.value);
  try {
    // throw new Error("Boom")
    data.value = await useDashboardApi().loadDashboardApi(start, end, topTracksCount.value);
    failed.value = false;
    loading.value = false;
  } catch (e) {
    console.error(e);
    failed.value = true;
  }
}
</script>

<template>
  <Splitter>
    <SplitterPanel class="flex justify-between flex-col" :size="20">
      <div class="flex flex-col gap-2 p-2">
        <span class="text-xl">Von</span>
        <MonthChooser v-model="selectedStartDate" @update:model-value="refresh()"></MonthChooser>
        <span class="text-xl">Bis</span>
        <MonthChooser v-model="selectedEndDate" @update:model-value="refresh()"></MonthChooser>
        <div class="flex justify-between flex-col gap-2">
          <span class="text-xl">TopTracks</span>
          <div class="flex flex-col items-center gap-2">
            <label>{{ topTracksCount }}</label>
            <Slider class="w-[80%]" v-model="topTracksCount" :step="1" :min="1" :max="10"></Slider>
          </div>
        </div>
      </div>
    </SplitterPanel>

    <SplitterPanel :size="80" class="flex flex-col">
      <template v-if="failed">
        <Message severity="error">Error on API call</Message>
      </template>
      <template v-else-if="loading">
        <ProgressSpinner />
      </template>
      <template v-else>
        <div class="flex p-2 gap-2">
          <Card>
            <template #title>Total Distance</template>
            <template #content>
              <p>{{ data?.totalDistance }} m</p>
            </template>
          </Card>
          <Card>
            <template #title>Total Runs</template>
            <template #content>
              <p>{{ data?.totalRuns }}</p>
            </template>
          </Card>
        </div>
        <Divider />
        <div>
          <label class="text-xl p-2">Top Tracks:</label>
        </div>
        <div class="flex flex-wrap overflow-auto grow shrink gap-2">
          <div class="w-[300px]" v-for="track in data?.topTracks">
            <TopTrack
              class=""
              :id="track.id"
              :name="track.name"
              :parents="track.parents"
              :count="track.count"
              :length="track.length"
            ></TopTrack>
          </div>
        </div>
      </template>
    </SplitterPanel>
  </Splitter>
</template>

<style scoped></style>
