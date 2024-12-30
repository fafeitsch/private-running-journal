import { journalEditor, journalList } from "../../wailsjs/go/models";
import { GetJournalListEntries } from "../../wailsjs/go/backend/App";
import { GetJournalEntry, SaveJournalEntry, DeleteJournalEntry } from "../../wailsjs/go/journalEditor/JournalEditor";
import SaveEntryDto = journalEditor.SaveEntryDto;
import SaveJournalEntryResultDto = journalEditor.SaveJournalEntryResultDto;

export function useJournalApi() {
  async function getJournalEntries(
    start: string,
    end: string,
  ): Promise<journalList.ListEntryDto[]> {
    return GetJournalListEntries(start, end);
  }

  async function getJournalEntry(id: string): Promise<journalEditor.EntryDto> {
    return GetJournalEntry(id);
  }

  async function saveEntry(entry: SaveEntryDto): Promise<SaveJournalEntryResultDto> {
    return SaveJournalEntry(entry);
  }

  async function deleteEntry(entryId: string) {
    return DeleteJournalEntry(entryId);
  }

  return {
    getListEntries: getJournalEntries,
    getJournalEntry,
    saveEntry,
    deleteEntry,
  };
}
