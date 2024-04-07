<script setup lang="ts">
import { ref } from "vue";
import dependencyLicenses from "./vendor.LICENSE.txt?raw";
import manualDependencyLicenses from "./manual.LICENSE.txt?raw";
import { useI18n } from "vue-i18n";

const version = __APP_VERSION__

const licenseInfo = ref(dependencyLicenses);
const manualLicenseInfo = ref(manualDependencyLicenses);
const { t } = useI18n();
</script>

<template>
  <div class="p-2 flex flex-column gap-2 overflow-hidden">
    <span>{{t('about.version')}}: <span class="bold">{{version}}</span></span>
    <span>{{ t("about.dependencies") }}</span>
    <div class="flex-grow-1 flex-shrink-1 overflow-auto">
      <pre
        >{{manualLicenseInfo}}
---
        {{ licenseInfo }}
    </pre
      >
    </div>
  </div>
</template>

<style scoped>
.bold {
  font-weight: bold;
}
</style>
