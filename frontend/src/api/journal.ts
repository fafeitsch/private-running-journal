import {journal, journalEditor, journalList} from "../../wailsjs/go/models";
import {
   DeleteJournalEntry,
  GetJournalEntry,
  GetJournalListEntries,
} from "../../wailsjs/go/backend/App";
import {SaveJournalEntry} from "../../wailsjs/go/journalEditor/JournalEditor"
import SaveEntryDto = journalEditor.SaveEntryDto;
import SaveJournalEntryResultDto = journalEditor.SaveJournalEntryResultDto;

export function useJournalApi() {
  async function getJournalEntries(start: string, end: string): Promise<journalList.ListEntryDto[]> {
    return GetJournalListEntries(start, end);
  }
  async function getJournalEntry(id: string): Promise<journal.Entry> {
    return GetJournalEntry(id);
  }
  async function saveEntry(entry: SaveEntryDto): Promise<SaveJournalEntryResultDto> {
    return SaveJournalEntry(entry)
  }
  async function deleteEntry(entryId: string) {
    return DeleteJournalEntry(entryId)
  }
  return { getListEntries: getJournalEntries, getListEntry: getJournalEntry, saveEntry, deleteEntry };
}
