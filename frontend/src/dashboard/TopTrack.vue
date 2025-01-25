<script setup lang="ts">
import { useI18n } from "vue-i18n";
import Card from "primevue/card";
import Tag from "primevue/tag";
import Breadcrumb from "primevue/breadcrumb";

const { t } = useI18n();
const props = defineProps<{
  id: string;
  name: string;
  parents: string[];
  length: number;
  count: number;
}>();

const breadModel: object[] = [];

props.parents.forEach((element) => {
  breadModel.push({ label: element });
});
</script>

<template>
  <Card>
    <template #title>{{ name }}</template>
    <template #subtitle>
      <Breadcrumb :model="breadModel"></Breadcrumb>
    </template>
    <template #content>
      <div class="flex flex-col">
        <span>Length: {{ length }} m</span>
        <span>Count: {{ count }}</span>
      </div>
    </template>
  </Card>
</template>

<style scoped></style>
