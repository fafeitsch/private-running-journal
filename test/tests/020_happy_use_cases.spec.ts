import { expect, test } from '@playwright/test';
import { globalSelectors } from '../selectors/global-selectors';
import { trackSelectors } from '../selectors/track-selectors';
import * as fs from 'fs';
import { journalSelectors } from '../selectors/journal-selectors';

test('should create new track and create journal entry with it, and delete track', async ({ page }) => {
  await page.goto('/');
  await page.getByTestId(globalSelectors.tracksTab).click();
  await page.getByLabel('Hinzufügen').click();
  await page.getByLabel('Streckenname').fill('Wurzelstrecke');

  await expect(page.getByLabel('Streckenname')).toHaveValue('Wurzelstrecke');
  await expect(page.getByLabel('Streckenname')).toBeVisible();
  await expect(page.getByLabel('Verwendungen')).toHaveValue('0');
  await expect(page.getByLabel('Streckenlänge')).toHaveValue('0,0');
  await page.getByLabel('Kommentar').fill("Stolpergefahr");
  await expect(page.getByLabel('Speichern')).toBeEnabled();

  await page.getByRole('button', { name: 'Vorwärts' }).click();
  await page.mouse.click(560, 500);
  await page.mouse.click(640, 450);
  await page.mouse.click(720, 390);
  await page.mouse.click(820, 290);
  await page.mouse.click(920, 200);
  await expect(page.getByLabel('Streckenlänge')).toHaveValue('0,5');
  await page.mouse.move(910, 220);
  await page.mouse.down();
  await page.mouse.move(560, 500);
  await page.mouse.up();
  await expect(page.getByTestId(trackSelectors.distanceMarker)).not.toBeVisible();
  await page.mouse.click(750, 380);
  await expect(page.getByTestId(trackSelectors.distanceMarker)).toBeVisible();
  await expect(page.getByLabel('Streckenlänge')).toHaveValue('1,0');

  await page.getByLabel('Speichern').click();
  await page.getByTestId(trackSelectors.createTrackButton).click();
  await expect(page).not.toHaveURL(/new/);
  const url = page.url();
  const id = url.substring(url.lastIndexOf('/') + 1);

  let trackNodes = page.getByTestId(trackSelectors.trackNode);
  await expect(trackNodes).toHaveCount(3);
  await expect(page.getByLabel('Speichern')).toBeDisabled();
  const track = JSON.parse(fs.readFileSync('testdata/tracks/' + id + '/info.json') as any);
  expect(track.name).toEqual('Wurzelstrecke');
  expect(track.comment).toEqual('Stolpergefahr');

  await page.goto('/');
  await page.getByTestId(globalSelectors.journalTab).click();

  await expect(page.getByTestId(journalSelectors.emptyState)).toBeVisible();
  await page.getByLabel('Hinzufügen').click();

  await page.getByTestId(journalSelectors.trackSelection).getByRole('button').click();
  await expect(page.getByTestId(journalSelectors.trackSelectionItem)).toHaveCount(3);
  await page.getByTestId(journalSelectors.trackSelectionItem).nth(2).click();
  await page.getByLabel('Speichern').click();

  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toHaveCount(1);
  const regex = /\d\d.\d\d.\d\d\d\d/;
  await expect(page.getByTestId(journalSelectors.journalEntryItem).nth(0)).toContainText(regex);
  await expect(page.getByTestId(journalSelectors.journalEntryItem).nth(0)).toContainText('Wurzelstrecke');
  await expect(page.getByTestId(journalSelectors.journalEntryItem).nth(0)).toContainText('1,0 km');

  await page.getByTestId(globalSelectors.tracksTab).click();
  await page.getByTestId(trackSelectors.trackNode).nth(2).click();
  await page.getByLabel('Löschen').click();
  await expect(page.getByTestId(trackSelectors.deleteTrackConfirmation)).toContainText('Die Strecke wird 1 mal verwendet.');
  await page.getByLabel('Löschen').nth(1).click();

  await expect(page.getByTestId(trackSelectors.trackNode)).toHaveCount(2);
  await page.getByTestId(globalSelectors.journalTab).click();
  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toHaveCount(1);
  await expect(page.getByTestId(journalSelectors.journalEntryItem).nth(0)).toContainText('Fehler beim Laden');

  await page.getByLabel('Löschen').click();
  await expect(page.getByTestId(journalSelectors.deleteEntryConfirmation)).toContainText('Der Eintrag wird gelöscht.');
  await page.getByLabel('Löschen').nth(1).click();
  await expect(page.getByTestId(journalSelectors.journalEntryItem)).toHaveCount(0);

  await expect(page.getByTestId(journalSelectors.emptyState)).toBeVisible();
});
