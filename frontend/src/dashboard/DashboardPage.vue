<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { storeToRefs } from "pinia";
import { dashboard } from "../../wailsjs/go/models";
import { useDashboardApi } from "../api/dasboard";
import ProgressSpinner from "primevue/progressspinner";
import TopTrack from "./TopTrack.vue";
import { getEndOfMonth, useDashboardStore } from "../store/dashboard-store";
import { useI18n } from "vue-i18n";
import MonthlyAnalytics from "./MonthlyAnalytics.vue";
import DistancePerMonthChart from "./DistancePerMonthChart.vue";
import RunsPerMonthChart from "./RunsPerMonthChart.vue";
import DateRangeSelection from "./DateRangeSelection.vue";
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
  <div class="flex flex-col p-2 gap-2">
    <div class="flex justify-end gap-2 items-center">
      <div v-if="!loading && !failed" class="flex gap-2">
        <span><span class="font-bold">{{t('dashboard.totalDistance')}}:</span> {{formattedTotal}}</span>
        <span><span class="font-bold">{{t('dashboard.totalRuns')}}:</span> {{data?.totalRuns}}</span>
        <span><span class="font-bold">{{t('dashboard.average')}}:</span> {{formattedAverage}}</span>
        <span><span class="font-bold">{{t('dashboard.median')}}:</span> {{formattedMedian}}</span>
      </div>
      <DateRangeSelection class="grow-0"
        v-model:from-date="selectedStartDate"
        v-model:to-date="selectedEndDate"
        @update:from-date="refresh"
        @update:to-date="refresh"
      />
    </div>
    <template v-if="failed">
      <Message severity="error">Error on API call</Message>
    </template>
    <template v-else-if="loading">
      <ProgressSpinner />
    </template>
    <div class="flex flex-col grow shrink gap-2 overflow-hidden" v-else>
      <h2 class="text-2xl">{{ t("dashboard.monthlyAnalytics") }}</h2>
      <Carousel
        :value="data.analytics"
        :num-visible="4"
        :show-navigators="data.analytics.length > 4"
        :responsive-options="[
          { breakpoint: '8000px', numVisible: 6, numScroll: 1 },
          { breakpoint: '1400px', numVisible: 4, numScroll: 1 },
          { breakpoint: '1000px', numVisible: 2, numScroll: 1 },
          { breakpoint: '600px', numVisible: 1, numScroll: 1 },
        ]"
      >
        <template #item="item">
          <MonthlyAnalytics :data="item.data" />
        </template>
      </Carousel>
      <div class="flex flex-col lg:flex-row grow shrink overflow-hidden gap-6">
        <div class="flex flex-col !min-h-[300px] lg:grow lg:shrink overflow-hidden">
          <h2 class="text-2xl">{{ t("dashboard.distancePerMonth") }}</h2>
          <div class="grow shrink overflow-hidden h-1/2">
            <DistancePerMonthChart :data="data.analytics" />
          </div>
          <h2 class="text-2xl">{{ t("dashboard.runsPerMonth") }}</h2>
          <div class="grow shrink overflow-hidden h-1/2">
            <RunsPerMonthChart :data="data.analytics" />
          </div>
        </div>
        <div class="flex flex-col gap-2 w-[610px]">
          <span class="text-2xl">{{ t("dashboard.topTracks") }}</span>
          <div class="flex flex-wrap overflow-auto grow shrink gap-2">
            <TopTrack
              v-for="track in data?.topTracks"
              :key="track.id"
              class="w-[300px]"
              :id="track.id"
              :name="track.name"
              :parents="track.parents"
              :count="track.count"
              :length="track.length"
            ></TopTrack>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
