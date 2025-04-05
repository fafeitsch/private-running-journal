<script lang="ts" setup>
import InputGroupAddon from "primevue/inputgroupaddon";
import InputGroup from "primevue/inputgroup";
import { useI18n } from "vue-i18n";
import MonthChooser from '../shared/MonthChooser.vue';

const { t } = useI18n();
const fromDate = defineModel<Date>("fromDate", { required: true });
const toDate = defineModel<Date>("toDate", { required: true });
</script>
<template>
  <div class="flex gap-2 flex-col sm:flex-row">
    <InputGroup class="!w-[300px]">
      <InputGroupAddon
        ><label for="from">{{ t("dashboard.from") }}</label></InputGroupAddon
      >
      <InputGroupAddon>
        <MonthChooser id="from" v-model="fromDate" hide-border></MonthChooser>
      </InputGroupAddon>
    </InputGroup>
    <InputGroup class="!w-[300px]">
      <InputGroupAddon
        ><label for="to"> {{ t("dashboard.to") }}</label></InputGroupAddon
      >
      <InputGroupAddon>
        <MonthChooser id="to" v-model="toDate" hide-border></MonthChooser>
      </InputGroupAddon>
    </InputGroup>
  </div>
</template>
