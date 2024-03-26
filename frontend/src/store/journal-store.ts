import { defineStore } from "pinia";
import { ref } from "vue";
import { journal } from "../../wailsjs/go/models";
import { useJournalApi } from "../api/journal";
import {useRouter} from 'vue-router';

export const useJournalStore = defineStore("journal", () => {
  const journalApi = useJournalApi();
  const listEntries = ref<journal.ListEntry[]>([]);
  const selectedEntryId = ref<string | undefined>(undefined);
  const router = useRouter()

  async function loadEntries() {
    listEntries.value = [];
    listEntries.value = await journalApi.getListEntries();
  }

  function addEntryToList(entry: journal.ListEntry) {
    listEntries.value.push(entry);
  }

  function deleteEntry(toDelete: string) {
    listEntries.value = listEntries.value.filter((entry) => toDelete !== entry.id);
    if (selectedEntryId.value === toDelete) {
      router.replace('/journal')
    }
  }

  return { listEntries, loadEntries, addEntryToList, selectedEntryId, deleteEntry };
});
