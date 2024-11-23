import {journal, journalList} from "../../wailsjs/go/models";
import {
  CreateJournalEntry, DeleteJournalEntry,
  GetJournalEntry,
  GetJournalListEntries,
  SaveJournalEntry
} from "../../wailsjs/go/backend/App";

export function useJournalApi() {
  async function getJournalEntries(start: string, end: string): Promise<journalList.ListEntryDto[]> {
    return GetJournalListEntries(start, end);
  }
  async function getJournalEntry(id: string): Promise<journal.Entry> {
    return GetJournalEntry(id);
  }
  async function createJournalEntry(entry: journal.Entry): Promise<journal.ListEntry> {
    return CreateJournalEntry(entry)
  }
  async function saveEntry(entry: journal.Entry) {
    return SaveJournalEntry(entry)
  }
  async function deleteEntry(entryId: string) {
    return DeleteJournalEntry(entryId)
  }
  return { getListEntries: getJournalEntries, getListEntry: getJournalEntry, createJournalEntry, saveEntry, deleteEntry };
}
