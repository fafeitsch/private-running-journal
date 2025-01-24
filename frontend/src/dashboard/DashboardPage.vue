<script setup lang="ts">
import { onMounted, ref } from "vue";
import { dashboard } from "../../wailsjs/go/models";
import { useDashboardApi } from "../api/dasboard";
import DashboardDto = dashboard.DashboardDto;

const data = ref<DashboardDto | undefined>(undefined);


onMounted(async () => {
  try {
    const start = new Date("2020-01-01")
    const end = new Date("2025-12-31")
    data.value = await useDashboardApi().loadDashboardApi(start, end, 10);
  } catch (e) {
    console.error(e);
  }
});
</script>

<template><div>{{data}}</div></template>

<style scoped></style>