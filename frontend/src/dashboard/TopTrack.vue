<script setup lang="ts">
 defineProps<{
  id: string;
  name: string;
  parents: string[];
  length: number;
  count: number;
}>();
</script>

<template>
  <Panel :pt="{header: {style: 'padding-bottom: 0'}}">
    <template #header>
      <div class="font-bold w-full flex justify-between">
        <span class="overflow-hidden overflow-ellipsis whitespace-nowrap" v-tooltip="{showDelay: 500, value: name}">{{name}}</span>
        <span>{{count}}&timesd;</span>
      </div>
    </template>
    <span class="text-muted-color">
    {{parents[parents.length - 1]}}
    </span>
  </Panel>
</template>

<style scoped></style>
