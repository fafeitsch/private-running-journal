<script setup lang="ts">
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputNumber from "primevue/inputnumber";
import InputText from "primevue/inputtext";
import Popover from "primevue/popover";
import { useI18n } from "vue-i18n";
import { computed, ref } from "vue";

const { t, locale } = useI18n();

const time = defineModel<string>("time", { required: true });
const laps = defineModel<number>("laps", { required: true });
const trackLength = defineModel<number | undefined>("trackLength", { required: true });
const customLength = defineModel<boolean | undefined>("customLength", { required: true });

const pace = computed(() => {
  if (trackLength.value == null || !/\d\d:\d\d:\d\d/.test(time.value)) {
    return "";
  }
  const [hours, minutes, seconds] = time.value.split(":").map((part) => Number(part));
  const secondsTotal = hours * 60 * 60 + minutes * 60 + seconds;
  const length = customLength.value ? trackLength.value : trackLength.value * laps.value;
  let rawPace = secondsTotal / length;
  const paceHours = Math.floor(Math.ceil(rawPace) / (60 * 60));
  rawPace = Math.ceil(rawPace) % (60 * 60);
  const paceMinutes = Math.floor(Math.ceil(rawPace) / 60);
  rawPace = Math.ceil(rawPace) % 60;
  return `${paceHours.toString().padStart(2, "0")}:${paceMinutes.toString().padStart(2, "0")}:${rawPace.toFixed(0).toString().padStart(2, "0")}`;
});

const lapsWarningOverlay = ref();
</script>

<template>
  <div class="flex gap-2">
    <InputGroup>
      <Button
        icon="pi pi-pencil"
        :text="!customLength"
        class="plain-button"
        @click="customLength = !customLength"
        :aria-label="t('journal.details.overwriteLength')"
        v-tooltip="{ value: t('journal.details.overwriteLength'), showDelay: 500 }"
        data-testid="edit-custom-length-button"
      ></Button>
      <InputGroupAddon>
        <label for="length">{{
          customLength ? t("journal.details.customTotalLength") : t("journal.details.length")
        }}</label>
      </InputGroupAddon>
      <InputNumber
        id="length"
        v-model="trackLength"
        :disabled="!customLength"
        :min-fraction-digits="1"
        :max-fraction-digits="1"
        :locale="locale"
        :pt="{ pcInput: { root: { 'data-testid': 'journal-length-input', style: 'min-width:80px' } } }"
      ></InputNumber>
      <InputGroupAddon>km</InputGroupAddon>
    </InputGroup>
    <InputGroup>
      <Button
          v-if="customLength && laps > 1"
          icon="pi pi-exclamation-triangle"
          text
          severity="warning"
          class="plain-button"
          @click="(evt) => lapsWarningOverlay.toggle(evt)"
          data-testid="laps-input-no-effect-warning-indicator"
      ></Button>
      <InputGroupAddon>
        <label for="laps">{{ t("journal.details.laps") }}</label>
      </InputGroupAddon>
      <InputNumber
          :pt="{ pcInput: { root: { 'data-testid': 'laps-input' } } }"
          id="laps"
          v-model="laps"
          :min="1"
      ></InputNumber>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="time">{{ t("journal.details.time") }}</label>
      </InputGroupAddon>
      <InputText id="time" v-model="time"></InputText>
    </InputGroup>
    <InputGroup>
      <InputGroupAddon>
        <label for="pace">{{ t("journal.details.pace") }}</label>
      </InputGroupAddon>
      <InputText id="pace" :value="pace" disabled></InputText>
    </InputGroup>
  </div>
  <Popover ref="lapsWarningOverlay">
    <div style="width: 300px" data-testid="laps-input-no-effect-warning">
      {{ t('journal.details.lengthOverwrittenHint')}}
    </div>
  </Popover>
</template>

<style scoped>
.plain-button {
  border: 1px solid lightgrey;
}
</style>
