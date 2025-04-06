import {expect, Page, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {journalSelectors} from '../selectors/journal-selectors';
import * as fs from 'fs';

test('should load and display all available journal entries and filter with the date filter', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()
  await expect(page.getByTestId(globalSelectors.journalTab)).toContainText("Tagebuch")

  await goToDate(page, "February 2025", 0);

  let journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(0)

  await goToDate(page, "May 2024", 1);

  await journalItem.nth(0).click()
  await expect(page.getByLabel('Kommentar')).toHaveValue('Regnerisch')
  await expect(page.getByLabel('Zeit')).toHaveValue('00:51:31')
  await expect(page.getByLabel('Pace')).toHaveValue('00:05:01')
  await expect(page.getByTestId(journalSelectors.lengthInput)).toHaveValue('10,3')
  await expect(page.getByTestId(journalSelectors.lapsInput)).toHaveValue('1')

  await page.getByLabel('Vorheriger Monat').click()
  await expect(journalItem).toHaveCount(0)
  await page.getByLabel('Nächster Monat').click()
  await expect(journalItem).toHaveCount(1)
});

test('it should not display journal entry if in wrong month', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()
  await expect(page.getByTestId(globalSelectors.journalTab)).toContainText("Tagebuch")

  let journalItem = page.getByTestId(journalSelectors.journalEntryItem);

  await goToDate(page, "May 2024", 1)

  await expect(page.getByTestId(journalSelectors.emptyState)).toBeVisible()
  await page.getByLabel("Hinzufügen").click()

  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click()
  await expect (page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(2)
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(1).click()
  await page.getByTestId(journalSelectors.entryDateInput).click()
  await page.getByTestId(journalSelectors.entryDateInputField).clear()
  await page.getByTestId(journalSelectors.entryDateInputField).fill("12.12.2024")
  await page.keyboard.press("Enter")
  await page.getByLabel("Speichern").click()

  await expect(journalItem).toHaveCount(1)

  await page.getByTestId(globalSelectors.monthInput.input).clear()
  await goToDate(page, "December 2024", 1)

  await page.getByLabel("Hinzufügen").click()

  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click()
  await expect(page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(2)
  await page.getByLabel('Höchberg').getByTestId('track-tree-selection-node-toggler').click()
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(1).click()
  await page.getByTestId(journalSelectors.entryDateInput).click()
  await page.getByTestId(journalSelectors.entryDateInputField).fill("13.12.2024")
  await page.getByLabel("Speichern").click()
  await expect(journalItem).toHaveCount(2)

  await page.getByLabel("Löschen").nth(0).click()
  await expect(page.getByTestId(journalSelectors.deleteEntryConfirmation)).toContainText('Der Eintrag wird gelöscht.')
  await page.getByLabel("Löschen").nth(1).click()
  await page.getByTestId(journalSelectors.journalEntryItem).click()
  await page.getByLabel("Löschen").nth(0).click()
  await expect(page.getByTestId(journalSelectors.deleteEntryConfirmation)).toContainText('Der Eintrag wird gelöscht.')
  await page.getByLabel("Löschen").nth(1).click()
  await expect(journalItem).toHaveCount(0)
})

test('should save file and live-update journal list after updating an entry', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()

  await goToDate(page, "May 2024", 1)
  const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await journalItem.nth(0).click()
  await page.getByTestId(journalSelectors.lapsInput).clear()
  await page.getByTestId(journalSelectors.lapsInput).pressSequentially('3')
  await page.getByTestId(journalSelectors.lapsInput).blur()
  await expect(page.getByLabel("Speichern")).toBeEnabled()
  await page.getByLabel('Speichern').click()
  let entry = JSON.parse(fs.readFileSync('testdata/journal/7c/7cd779e7-10ed-4a88-bbb2-42edeb4ad43e/entry.json') as any)
  expect(entry.customLength).toEqual(undefined)
  expect(entry.track).toEqual('544dadf2-e83c-4b14-8768-fb2b2d36483f')
  expect(entry.time).toEqual('00:51:31')
  expect(entry.comment).toEqual('Regnerisch')
  expect(entry.laps).toEqual(3)
  await expect(journalItem.nth(0)).toContainText('30,9 km')
  await page.getByTestId(journalSelectors.lapsInput).clear()
  await page.getByTestId(journalSelectors.lapsInput).pressSequentially('1')
  await page.getByTestId(journalSelectors.lapsInput).blur()
  await page.getByLabel('Speichern').click()
  await expect(journalItem.nth(0)).toContainText('10,3 km')
})

