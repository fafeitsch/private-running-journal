<script setup lang="ts">
import { dashboard } from "../../wailsjs/go/models";
import { computed } from "vue";
import { useI18n } from "vue-i18n";

const { n, t } = useI18n();
const props = defineProps<{ data: dashboard.MonthlyAnalytics }>();

const formattedMonth = computed(() => n(props.data.month, { minimumIntegerDigits: 2 }));
const formattedTotal = computed(
  () => n(props.data.totalDistance / 1000, { maximumFractionDigits: 0 }) + " km",
);

const formattedAverage = computed(
  () => n(props.data.averageDistance / 1000, { maximumFractionDigits: 0 }) + " km",
);

const formattedMedian = computed(
  () => n(props.data.medianDistance / 1000, { maximumFractionDigits: 0 }) + " km",
);
</script>

<template>
  <Panel :header="formattedMonth + '/' + data.year" class="mx-2">
    <p>{{ t("dashboard.totalDistance") }}: {{ formattedTotal }}</p>
    <p>{{ t("dashboard.totalRuns") }}: {{ data.totalRuns }}</p>
    <p>{{ t("dashboard.average") }}: {{ formattedAverage }}</p>
    <p>{{ t("dashboard.median") }}: {{ formattedMedian }}</p>
  </Panel>
</template>

<style scoped></style>
