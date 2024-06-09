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
    let month = selectedMonth.value;
    const clone = new Date(month.getTime());
    const end = new Date(clone.setMonth(month.getMonth() + 1));

    const date = Date.parse(entry.date);
    if (date >= month.getTime() && date <= end.getTime()) {
      listEntries.value.push(entry);
    }
  }

  function deleteEntry(toDelete: string) {
    listEntries.value = listEntries.value.filter((entry) => toDelete !== entry.id);
    if (selectedEntryId.value === toDelete) {
      router.replace("/journal");
    }
  }

  return { listEntries, loadEntries, addEntryToList, selectedEntryId, deleteEntry, selectedMonth };
});

function getCurrentMonth() {
  const date = new Date();
  return new Date(date.getFullYear(), date.getMonth(), 1);
}
