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
  <Sidebar v-model:visible="sidebarVisible" :header="t('sidenav.menu')" :position="'top'" :modal="false" :dismissable="false">
    <ul class="list-none p-0 m-0">
      <li class="list-none p-0 m-0 flex flex-column">
        <RouterLink
          v-ripple
          class="cursor-pointer p-ripple transition-colors hover:surface-100 transition-duration-150 text-700 p-3"
          to="/"
          active-class="surface-200"
          ><span class="font-medium">{{ $t("sidenav.journal") }}</span></RouterLink
        >
      </li>
    </ul>
  </Sidebar>
  <Button icon="pi pi-bars" class="absolute right-0 bottom-0 m-1" @click="toggleSidebar"></Button>
  <router-view class="h-full w-full"></router-view>
</template>

<style></style>
