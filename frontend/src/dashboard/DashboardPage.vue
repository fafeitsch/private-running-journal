<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { dashboard } from "../../wailsjs/go/models";
import { useDashboardApi } from "../api/dasboard";
import { Slider, Splitter } from "primevue";
import Divider from "primevue/divider";
import ProgressSpinner from "primevue/progressspinner";
import MonthChooser from "./MonthChooser.vue";
import TopTrack from "./TopTrack.vue";
import { getEndOfMonth, useDashboardStore } from "../store/dashboard-store";
import { useI18n } from "vue-i18n";
import MonthlyAnalytics from "./MonthlyAnalytics.vue";
import DashboardDto = dashboard.DashboardDto;

const data = ref<DashboardDto>(
  new DashboardDto({
    totalRuns: 0,
    totalDistance: 0,
    analytics: [],
    medianDistance: 0,
    averageDistance: 0,
    topTracks: [],
  }),
);
const failed = ref<boolean>(false);
const loading = ref<boolean>(true);

const dashboardStore = useDashboardStore();
const { selectedStartDate, selectedEndDate, topTracksCount } = storeToRefs(dashboardStore);
const { t, d, n } = useI18n();

const formattedTotal = computed(
  () => n(data.value?.totalDistance / 1000, { maximumFractionDigits: 0 }) + " km",
);

const formattedAverage = computed(
  () => n(data.value?.averageDistance / 1000, { maximumFractionDigits: 0 }) + " km",
);

const formattedMedian = computed(
  () => n(data.value?.medianDistance / 1000, { maximumFractionDigits: 0 }) + " km",
);

onMounted(async () => {
  await refresh();
});

// Still as watcher because on update:model-value it called
// refresh not only on step change
watch(topTracksCount, (value) => refresh());

async function refresh() {
  loading.value = true;
  const start = new Date(selectedStartDate.value.getTime());
  start.setMinutes(start.getMinutes() - start.getTimezoneOffset());
  const end = getEndOfMonth(selectedEndDate.value);
  end.setMinutes(end.getMinutes() - end.getTimezoneOffset());
  try {
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
        {{ selectedStartDate }}
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
      <div class="flex flex-col gap-2 p-2" v-else>
        <div class="flex p-2 gap-2 w-full">
          <Panel class="flex-grow w-1/2" :header="t('dashboard.totalDistance')">{{
            formattedTotal
          }}</Panel>
          <Panel class="flex-grow w-1/2" :header="t('dashboard.totalRuns')">
            {{ data?.totalRuns }}
          </Panel>
          <Panel class="flex-grow w-1/2" :header="t('dashboard.average')">
            {{ formattedAverage }}
          </Panel>
          <Panel class="flex-grow w-1/2" :header="t('dashboard.median')">
            {{ formattedMedian }}
          </Panel>
        </div>
        <h2 class="text-2xl">{{ t("dashboard.monthlyAnalytics") }}</h2>
        <Carousel
          :value="data.analytics"
          :num-visible="4"
          :show-navigators="data.analytics.length > 4"
          :responsive-options="[
            { breakpoint: '8000px', numVisible: 4, numScroll: 1 },
            { breakpoint: '1200px', numVisible: 2, numScroll: 1 },
            { breakpoint: '700px', numVisible: 1, numScroll: 1 },
          ]"
        >
          <template #item="item">
            <MonthlyAnalytics :data="item.data" />
          </template>
        </Carousel>
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
      </div>
    </SplitterPanel>
  </Splitter>
</template>

<style scoped></style>
