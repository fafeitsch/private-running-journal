<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const selectedMonth = defineModel<Date>({ required: true });

function nextMonth() {
  selectedMonth.value = new Date(selectedMonth.value.setMonth(selectedMonth.value.getMonth() + 1));
}

function previousMonth() {
  selectedMonth.value = new Date(selectedMonth.value.setMonth(selectedMonth.value.getMonth() - 1));
}
</script>

<template>
  <div class="flex gap-1 align-items-center">
    <Button
      icon="pi pi-arrow-left"
      text
      :aria-label="t('shared.previousMonth')"
      @click="previousMonth()"
    ></Button>
    <Calendar
      v-model="selectedMonth"
      view="month"
      :date-format="'MM yy'"
      class="flex-grow-1"
      show-button-bar
      :pt="{
        input: { 'data-testid': 'month-chooser-component-input' },
        todayButton: {
          root: { 'data-testid': 'month-chooser-component-today-button' },
        },
        clearButton: {
          root: {style: 'display:none'}
        }
      }"
    ></Calendar>
    <Button
      icon="pi pi-arrow-right"
      text
      :aria-label="t('shared.nextMonth')"
      @click="nextMonth()"
    ></Button>
  </div>
</template>

<style scoped></style>
