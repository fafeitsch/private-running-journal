import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { useTracksApi } from "../api/tracks";
import {dashboard} from "../../wailsjs/go/models";
import DashboardDto = dashboard.DashboardDto;


export const useDashboardStore = defineStore("dashboard", () => {
  const selectedStartDate = ref<Date>(getCurrentMonth());
  const selectedEndDate = ref<Date>(getCurrentMonthEnd());
  const topTracksCount = ref<number>(10);

  return {
    selectedStartDate,
    selectedEndDate,
    topTracksCount
  };
});

function getCurrentMonth() {
  const date = new Date();
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function getCurrentMonthEnd() {
  const date = new Date();
  return getEndOfMonth(date);
}

export function getEndOfMonth(date: Date) {
  return new Date(date.getFullYear(), date.getMonth() + 1, 0);
}