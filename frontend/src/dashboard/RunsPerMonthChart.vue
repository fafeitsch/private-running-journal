<script setup lang="ts">
import { dashboard } from "../../wailsjs/go/models";
import { computed } from "vue";
import { ChartOptions } from "chart.js";
import { useI18n } from "vue-i18n";

const props = defineProps<{ data: dashboard.MonthlyAnalytics[] }>();
const { t, n } = useI18n();

const options: ChartOptions = {
  maintainAspectRatio: false,
  responsive: true,
  plugins: { legend: { display: false }, tooltip: { enabled: false } },
  scales: {
    y: {
      ticks: {padding: 15}
    },
  }
};

const chartData = computed(() => ({
  labels: props.data.map((item) => `${n(item.month, { minimumIntegerDigits: 2 })}/${item.year}`),
  datasets: [
    {
      data: props.data.map((item) => item.totalRuns),
      backgroundColor: "rgba(75, 192, 192, 0.2)",
      borderColor: "rgba(75, 192, 192, 1)",
      borderWidth: 1,
    },
  ],
}));
</script>

<template>
  <Chart class="h-full" :data="chartData" type="bar" :options="options"></Chart>
</template>

<style scoped></style>
