import {expect, Page, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {journalSelectors} from '../selectors/journal-selectors';
import * as fs from 'fs';

test('should load and display all available journal entries and filter with the date filter', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()
  await expect(page.getByTestId(globalSelectors.journalTab)).toContainText("Tagebuch")

  let journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(0)

  await goToMay2024(page);

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
  await expect(journalItem).toHaveCount(0)

  await goToMay2024(page)

  await expect(page.getByTestId(journalSelectors.emptyState)).toBeVisible()
  await page.getByLabel("Hinzufügen").click()

  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click()
  await expect (page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(2)
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(1).click()
  await page.getByTestId(journalSelectors.todayButton).click()
  await page.getByTestId(journalSelectors.createEntryButton).click()

  await expect(journalItem).toHaveCount(1)

  await page.getByTestId(globalSelectors.monthInput.input).clear()
  await page.getByTestId(globalSelectors.monthInput.todayButton).click()

  await expect(journalItem).toHaveCount(1)

  await page.getByLabel("Hinzufügen").click()

  await page.getByTestId(journalSelectors.trackSelection).nth(1).getByRole('button').click()
  await expect (page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(2)
  await page.getByTestId(journalSelectors.trackSelection).nth(1).click()
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(0).click()
  await page.getByTestId(journalSelectors.todayButton).click()
  await page.getByTestId(journalSelectors.createEntryButton).click()
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
  await goToMay2024(page)
  const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await journalItem.nth(0).click()
  await page.getByTestId(journalSelectors.lapsInput).fill('3')
  await page.getByTestId(journalSelectors.lapsInput).blur()
  await page.getByLabel('Speichern').click()
  let entry = JSON.parse(fs.readFileSync('testdata/journal/2024/05/01/entry.json') as any)
  expect(entry.customLength).toEqual(undefined)
  expect(entry.track).toEqual('Höchberg/stadtrunde')
  expect(entry.time).toEqual('00:51:31')
  expect(entry.comment).toEqual('Regnerisch')
  expect(entry.laps).toEqual(3)
  await expect(journalItem.nth(0)).toContainText('30,9 km')
  await page.getByTestId(journalSelectors.lapsInput).fill('1')
  await page.getByTestId(journalSelectors.lapsInput).blur()
  await page.getByLabel('Speichern').click()
  await expect(journalItem.nth(0)).toContainText('10,3 km')
})

test('should correctly load and display entries with overwritten length', async ( {page}) => {
  await page.goto('/');
  await goToJune2024(page)
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

  let entry = JSON.parse(fs.readFileSync('testdata/journal/2024/06/10/entry.json') as any)
  expect(entry.customLength).toEqual(undefined)
  expect(entry.track).toEqual('Höchberg/farmrunde')
  expect(entry.time).toEqual('01:56:29')
  expect(entry.comment).toEqual('Hitze, deswegen mit vielen Schleifen')
  expect(entry.laps).toEqual(2)

  await page.getByTestId(journalSelectors.editCustomLengthButton).click()
  await expect(page.getByTestId(journalSelectors.customLengthWarningIndicator)).toBeVisible()
  await expect(page.getByTestId(journalSelectors.lengthInput)).toHaveValue('10,8')
  await page.getByTestId(journalSelectors.lengthInput).clear()
  await page.getByTestId(journalSelectors.lengthInput).fill('21,0')
  await page.getByLabel('Speichern').click()

  entry = JSON.parse(fs.readFileSync('testdata/journal/2024/06/10/entry.json') as any)
  expect(entry.customLength).toEqual(21000)
  expect(entry.track).toEqual('Höchberg/farmrunde')
  expect(entry.time).toEqual('01:56:29')
  expect(entry.comment).toEqual('Hitze, deswegen mit vielen Schleifen')
  expect(entry.laps).toEqual(2)
})

async function goToMay2024(page: Page) {
  await (page.getByTestId(globalSelectors.monthInput.input)).clear()
  await (page.getByTestId(globalSelectors.monthInput.input)).fill('May 2024')
  await (page.getByTestId(globalSelectors.monthInput.input)).press('Enter')
  const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(1)
}

async function goToJune2024(page: Page) {
  await (page.getByTestId(globalSelectors.monthInput.input)).clear()
  await (page.getByTestId(globalSelectors.monthInput.input)).fill('June 2024')
  await (page.getByTestId(globalSelectors.monthInput.input)).press('Enter')
  const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(1)
}
