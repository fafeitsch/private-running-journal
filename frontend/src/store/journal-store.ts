import { defineStore } from "pinia";
import { ref } from "vue";
import { journal } from "../../wailsjs/go/models";
import { useJournalApi } from "../api/journal";
import { useRouter } from "vue-router";

export const useJournalStore = defineStore("journal", () => {
  const journalApi = useJournalApi();
  const listEntries = ref<journal.ListEntry[]>([]);
  const selectedEntryId = ref<string | undefined>(undefined);
  const selectedMonth = ref<Date>(getCurrentMonth());
  const router = useRouter();

  async function loadEntries(start: string, end: string) {
    listEntries.value = [];
    listEntries.value = await journalApi.getListEntries(start, end);
  }

  function addEntryToList(entry: journal.ListEntry) {
    const date = new Date(Date.parse(entry.date));
    const month = new Date(date.getFullYear(), date.getMonth(), 1);
    selectedMonth.value = new Date(date.getFullYear(), date.getMonth(), 1);
    listEntries.value.push(entry);
  }

  function updateEntry(updatedEntry: journal.ListEntry) {
    listEntries.value = listEntries.value.map((entry) =>
      updatedEntry.id === entry.id ? updatedEntry : entry,
    );
  }

  function deleteEntry(toDelete: string) {
    listEntries.value = listEntries.value.filter((entry) => toDelete !== entry.id);
    if (selectedEntryId.value === toDelete) {
      router.replace("/journal");
    }
  }

  return {
    listEntries,
    loadEntries,
    addEntryToList,
    selectedEntryId,
    deleteEntry,
    selectedMonth,
    updateEntry,
  };
});

function getCurrentMonth() {
  const date = new Date();
  return new Date(date.getFullYear(), date.getMonth(), 1);
}