test('should correctly load and display entries with overwritten length', async ( {page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()
  await goToDate(page, "June 2024", 1)
  const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await journalItem.nth(0).click()
  await expect(page.getByTestId(journalSelectors.lapsInput)).toHaveValue('2')
  await expect(page.getByTestId(journalSelectors.lengthInput)).toHaveValue('21,0')
  await page.getByTestId(journalSelectors.customLengthWarningIndicator).click()
  await expect(page.getByTestId(journalSelectors.customLengthWarning)).toContainText('Die Streckenlänge wurde manuell überschrieben')
  await page.getByTestId(journalSelectors.customLengthWarningIndicator).click()
  await expect(page.getByTestId(journalSelectors.customLengthWarning)).not.toBeVisible()
  await page.getByTestId(journalSelectors.editCustomLengthButton).click()
  await expect(page.getByTestId(journalSelectors.customLengthWarningIndicator)).not.toBeVisible()
  await expect(page.getByTestId(journalSelectors.lengthInput)).toHaveValue('10,8')

  await page.getByLabel('Speichern').click()

  let entry = JSON.parse(fs.readFileSync('testdata/journal/d6/d6029122-d364-42c9-9784-0eb77333c43a/entry.json') as any)
  expect(entry.customLength).toEqual(undefined)
  expect(entry.track).toEqual('174b5f8d-3132-4696-bead-692335156ea3')
  expect(entry.time).toEqual('01:56:29')
  expect(entry.comment).toEqual('Hitze, deswegen mit vielen Schleifen')
  expect(entry.laps).toEqual(2)

  await page.getByTestId(journalSelectors.editCustomLengthButton).click()
  await expect(page.getByTestId(journalSelectors.customLengthWarningIndicator)).toBeVisible()
  await expect(page.getByTestId(journalSelectors.lengthInput)).toHaveValue('10,8')
  await page.getByTestId(journalSelectors.lengthInput).clear()
  await page.getByTestId(journalSelectors.lengthInput).pressSequentially("21,0")
  await expect(page.getByTestId(journalSelectors.lengthInput)).toHaveValue('21,0')
  await page.getByTestId(journalSelectors.lengthInput).blur()
  await page.getByLabel('Speichern').click()

  entry = JSON.parse(fs.readFileSync('testdata/journal/d6/d6029122-d364-42c9-9784-0eb77333c43a/entry.json') as any)
  expect(entry.customLength).toEqual(21000)
  expect(entry.track).toEqual('174b5f8d-3132-4696-bead-692335156ea3')
  expect(entry.time).toEqual('01:56:29')
  expect(entry.comment).toEqual('Hitze, deswegen mit vielen Schleifen')
  expect(entry.laps).toEqual(2)
})

test('should be possible to add new entry if there is an load error', async ({page}) => {
  await page.goto('#/journal/does-not-exist');
  await expect(page.getByTestId(journalSelectors.loadError)).toHaveCount(1)
  await page.getByLabel('Hinzufügen').click()
  await expect(page.getByTestId(journalSelectors.loadError)).toHaveCount(0)
  await expect(page.getByTestId(journalSelectors.trackSelection)).toHaveCount(1)
})

test('should be possible to view other entry on error and then create new entry and jump to current month', async ({page}) => {
  await page.goto('#/journal/does-not-exist');
  await goToDate(page, "May 2024", 1)
  await expect(page.getByTestId(journalSelectors.loadError)).toHaveCount(1)
  await page.getByTestId(journalSelectors.journalEntryItem).nth(0).click()
  await expect(page.getByLabel("Kommentar")).toHaveValue("Regnerisch")
  await page.getByLabel('Hinzufügen').click()

  await expect(page.getByLabel("Kommentar")).toHaveValue("")
  await expect(page.getByTestId(journalSelectors.lapsInput)).toHaveValue("1")
  await expect(page.getByLabel("Zeit")).toHaveValue("")
  const year = `${new Date().getFullYear()}`.padStart(4, "0");
  const month = `${new Date().getMonth() + 1}`.padStart(2, "0");
  const day = `${new Date().getDate()}`.padStart(2, "0");
  await expect(page.getByTestId(journalSelectors.entryDateInput).getByRole('combobox')).toHaveValue(`${day}.${month}.${year}`)
  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click()
  await expect(page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(5)
  await page.getByLabel('Höchberg').getByTestId('track-tree-selection-node-toggler').nth(0).click()
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(1).click()
  await page.getByTestId(journalSelectors.entryDateInput).click()
  await page.getByTestId(journalSelectors.todayButton).click()
  await page.getByLabel("Speichern").click()

  await expect(page.getByTestId(journalSelectors.journalEntryItem).getByText(`${day}.${month}.${year}Dummy` )).toBeVisible()

  await page.getByLabel("Löschen").nth(0).click()
  await expect(page.getByTestId(journalSelectors.deleteEntryConfirmation)).toContainText('Der Eintrag wird gelöscht.')
  await page.getByLabel("Löschen").nth(1).click()
})

test('should not be possible to save journal entry without track', async ({page}) => {
  await page.goto('#/journal/does-not-exist');
  await goToDate(page, "May 2024", 1)
  await page.getByLabel("Hinzufügen").click()
  await page.getByLabel("Kommentar").fill("Sehr warm")
  await expect(page.getByLabel("Speichern")).toBeDisabled()
  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click()
  await expect(page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(2)
  await page.getByLabel('Höchberg').getByTestId('track-tree-selection-node-toggler').click()
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(1).click()
  await expect(page.getByLabel("Speichern")).toBeEnabled()
})

async function goToDate(page: Page, target: string, entries: number) {
  await (page.getByTestId(globalSelectors.monthInput.input)).clear()
  await (page.getByTestId(globalSelectors.monthInput.input)).fill(target)
  await (page.getByTestId(globalSelectors.monthInput.input)).press('Enter')
  const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(entries)
}
