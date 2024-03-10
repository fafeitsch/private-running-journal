import { defineStore } from "pinia";
import {computed, ref} from "vue";
import { journal } from "../../wailsjs/go/models";
import { useJournalApi } from "../api/journal";
import {TreeNode} from 'primevue/treenode';

export const useJournalStore = defineStore("journal", () => {
  const journalApi = useJournalApi();
  const listEntries = ref<journal.ListEntry[]>([]);
  const selectedEntryId = ref<string | undefined>(undefined);

  async function loadEntries() {
    listEntries.value = [];
    listEntries.value = await journalApi.getListEntries();
  }

  function addEntryToList(entry: journal.ListEntry) {
    listEntries.value.push(entry);
  }

  return { listEntries, loadEntries, addEntryToList, selectedEntryId };
});
