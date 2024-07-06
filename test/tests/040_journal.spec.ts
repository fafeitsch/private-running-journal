import {expect, Page, test} from '@playwright/test';
import {globalSelectors} from '../selectors/global-selectors';
import {journalSelectors} from '../selectors/journal-selectors';

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
  await expect(page.getByLabel('Länge')).toHaveValue('10.3 km')
  await expect(page.getByTestId(journalSelectors.lapsInput)).toHaveValue('1')

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

test('should live-update journal list after updating an entry', async ({page}) => {
  await page.goto('/');
  goToMay2024(page)
    const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
  await journalItem.nth(0).click()
  await page.getByTestId(journalSelectors.lapsInput).fill('3')
  await page.getByTestId(journalSelectors.lapsInput).blur()
  await page.getByLabel('Speichern').click()
  await expect(journalItem.nth(0)).toContainText('30,9 km')
  await page.getByTestId(journalSelectors.lapsInput).fill('1')
  await page.getByTestId(journalSelectors.lapsInput).blur()
  await page.getByLabel('Speichern').click()
  await expect(journalItem.nth(0)).toContainText('10,3 km')
})

  async function goToMay2024(page: Page) {
    await (page.getByTestId(globalSelectors.monthInput.input)).clear()
    await (page.getByTestId(globalSelectors.monthInput.input)).fill('May 2024')
    await (page.getByTestId(globalSelectors.monthInput.input)).press('Enter')
    const journalItem = page.getByTestId(journalSelectors.journalEntryItem);
    await expect(journalItem).toHaveCount(1)
  }
