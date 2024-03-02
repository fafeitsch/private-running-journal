import { journal } from "../../wailsjs/go/models";
import { GetJournalEntry, GetJournalListEntries } from "../../wailsjs/go/backend/App";

export function useJournalApi() {
  async function getJournalEntries(): Promise<journal.ListEntry[]> {
    return GetJournalListEntries();
  }
  async function getJournalEntry(id: string): Promise<journal.Entry> {
    return GetJournalEntry(id);
  }
  return { getListEntries: getJournalEntries, getListEntry: getJournalEntry };
}
