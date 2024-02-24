<script lang="ts" setup>
import Sidebar from "primevue/sidebar";
import Button from "primevue/button";
import { ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const sidebarVisible = ref(false);

function toggleSidebar() {
  sidebarVisible.value = !sidebarVisible.value;
}
</script>

<template>
  <Sidebar v-model:visible="sidebarVisible" :header="t('sidenav.menu')">{{ t("menu") }}</Sidebar>
  {{ sidebarVisible }}
  <Button icon="pi pi-bars" class="absolute left-0 m-1" @click="toggleSidebar"></Button>
</template>

<style></style>
