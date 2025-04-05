<script setup lang="ts">
import { useI18n } from "vue-i18n";
import DatePicker from "primevue/datepicker";

defineProps<{hideBorder?: boolean}>()

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
  <div class="flex gap-1 items-center">
    <Button
      icon="pi pi-arrow-left"
      text
      :aria-label="t('shared.previousMonth')"
      @click="previousMonth()"
    ></Button>
    <DatePicker
      v-model="selectedMonth"
      view="month"
      :date-format="'MM yy'"
      class="grow"
      show-button-bar
      :pt="{
        pcInputText: { root: { 'data-testid': 'month-chooser-component-input', style: hideBorder ? 'border: none' : ''}},
        pcTodayButton: {
          root: { 'data-testid': 'month-chooser-component-today-button' },
        },
        pcClearButton: {
          root: { style: 'display:none' },
        },
      }"
    ></DatePicker>
    <Button
      icon="pi pi-arrow-right"
      text
      :aria-label="t('shared.nextMonth')"
      @click="nextMonth()"
    ></Button>
  </div>
</template>

<style scoped></style>
