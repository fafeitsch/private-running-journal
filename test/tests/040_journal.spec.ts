import {expect, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {trackSelectors} from '../selectors/track-selectors';
import {journalSelectors} from '../selectors/journal-selectors';

test('should load and display all available journal entries and filter with the date filter', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()
  await expect(page.getByTestId(globalSelectors.journalTab)).toContainText("Tagebuch")

  let journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(0)

  await(page.getByTestId(globalSelectors.monthInput.input)).clear()
  await(page.getByTestId(globalSelectors.monthInput.input)).fill('May 2024')
  await(page.getByTestId(globalSelectors.monthInput.input)).press('Enter')
  await expect(journalItem).toHaveCount(1)

  await journalItem.nth(0).click()
  await expect(page.getByLabel('Kommentar')).toHaveValue('Regnerisch')
  await expect(page.getByLabel('Zeit')).toHaveValue('00:51:31')
  await expect(page.getByLabel('Pace')).toHaveValue('00:05:01')
  await expect(page.getByLabel('Länge')).toHaveValue('10.3 km')

  await page.getByLabel('Nächster Monat').click()
  await expect(journalItem).toHaveCount(0)
  await page.getByLabel('Vorheriger Monat').click()
  await expect(journalItem).toHaveCount(1)
});

test('it should not display journal entry if in wrong month', async ({page}) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click()
  await expect(page.getByTestId(globalSelectors.journalTab)).toContainText("Tagebuch")

  let journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await expect(journalItem).toHaveCount(0)

  await(page.getByTestId(globalSelectors.monthInput.input)).clear()
  await(page.getByTestId(globalSelectors.monthInput.input)).fill('May 2024')
  await(page.getByTestId(globalSelectors.monthInput.input)).press('Enter')
  await expect(journalItem).toHaveCount(1)

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
