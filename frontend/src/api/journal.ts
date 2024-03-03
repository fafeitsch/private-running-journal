import { journal } from "../../wailsjs/go/models";
import {CreateJournalEntry, GetJournalEntry, GetJournalListEntries} from "../../wailsjs/go/backend/App";

export function useJournalApi() {
  async function getJournalEntries(): Promise<journal.ListEntry[]> {
    return GetJournalListEntries();
  }
  async function getJournalEntry(id: string): Promise<journal.Entry> {
    return GetJournalEntry(id);
  }
  async function createJournalEntry(date: string, trackId: string): Promise<journal.ListEntry> {
    return CreateJournalEntry(date, trackId)
  }
  return { getListEntries: getJournalEntries, getListEntry: getJournalEntry, createJournalEntry };
}
