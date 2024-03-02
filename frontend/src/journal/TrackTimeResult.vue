<script setup lang="ts">
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputNumber from "primevue/inputnumber";
import InputText from "primevue/inputtext";
import { useI18n } from "vue-i18n";
import { computed, watch } from "vue";
import { toRefs } from "vue";

const { t } = useI18n();

const time = defineModel<string>("time", { required: true });
const laps = defineModel<number>("laps", { required: true });

const props = defineProps<{ trackLength: number | undefined }>();

const { trackLength } = toRefs(props);

const formattedLength = computed(() => {
  if(trackLength.value == null) {
    return ''
  }
  return ((laps.value * trackLength.value) / 1000).toFixed(1) + " km";
});

const pace = computed(() => {
  if (trackLength.value == null || !/\d\d:\d\d:\d\d/.test(time.value)) {
    return "";
  }
  const [hours, minutes, seconds] = time.value.split(":").map((part) => Number(part));
  const secondsTotal = hours * 60 * 60 + minutes * 60 + seconds;
  let rawPace = secondsTotal / ((trackLength.value * laps.value) / 1000);
  const paceHours = Math.floor(Math.ceil(rawPace) / (60 * 60));
  rawPace = Math.ceil(rawPace) % (60 * 60);
  const paceMinutes = Math.floor(Math.ceil(rawPace) / 60);
  rawPace = Math.ceil(rawPace) % 60;
  return `${paceHours.toString().padStart(2, "0")}:${paceMinutes.toString().padStart(2, "0")}:${rawPace.toFixed(0).toString().padStart(2, "0")}`;
});
</script>

<template>
  <div class="flex gap-2">
    <InputGroup>
      <InputGroupAddon>
        <label for="length">{{ t("journal.details.length") }}</label>
      </InputGroupAddon>
      <InputText id="length" :value="formattedLength" :disabled="true"></InputText>
      <InputGroupAddon>km</InputGroupAddon>
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
    <InputGroup>
      <InputGroupAddon>
        <label for="laps">{{ t("journal.details.laps") }}</label>
      </InputGroupAddon>
      <InputNumber id="laps" v-model="laps" :min="1"></InputNumber>
    </InputGroup>
  </div>
</template>

<style scoped></style>
