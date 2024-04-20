<script setup lang="ts">
import { settings } from "../../wailsjs/go/models.js";
import { defineModel, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";

const { t } = useI18n();

const model = defineModel<settings.GitSettings>({ required: true });

const emit = defineEmits<{
  (e: "update:model-value", props: settings.GitSettings): void;
}>();

watch(
  [() => model.value.enabled, () => model.value.pushAfterCommit, () => model.value.pullOnStartUp],
  () => {
    emit("update:model-value", model.value);
  },
);

const gitDescription = ref();

function moreInformationShown() {
  console.log("test");
  console.log(gitDescription.value);
  setTimeout(() => gitDescription.value?.scrollIntoView({ behavior: "smooth" }), 300);
}
</script>
<template>
  <Panel :header="'Git Integration'" :pt="{ content: { class: 'flex flex-column gap-4' } }">
    {{ t("settings.git.help") }}
    <div class="flex gap-2">
      <InputGroup>
        <InputGroupAddon>
          <Checkbox
            id="gitEnabledInput"
            v-model="model.enabled"
            :disabled="false"
            :binary="true"
          ></Checkbox>
        </InputGroupAddon>
        <InputGroupAddon class="flex-grow-1 justify-content-start">
          <label for="gitEnabledInput">{{ t("settings.git.enabled") }}</label>
        </InputGroupAddon>
      </InputGroup>
      <InputGroup>
        <InputGroupAddon>
          <Checkbox
            id="gitPushEnabledInput"
            v-model="model.pushAfterCommit"
            :disabled="!model.enabled"
            :binary="true"
          ></Checkbox>
        </InputGroupAddon>
        <InputGroupAddon class="flex-grow-1 justify-content-start">
          <label for="gitPushEnabledInput">{{ t("settings.git.pushEnabled") }}</label>
        </InputGroupAddon>
      </InputGroup>
      <InputGroup>
        <InputGroupAddon>
          <Checkbox
            id="gitPullEnabledInput"
            v-model="model.pullOnStartUp"
            :disabled="!model.enabled"
            :binary="true"
          ></Checkbox>
        </InputGroupAddon>
        <InputGroupAddon class="flex-grow-1 justify-content-start">
          <label for="gitPullEnabledInput">{{ t("settings.git.pullEnabled") }}</label>
        </InputGroupAddon>
      </InputGroup>
    </div>
    <Fieldset :legend="t('settings.git.moreInformation.title')" toggleable @toggle="moreInformationShown()">
      <p>
        {{ t("settings.git.moreInformation.header") }}
      </p>
      <ul>
        <li>{{ t("settings.git.moreInformation.executable") }}</li>
        <li>{{ t("settings.git.moreInformation.directory") }}</li>
        <li ref="gitDescription">
          {{ t("settings.git.moreInformation.pushPull") }}
        </li>
      </ul>
    </Fieldset>
  </Panel>
</template>
<style scoped>
a {
  text-decoration: underline;
}
</style>
